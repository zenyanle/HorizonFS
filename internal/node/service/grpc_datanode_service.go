package service

import (
	"HorizonFS/internal/node/kvstore"
	"HorizonFS/internal/node/types"
	"HorizonFS/pkg/logger"
	"HorizonFS/pkg/proto"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"math"
	"sync"
	"time"
)

type DataNodeService struct {
	proto.UnimplementedDataNodeServiceServer
	Kv *kvstore.BadgerStore
}

func (s *DataNodeService) TransferData(stream grpc.BidiStreamingServer[proto.DownloadDataRequest, proto.DataChunk]) error {
	handshakeData, err := stream.Recv()
	if err != nil {
		logger.Error("handshake error ", err)
		return err
	}
	// res := s.Pool.Get(handshakeData.StreamId)
	// if res == nil {
	//	return fmt.Errorf("No pool object")
	// }
	c := NewSendController(stream, handshakeData.Request, s.Kv)
	c.Start()
	return nil
}

func (s *DataNodeService) UploadData(stream grpc.BidiStreamingServer[proto.DataChunk, proto.UploadDataResponse]) error {
	handshakeData, err := stream.Recv()
	if err != nil {
		logger.Error("handshake error ", err)
		return err
	}
	if handshakeData.ChunkId != -1 {
		logger.Error("handshake error ", err)
		return fmt.Errorf("handshake not -1")
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	var uploadErr error = nil

	fp := kvstore.NewFileProcessor(handshakeData.FileName, handshakeData.ChunkTotal, handshakeData.FileSize, &uploadErr, &wg, s.Kv)

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			logger.Info("Stream closed by server")
			break
		}
		if err != nil {
			logger.Error("Error while receiving data: ", err)
			uploadErr = err
			break
		}

		// 检查分片 ID 是否有效
		if msg.ChunkId < 0 || msg.ChunkId >= msg.ChunkTotal {
			logger.Error("Invalid chunk ID: ", msg.ChunkId)
			continue
		}

		/*			if msg.ChunkId == 0 {
					logger.Infof("%x", msg.Payload)
				}*/

		msgCopy := &types.DataChunk{
			FileName:   msg.FileName,
			ChunkId:    msg.ChunkId,
			ChunkTotal: msg.ChunkTotal,
			Payload:    make([]byte, len(msg.Payload)), // 创建 Payload 的深拷贝
			FileSize:   msg.FileSize,
		}

		copy(msgCopy.Payload, msg.Payload)

		logger.Info("第几片：", msg.ChunkId)
		logger.Info("总共几片：", msg.ChunkTotal)

		fp.ProcessChunk(int(msgCopy.ChunkId), msgCopy.Payload)
	}

	if !fp.IsComplete() {
		fp.Once.Do(wg.Done)
	}

	wg.Wait()

	if uploadErr != nil {
		logger.Error(uploadErr)
		stream.Send(&proto.UploadDataResponse{Response: uploadErr.Error()})
		return uploadErr
	}

	stream.Send(&proto.UploadDataResponse{Response: "Success"})

	return uploadErr
}

const ChunkSize = 64 * 1024

type SendController struct {
	stream   grpc.BidiStreamingServer[proto.DownloadDataRequest, proto.DataChunk]
	sendChan chan types.SenderData // 用于外部发送数据的通道
	done     chan struct{}         // 用于关闭的信号
	// mu           sync.Mutex
	kv   *kvstore.BadgerStore
	data string
	wg   sync.WaitGroup
	once *sync.Once
}

func NewSendController(stream grpc.BidiStreamingServer[proto.DownloadDataRequest, proto.DataChunk], data string, kv *kvstore.BadgerStore) *SendController {

	return &SendController{
		stream:   stream,
		sendChan: make(chan types.SenderData, 1000), // 缓冲区大小可调整
		done:     make(chan struct{}),
		kv:       kv,
		data:     data,
		once:     &sync.Once{},
	}
}

func (c *SendController) Start() {
	c.wg.Add(1)

	go c.handleSend()

	err := c.sendChunk(c.data)
	if err != nil {
		logger.Error(err)
		c.close()
		return
	}
	c.wg.Wait()
	c.close()
}

func (c *SendController) handleSend() {
	defer c.once.Do(c.wg.Done)

	for {
		select {
		case data := <-c.sendChan:
			// 发送数据到stream
			// c.mu.Lock() // sendmsg是线程安全的
			/*			if data.ChunkId == 0 {
						logger.Infof("%x", data.Payload)
					}*/
			err := c.stream.Send(&proto.DataChunk{
				FileName:   data.FileName,
				ChunkTotal: data.ChunkTotal,
				ChunkId:    data.ChunkId,
				Payload:    data.Payload,
			})

			if data.ChunkId == data.ChunkTotal-1 {
				c.once.Do(c.wg.Done)
			}

			// c.mu.Unlock()
			if err != nil {
				logger.Error("Failed to send data", err)
				return
			}
		case <-c.done:
			return
		}
	}
}

func (c *SendController) send(data types.SenderData) {
	select {
	case c.sendChan <- data:
	case <-time.After(time.Second):
		logger.Warn("Send channel full, message dropped after timeout")
	}
}

func (c *SendController) sendDataReader(reader io.Reader, chunkNum int32, fileName string) error {
	var i int32 = 0
	buffer := make([]byte, ChunkSize)
	for {
		realLen, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break // End of data
			}
			return err
		}

		copiedSlice := make([]byte, realLen) // 预先分配目标切片，大小相同
		copy(copiedSlice, buffer[:realLen])  // 使用内置的 copy() 函数高效复制

		dataSlice := types.SenderData{
			FileName:   fileName,
			ChunkTotal: chunkNum,
			ChunkId:    i,
			Payload:    copiedSlice, // 创建一个新的字节切片，其中包含 buffer 中前 realLen 个字节的内容
		}
		c.send(dataSlice)
		i += 1
	}

	return nil
}

func (c *SendController) sendChunk(id string) error {
	rd, bytelen, err := c.kv.GetReader(id)
	if err != nil {
		logger.Error(err)
		return err
	}

	numBlocks := int64(math.Ceil(float64(bytelen) / float64(ChunkSize)))

	err = c.sendDataReader(rd, int32(numBlocks), id)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (c *SendController) close() {
	close(c.done)
	close(c.sendChan)
}
