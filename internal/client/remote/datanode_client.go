package remote

import (
	"HorizonFS/internal/client/store"
	"HorizonFS/pkg/logger"
	"HorizonFS/pkg/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"math"
	"os"
	"sync"
	"time"
)

const (
	ChunkSize = 64 * 1024
)

// DataNodeClient 用于处理与TransferService的连接
type DataNodeClient struct {
	conn   *grpc.ClientConn
	client proto.DataNodeServiceClient
	ctx    context.Context
	cancel context.CancelFunc
}

// NewDataNodeClient 创建一个新的客户端连接到传输服务
func NewDataNodeClient(serverAddr string) (*DataNodeClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)

	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cancel()
		return nil, fmt.Errorf("连接服务器失败: %v", err)
	}

	client := proto.NewDataNodeServiceClient(conn)

	return &DataNodeClient{
		conn:   conn,
		client: client,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Close 关闭客户端连接
func (c *DataNodeClient) Close() {
	c.cancel()
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *DataNodeClient) DownloadData(requestID string) error {
	stream, err := c.client.DownloadData(c.ctx)
	if err != nil {
		return fmt.Errorf("创建流失败: %v", err)
	}

	// 发送握手信息
	err = stream.Send(&proto.DownloadDataRequest{
		Request: requestID,
	})
	if err != nil {
		return fmt.Errorf("发送握手失败: %v", err)
	}

	// 处理接收到的数据块
	wg := sync.WaitGroup{}
	wg.Add(1)

	var uploadErr error = nil

	// 接收第一个块以获取文件信息
	chunk, err := stream.Recv()
	if err != nil {
		logger.Error("接收第一个数据块失败: ", err)
		return err
	}

	// 使用file_processor处理文件
	fp := store.NewFileProcessor(chunk.FileName, chunk.ChunkTotal, chunk.FileSize, &uploadErr, &wg)

	// 处理第一个块
	payloadCopy := make([]byte, len(chunk.Payload))
	copy(payloadCopy, chunk.Payload)
	fp.ProcessChunk(int(chunk.ChunkId), payloadCopy)

	// 处理剩余的块
	go func() {
		for {
			chunk, err := stream.Recv()
			if err == io.EOF {
				logger.Info("数据传输完成")
				break
			}
			if err != nil {
				logger.Error("接收数据块时出错: ", err)
				uploadErr = err
				break
			}

			// 创建载荷的深拷贝
			payloadCopy := make([]byte, len(chunk.Payload))
			copy(payloadCopy, chunk.Payload)

			logger.Info("第几片: ", chunk.ChunkId)
			logger.Info("总共几片: ", chunk.ChunkTotal)

			// 处理数据块
			fp.ProcessChunk(int(chunk.ChunkId), payloadCopy)
		}

		// 如果尚未完成，确保waitgroup完成
		if !fp.IsComplete() {
			fp.Once.Do(wg.Done)
		}
	}()

	// 等待所有数据块处理完成
	wg.Wait()

	if uploadErr != nil {
		logger.Error("传输过程中发生错误:", uploadErr)
		return uploadErr
	}

	return nil
}

// UploadData 上传文件到服务器
func (c *DataNodeClient) UploadData(filePath string) error {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 计算总块数
	fileSize := fileInfo.Size()
	chunkTotal := int32(math.Ceil(float64(fileSize) / float64(ChunkSize)))

	// 创建流
	stream, err := c.client.UploadData(c.ctx)
	if err != nil {
		return fmt.Errorf("创建上传流失败: %v", err)
	}

	// 创建一个错误通道和完成通道
	errChan := make(chan error, 1)
	// doneChan := make(chan struct{})
	ackChan := make(chan struct{})

	// 启动一个goroutine来接收服务器响应
	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				// close(doneChan)
				return
			}
			if err != nil {
				errChan <- fmt.Errorf("接收服务器响应失败: %v", err)
				return
			}
			logger.Info("收到服务器响应: ", response.Response)

			if response.Response != "Success" {
				errChan <- fmt.Errorf("发送端出错")
				return
			}

			close(ackChan)
			return
		}
	}()

	// 发送握手
	err = stream.Send(&proto.DataChunk{
		FileName:   fileInfo.Name(),
		ChunkId:    -1, // 握手块
		ChunkTotal: chunkTotal,
		FileSize:   fileSize,
		Payload:    nil,
	})
	if err != nil {
		return fmt.Errorf("发送握手失败: %v", err)
	}

	// 读取并发送文件块
	buffer := make([]byte, ChunkSize)
	var chunkId int32 = 0

	for {
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("读取文件错误: %v", err)
		}

		// 创建数据副本
		payload := make([]byte, bytesRead)
		copy(payload, buffer[:bytesRead])

		// 发送块
		err = stream.Send(&proto.DataChunk{
			FileName:   fileInfo.Name(),
			ChunkId:    chunkId,
			ChunkTotal: chunkTotal,
			FileSize:   fileSize,
			Payload:    payload,
		})
		if err != nil {
			return fmt.Errorf("发送数据块 %d 错误: %v", chunkId, err)
		}

		logger.Info("第几片: ", chunkId)
		logger.Info("总共几片: ", chunkTotal)
		chunkId++

		// 检查是否有错误
		select {
		case err := <-errChan:
			return err
		default:
			// 继续
		}
	}

	// 通知服务器客户端已完成发送
	err = stream.CloseSend()
	if err != nil {
		return fmt.Errorf("关闭发送流错误: %v", err)
	}

	// 等待接收完成或出错
	select {
	case err := <-errChan:
		return err
	// case <-doneChan:
	// 	logger.Info("上传完成")
	case <-ackChan:
		logger.Info("服务端确认成功")
	case <-time.After(30 * time.Second):
		return fmt.Errorf("等待服务器响应超时")
	}

	return nil
}

