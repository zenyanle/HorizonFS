package main

import (
	log "HorizonFS/pkg/logger"
	"context"
	"encoding/json" // 用于解析 ProposeSet 的 value 参数
	"flag"
	"os"
	"strings"
	"time"

	pb "HorizonFS/pkg/proto" // 你的 proto 包路径 (请确保正确)

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status" // 用于解析 gRPC 错误
	// 如果需要更可靠的 JSON<->Proto 转换，可以引入:
	// "google.golang.org/protobuf/encoding/protojson"
)

func main() {
	// --- 1. 定义命令行参数 ---
	serverAddr := flag.String("addr", "localhost:50051", "元数据服务地址和端口")
	method := flag.String("method", "", "要调用的方法 (GetMetadata, ProposeSet, ProposeDelete)")
	key := flag.String("key", "", "操作的 Key (例如: 文件名)")
	// value 参数现在只用于 ProposeSet
	value := flag.String("value", "", "对于 ProposeSet: ChunkInfo 的 JSON 字符串 (例如: '{\"chunk_id\":\"id1\",\"chunk_location\":\"loc1\",\"chunk_status\":\"available\"}')")
	timeout := flag.Duration("timeout", 10*time.Second, "RPC 调用超时时间")

	// 解析命令行参数
	flag.Parse()

	// --- 2. 校验基本参数 ---
	if *serverAddr == "" || *method == "" {
		log.Println("错误: 必须提供服务器地址 (--addr) 和方法名 (--method)")
		flag.Usage()
		os.Exit(1)
	}
	// 根据方法校验必需的参数
	methodLower := strings.ToLower(*method)
	if methodLower != "getmetadata" && methodLower != "proposedelete" && methodLower != "proposeset" {
		log.Printf("错误: 未知的方法名 '%s'", *method)
		log.Println("支持的方法: GetMetadata, ProposeSet, ProposeDelete")
		os.Exit(1)
	}
	if (methodLower == "getmetadata" || methodLower == "proposedelete") && *key == "" {
		log.Fatalf("错误: %s 方法需要 --key 参数", *method)
	}
	if methodLower == "proposeset" && (*key == "" || *value == "") {
		log.Fatalf("错误: ProposeSet 方法需要 --key 和 --value (JSON) 参数")
	}

	// --- 3. 建立 gRPC 连接 ---
	log.Printf("正在连接到元数据服务器: %s", *serverAddr)
	conn, err := grpc.Dial(*serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // 本地测试，无加密
		grpc.WithBlock(),                // 阻塞直到连接成功或超时 (由下面的 context 控制)
		grpc.WithTimeout(5*time.Second), // 设置 Dial 的超时，避免无限等待
	)
	if err != nil {
		log.Fatalf("无法连接 gRPC 服务器: %v", err)
	}
	defer conn.Close()
	log.Println("连接成功!")

	// --- 4. 创建上下文和 MetadataService 客户端 ---
	ctx, cancel := context.WithTimeout(context.Background(), *timeout) // 使用命令行指定的超时
	defer cancel()

	metadataClient := pb.NewMetadataServiceClient(conn) // 创建客户端

	// --- 5. 根据方法名调用对应的函数 ---
	switch methodLower {
	case "getmetadata":
		callGetMetadata(ctx, metadataClient, *key)
	case "proposeset":
		callProposeSet(ctx, metadataClient, *key, *value)
	case "proposedelete":
		callProposeDelete(ctx, metadataClient, *key)
		// default 分支在上面已处理
	}

	log.Println("操作完成.")
}

// --- 调用 GetMetadata ---
func callGetMetadata(ctx context.Context, client pb.MetadataServiceClient, key string) {
	log.Printf("调用 GetMetadata, Key: %s", key)
	req := &pb.GetMetadataRequest{Key: key}

	startTime := time.Now() // 记录开始时间
	res, err := client.GetMetadata(ctx, req)
	duration := time.Since(startTime) // 计算耗时

	if err != nil {
		// 尝试解析 gRPC 错误状态
		st, ok := status.FromError(err)
		if ok {
			log.Fatalf("GetMetadata 调用失败: code=%s, message='%s', 耗时: %v", st.Code(), st.Message(), duration)
		} else {
			log.Fatalf("GetMetadata 调用失败 (非 gRPC 错误): %v, 耗时: %v", err, duration)
		}
	}

	if res.Found {
		// 使用 %+v 可以更详细地打印结构体，包括字段名
		log.Printf("响应: Found=true, Value=%+v, 耗时: %v", res.Value, duration)
	} else {
		log.Printf("响应: Found=false, Key '%s' 未找到, 耗时: %v", key, duration)
	}
}

// --- 调用 ProposeSet ---
func callProposeSet(ctx context.Context, client pb.MetadataServiceClient, key, valueJSON string) {
	// 解析 JSON 到 ChunkInfo 结构体
	var chunkInfo pb.ChunkInfo
	// 尝试使用标准库 json unmarshal
	// 注意：如果字段名不匹配 (proto snake_case vs json camelCase)，可能需要 protojson
	if err := json.Unmarshal([]byte(valueJSON), &chunkInfo); err != nil {
		log.Fatalf("错误: 解析 --value 的 JSON 失败: %v\n请确保 JSON 格式正确且字段名与 proto 定义匹配 (例如: {\"chunk_id\":\"id1\",\"chunk_location\":\"loc1\",\"chunk_status\":\"available\"})", err)
	}
	// 使用 protojson (如果引入了库):
	// if err := protojson.Unmarshal([]byte(valueJSON), &chunkInfo); err != nil {
	//     log.Fatalf("Error parsing --value JSON with protojson: %v", err)
	// }

	log.Printf("调用 ProposeSet, Key: %s, Value: %+v", key, chunkInfo)
	req := &pb.ProposeSetRequest{
		Key:   key,
		Value: &chunkInfo, // 传入解析后的 proto 结构体指针
	}

	startTime := time.Now()
	res, err := client.ProposeSet(ctx, req)
	duration := time.Since(startTime)

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			log.Fatalf("ProposeSet 调用失败: code=%s, message='%s', 耗时: %v", st.Code(), st.Message(), duration)
		} else {
			log.Fatalf("ProposeSet 调用失败: %v, 耗时: %v", err, duration)
		}
	}
	log.Printf("响应: Success=%t, ErrorMessage='%s', 耗时: %v", res.Success, res.ErrorMessage, duration)
	// 提醒用户 Propose 成功的弱保证性质
	log.Println("注意: ProposeSet 成功响应仅表示提案已被 Leader 节点接收，不保证最终被 commit。")
}

// --- 调用 ProposeDelete ---
func callProposeDelete(ctx context.Context, client pb.MetadataServiceClient, key string) {
	log.Printf("调用 ProposeDelete, Key: %s", key)
	req := &pb.ProposeDeleteRequest{Key: key}

	startTime := time.Now()
	res, err := client.ProposeDelete(ctx, req)
	duration := time.Since(startTime)

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			log.Fatalf("ProposeDelete 调用失败: code=%s, message='%s', 耗时: %v", st.Code(), st.Message(), duration)
		} else {
			log.Fatalf("ProposeDelete 调用失败: %v, 耗时: %v", err, duration)
		}
	}
	log.Printf("响应: Success=%t, ErrorMessage='%s', 耗时: %v", res.Success, res.ErrorMessage, duration)
	log.Println("注意: ProposeDelete 成功响应仅表示提案已被 Leader 节点接收，不保证最终被 commit。")
}
