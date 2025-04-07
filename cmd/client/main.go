package main

import (
	"HorizonFS/internal/client/remote"
	"HorizonFS/pkg/logger"
	"path/filepath"
)

// 示例用法
func main() {
	client, err := remote.NewDataNodeClient("localhost:22380")
	if err != nil {
		logger.Fatal("创建客户端失败:", err)
	}
	defer client.Close()

	// 下载数据
	err = client.DownloadData("test.png")
	if err != nil {
		logger.Error("传输失败:", err)
	}

	// 上传文件
	err = client.UploadData(filepath.Join("test", "sample", "test.png"))
	if err != nil {
		logger.Error("上传失败:", err)
	}
}
