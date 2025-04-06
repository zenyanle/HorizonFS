package kvstore

import (
	"HorizonFS/pkg/logger"
	"sync"
	"sync/atomic"
)

// FileProcessor 处理流式传输的文件分片
type FileProcessor struct {
	fileName      string
	totalChunks   int32
	receivedCount atomic.Int32
	fileSize      int64

	// 分片存储
	chunks map[int][]byte
	mu     sync.RWMutex

	// 收到了全部的分片
	isComplete atomic.Bool

	// 分片组装完成，生命周期结束
	isDone atomic.Bool

	uploadErr error

	wg *sync.WaitGroup

	Once *sync.Once

	Kv *BadgerStore
}

// NewFileProcessor 创建新的文件处理器
func NewFileProcessor(fileName string, totalChunks int32, fileSize int64, uploadErr *error, wg *sync.WaitGroup, kv *BadgerStore) *FileProcessor {
	return &FileProcessor{
		fileName:    fileName,
		totalChunks: totalChunks,
		fileSize:    fileSize,
		chunks:      make(map[int][]byte),
		uploadErr:   *uploadErr,
		wg:          wg,
		Once:        &sync.Once{},
		Kv:          kv,
	}
}

// ProcessChunk 处理单个文件分片，并发安全
func (fp *FileProcessor) ProcessChunk(chunkIndex int, data []byte) {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	// 再次检查是否已完成，防止竞态
	if fp.isComplete.Load() {
		logger.Errorf("file %s processing already completed", fp.fileName)
		return
	}

	// 检查是否已存储该分片
	if _, exists := fp.chunks[chunkIndex]; exists {
		return // 避免重复存储
	}

	// 存储分片
	// fp.chunks[chunkIndex] = make([]byte, len(data))
	// copy(fp.chunks[chunkIndex], data)
	fp.chunks[chunkIndex] = data

	/*	if chunkIndex == 0 {
		logger.Infof("%x", data)
	}*/

	// 增加计数
	count := fp.receivedCount.Add(1)

	// 确保 `WriteToFile` 只被调用一次
	if count == atomic.LoadInt32(&fp.totalChunks) {
		fp.isComplete.Store(true)
		go fp.WriteToFile()
	}
}

// WriteToFile 将完整文件写入磁盘
func (fp *FileProcessor) WriteToFile() {
	defer fp.Once.Do(fp.wg.Done)

	tmp := make([]byte, fp.fileSize)
	offset := 0
	for i := 0; i < int(fp.totalChunks); i++ {
		chunk := fp.chunks[i]
		copy(tmp[offset:], chunk)
		offset += len(chunk)
	}

	err := fp.Kv.SetValue(fp.fileName, tmp)
	if err != nil {
		logger.Error("kv set error ", err)
		return
	}

	// 清理并通知完成
	fp.cleanUp()
	fp.isDone.Store(true)
	return
}

// IsComplete 检查文件是否接收完成
func (fp *FileProcessor) IsComplete() bool {
	return fp.isComplete.Load()
}

func (fp *FileProcessor) IsDone() bool {
	return fp.isDone.Load()
}

// cleanUp 清理资源
func (fp *FileProcessor) cleanUp() {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	// 清空分片数据
	fp.chunks = nil
}
