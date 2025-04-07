package remote

import (
	"HorizonFS/pkg/logger"
	"HorizonFS/pkg/proto"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func main() {
	serverAddr := "localhost:50051"
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("无法连接 gRPC 服务器: %v", err)
	}
	// 在 main 函数结束时关闭连接
	defer conn.Close()
	metadataClient := proto.NewMetadataServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // 10 秒超时
	defer cancel()
	metadataClient.GetMetadata(ctx)
}
