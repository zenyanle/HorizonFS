package main

import (
	"HorizonFS/internal/client/remote"
	"HorizonFS/internal/node/metadata"
	"HorizonFS/pkg/logger"
	"HorizonFS/pkg/proto"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/grpc/status"
)

func main() {
	// --- 核心参数 ---
	addr := flag.String("addr", "", "目标服务器地址 (host:port) (必需)")
	method := flag.String("method", "", "调用方法 (upload, download, getmeta, setmeta, delmeta, getfilemeta) (必需)")
	timeout := flag.Duration("timeout", 30*time.Second, "RPC 调用超时时长")

	// --- 各方法所需参数 ---
	localPath := flag.String("local-path", "", "本地文件路径 (用于上传/下载)")
	remoteName := flag.String("remote-name", "", "远程文件/键名 (用于上传/下载/获取元数据/设置元数据/删除元数据/获取文件元数据)")
	valueJSON := flag.String("value", "", "值的 JSON 字符串 (用于设置元数据，期望单个 ChunkInfo)")

	flag.Parse()

	// --- 参数校验 ---
	if *addr == "" {
		logger.Fatal("错误: 必须提供 -addr 参数")
	}
	if *method == "" {
		logger.Fatal("错误: 必须提供 -method 参数")
	}
	methodLower := strings.ToLower(*method)

	// --- 根据方法调用相应逻辑 ---
	switch methodLower {
	case "upload":
		if *localPath == "" {
			logger.Fatal("错误: 必须提供 -local-path 参数用于上传")
		}
		targetRemoteName := *remoteName
		if targetRemoteName == "" {
			targetRemoteName = filepath.Base(*localPath)
			logger.Infof("-remote-name 未提供，默认使用: %s", targetRemoteName)
		}
		err := callUpload(*addr, *localPath, *timeout)
		if err != nil {
			logger.Fatalf("上传失败: %v", err)
		}
		logger.Info("上传成功。")

	case "download":
		if *remoteName == "" {
			logger.Fatal("错误: 必须提供 -remote-name 参数用于下载")
		}
		if *localPath == "" {
			logger.Fatal("错误: 必须提供 -local-path 参数用于下载")
		}
		err := callDownload(*addr, *remoteName, *localPath, *timeout)
		if err != nil {
			logger.Fatalf("下载失败: %v", err)
		}
		logger.Info("下载成功。")

	case "getmeta":
		if *remoteName == "" {
			logger.Fatal("错误: 必须提供 -remote-name (键) 参数用于获取元数据")
		}
		chunkInfo, err := callGetMetadataLegacy(*addr, 0, *remoteName, *timeout)
		if err != nil {
			logger.Fatalf("获取元数据失败: %v", err)
		}
		metaJSON, _ := json.MarshalIndent(chunkInfo, "", "  ")
		fmt.Printf("获取元数据结果，键 '%s':\n%s\n", *remoteName, string(metaJSON))
		logger.Info("获取元数据成功。")

	case "setmeta":
		if *remoteName == "" {
			logger.Fatal("错误: 必须提供 -remote-name (键) 参数用于设置元数据")
		}
		if *valueJSON == "" {
			logger.Fatal("错误: 必须提供 -value (JSON ChunkInfo) 参数用于设置元数据")
		}
		var protoChunk proto.ChunkInfo
		if err := json.Unmarshal([]byte(*valueJSON), &protoChunk); err != nil {
			logger.Fatalf("解析 -value JSON 到 proto.ChunkInfo 失败: %v\n示例: '{\"chunk_id\":\"c1\",\"chunk_location\":\"localhost:9121\",\"chunk_status\":\"ok\"}'", err)
		}
		internalChunkInfo := metadata.ChunkInfo{
			ChunkId:       protoChunk.ChunkId,
			ChunkLocation: protoChunk.ChunkLocation,
			ChunkStatus:   protoChunk.ChunkStatus,
		}
		err := callProposeSetLegacy(*addr, 0, *remoteName, internalChunkInfo, *timeout)
		if err != nil {
			logger.Fatalf("设置元数据失败: %v", err)
		}
		logger.Info("设置元数据成功。")

	case "delmeta":
		if *remoteName == "" {
			logger.Fatal("错误: 必须提供 -remote-name (键) 参数用于删除元数据")
		}
		err := callProposeDeleteLegacy(*addr, 0, *remoteName, *timeout)
		if err != nil {
			logger.Fatalf("删除元数据失败: %v", err)
		}
		logger.Info("删除元数据成功。")

	case "getfilemeta":
		if *remoteName == "" {
			logger.Fatal("错误: 必须提供 -remote-name (键) 参数用于获取文件元数据")
		}
		chunkInfoList, err := callGetFileMetadata(*addr, *remoteName, *timeout)
		if err != nil {
			logger.Fatalf("获取文件元数据失败: %v", err)
		}
		metaJSON, _ := json.MarshalIndent(chunkInfoList, "", "  ")
		fmt.Printf("获取文件元数据结果，键 '%s':\n%s\n", *remoteName, string(metaJSON))
		logger.Info("获取文件元数据成功。")

	default:
		logger.Errorf("未知的方法: %s", *method)
		fmt.Println("可用方法: upload, download, getmeta, setmeta, delmeta, getfilemeta")
		os.Exit(1)
	}
}