// UploadDataFromReader 从Reader上传数据到服务器
func (c *DataNodeClient) UploadDataFromReader(reader io.Reader, fileName string, fileSize int64) error {
	// 计算总块数
	chunkTotal := int32(math.Ceil(float64(fileSize) / float64(ChunkSize)))

	// 创建流
	stream, err := c.client.UploadData(c.ctx)
	if err != nil {
		return fmt.Errorf("创建上传流失败: %v", err)
	}

	// 创建一个错误通道和完成通道
	errChan := make(chan error, 1)
	doneChan := make(chan struct{})

	// 启动一个goroutine来接收服务器响应
	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				close(doneChan)
				return
			}
			if err != nil {
				errChan <- fmt.Errorf("接收服务器响应失败: %v", err)
				return
			}
			logger.Info("收到服务器响应: ", response)
		}
	}()

	// 发送握手
	err = stream.Send(&proto.DataChunk{
		FileName:   fileName,
		ChunkId:    -1, // 握手块
		ChunkTotal: chunkTotal,
		FileSize:   fileSize,
		Payload:    nil,
	})
	if err != nil {
		return fmt.Errorf("发送握手失败: %v", err)
	}

	// 读取并发送数据块
	buffer := make([]byte, ChunkSize)
	var chunkId int32 = 0

	for {
		bytesRead, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("读取数据错误: %v", err)
		}

		if bytesRead == 0 {
			break
		}

		// 创建数据副本
		payload := make([]byte, bytesRead)
		copy(payload, buffer[:bytesRead])

		// 发送块
		err = stream.Send(&proto.DataChunk{
			FileName:   fileName,
			ChunkId:    chunkId,
			ChunkTotal: chunkTotal,
			FileSize:   fileSize,
			Payload:    payload,
		})
		if err != nil {
			return fmt.Errorf("发送数据块 %d 错误: %v", chunkId, err)
		}

		logger.Info("第几片: ", chunkId)
		logger.Info("总共几片: ", chunkTotal)
		chunkId++

		// 检查是否有错误
		select {
		case err := <-errChan:
			return err
		default:
			// 继续
		}
	}

	// 通知服务器客户端已完成发送
	err = stream.CloseSend()
	if err != nil {
		return fmt.Errorf("关闭发送流错误: %v", err)
	}

	// 等待接收完成或出错
	select {
	case err := <-errChan:
		return err
	case <-doneChan:
		logger.Info("上传完成")
	case <-time.After(30 * time.Second):
		return fmt.Errorf("等待服务器响应超时")
	}

	return nil
}
