package store

import (
	"HorizonFS/pkg/logger"
	"fmt"
	"os"
	"path/filepath"
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
}

// NewFileProcessor 创建新的文件处理器
func NewFileProcessor(fileName string, totalChunks int32, fileSize int64, uploadErr *error, wg *sync.WaitGroup) *FileProcessor {
	return &FileProcessor{
		fileName:    fileName,
		totalChunks: totalChunks,
		fileSize:    fileSize,
		chunks:      make(map[int][]byte),
		uploadErr:   *uploadErr,
		wg:          wg,
		Once:        &sync.Once{},
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

	logger.Info("创建文件了")
	path := filepath.Join("test", "downloaded", fp.fileName)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fp.uploadErr = err
		logger.Errorf("failed to create file %s: %v", fp.fileName, err)
		return
	}
	defer func() {
		cerr := file.Close()
		if cerr != nil {
			fp.uploadErr = cerr
			logger.Error(cerr)
		}
	}()

	if err := file.Truncate(fp.fileSize); err != nil {
		fp.uploadErr = err
		logger.Errorf("failed to truncate file %s: %v", fp.fileName, err)
		return
	}
	logger.Info("在写文件了")
	logger.Info("总分片： ", fp.totalChunks)
	// 跟踪写入位置
	var offset int64 = 0

	// 按顺序写入分片
	for i := 0; i < int(fp.totalChunks); i++ {
		fp.mu.RLock()
		chunk, ok := fp.chunks[i]
		fp.mu.RUnlock()
		// logger.Info(offset)
		if !ok {
			logger.Info("OK断点")
			fp.uploadErr = fmt.Errorf("missing chunk %d for file %s", i, fp.fileName)
			logger.Errorf("missing chunk %d for file %s", i, fp.fileName)
			return
		}

		//logger.Infof("%x", chunk)

		// 使用当前offset写入
		if _, err := file.WriteAt(chunk, offset); err != nil {
			fp.uploadErr = err
			logger.Errorf("failed to write chunk %d: %v", i, err)
			return
		}
		// logger.Info("offset更新点")
		// 更新offset
		offset += int64(len(chunk))
		// logger.Info(offset)
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