func callUpload(addr, localPath string, timeout time.Duration) error {
	dataClient, err := remote.NewDataNodeClient(addr)
	if err != nil {
		return fmt.Errorf("创建 DataNodeClient 失败 %s: %w", addr, err)
	}
	defer dataClient.Close()

	logger.Infof("调用 UploadData 上传本地文件: %s 到服务器 %s", localPath, addr)
	startTime := time.Now()
	err = dataClient.UploadData(localPath)
	duration := time.Since(startTime)
	logger.Infof("UploadData 持续时间: %v", duration)
	return err
}

func callDownload(addr, remoteName, localPath string, timeout time.Duration) error {
	dataClient, err := remote.NewDataNodeClient(addr)
	if err != nil {
		return fmt.Errorf("创建 DataNodeClient 失败 %s: %w", addr, err)
	}
	defer dataClient.Close()

	logger.Infof("调用 DownloadData 下载远程文件: %s 从服务器 %s 到本地: %s", remoteName, addr, localPath)

	startTime := time.Now()
	err = dataClient.DownloadData(remoteName)
	duration := time.Since(startTime)
	logger.Infof("DownloadData 持续时间: %v", duration)
	return err
}

func callGetMetadataLegacy(addr string, index int64, key string, timeout time.Duration) (metadata.ChunkInfo, error) {
	metaClient, err := remote.NewMetadataClient(addr)
	if err != nil {
		return metadata.ChunkInfo{}, fmt.Errorf("创建 MetadataClient 失败 %s: %w", addr, err)
	}
	defer metaClient.Close()

	logger.Infof("调用 GetMetadata (旧版) 获取键: %s 的元数据，服务器 %s", key, addr)
	startTime := time.Now()
	result, err := metaClient.GetMetadata(index, key)
	duration := time.Since(startTime)
	logger.Infof("GetMetadata (旧版) 持续时间: %v", duration)
	return result, err
}

func callProposeSetLegacy(addr string, index int64, key string, value metadata.ChunkInfo, timeout time.Duration) error {
	metaClient, err := remote.NewMetadataClient(addr)
	if err != nil {
		return fmt.Errorf("创建 MetadataClient 失败 %s: %w", addr, err)
	}
	defer metaClient.Close()

	logger.Infof("调用 ProposeSet (旧版) 设置键: %s 的元数据，服务器 %s 值: %+v", key, addr, value)
	startTime := time.Now()
	err = metaClient.ProposeSet(index, key, value)
	duration := time.Since(startTime)
	logger.Infof("ProposeSet (旧版) 持续时间: %v", duration)
	return err
}

func callProposeDeleteLegacy(addr string, index int64, key string, timeout time.Duration) error {
	metaClient, err := remote.NewMetadataClient(addr)
	if err != nil {
		return fmt.Errorf("创建 MetadataClient 失败 %s: %w", addr, err)
	}
	defer metaClient.Close()

	logger.Infof("调用 ProposeDelete (旧版) 删除键: %s 的元数据，服务器 %s", key, addr)
	startTime := time.Now()
	err = metaClient.ProposeDelete(index, key)
	duration := time.Since(startTime)
	logger.Infof("ProposeDelete (旧版) 持续时间: %v", duration)
	return err
}

func callGetFileMetadata(addr string, key string, timeout time.Duration) ([]metadata.ChunkInfo, error) {
	metaClient, err := remote.NewMetadataClient(addr)
	if err != nil {
		return nil, fmt.Errorf("创建 MetadataClient 失败 %s: %w", addr, err)
	}
	defer metaClient.Close()

	logger.Infof("调用 GetFileMetadata 获取键: %s 的文件元数据，服务器 %s", key, addr)
	startTime := time.Now()
	result, err := metaClient.GetFileMetadata(key)
	duration := time.Since(startTime)
	logger.Infof("GetFileMetadata 持续时间: %v", duration)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			return nil, fmt.Errorf("gRPC 错误: 代码 = %s, 信息 = %s", st.Code(), st.Message())
		}
		return nil, err
	}
	return result, nil
}
