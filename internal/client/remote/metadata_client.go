package remote

import (
	"HorizonFS/internal/node/metadata"
	"HorizonFS/pkg/logger"
	"HorizonFS/pkg/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"time"
)

type MetadataClient struct {
	conn   *grpc.ClientConn
	client proto.MetadataServiceClient
	ctx    context.Context
	cancel context.CancelFunc
}

func NewMetadataClient(serverAddr string) (*MetadataClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cancel()
		return nil, fmt.Errorf("连接服务器失败: %v", err)
	}

	client := proto.NewMetadataServiceClient(conn)

	return &MetadataClient{
		conn:   conn,
		client: client,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (c *MetadataClient) GetMetadata(key string) (metadata.ChunkInfo, error) {
	req := &proto.GetMetadataRequest{Key: key}

	res, err := c.client.GetMetadata(c.ctx, req)

	if err != nil {
		// 尝试解析 gRPC 错误状态
		st, ok := status.FromError(err)
		if ok {
			logger.Errorf("GetMetadata 调用失败: code=%s, message='%s'", st.Code(), st.Message())
			return metadata.ChunkInfo{}, err
		} else {
			logger.Errorf("GetMetadata 调用失败 (非 gRPC 错误): %v", err)
			return metadata.ChunkInfo{}, err
		}
	}
	return metadata.ChunkInfo{
		ChunkId:       res.Value.ChunkId,
		ChunkLocation: res.Value.ChunkLocation,
		ChunkStatus:   res.Value.ChunkStatus,
	}, nil
}

func (c *MetadataClient) ProposeSet(key string, chunkInfo metadata.ChunkInfo) error {

	req := &proto.ProposeSetRequest{
		Key: key,
		Value: &proto.ChunkInfo{
			ChunkId:       chunkInfo.ChunkId,
			ChunkLocation: chunkInfo.ChunkLocation,
			ChunkStatus:   chunkInfo.ChunkStatus,
		},
	}

	_, err := c.client.ProposeSet(c.ctx, req)

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			logger.Errorf("ProposeSet 调用失败: code=%s, message='%s'", st.Code(), st.Message())
			return err
		} else {
			logger.Errorf("ProposeSet 调用失败: %v", err)
			return err
		}
	}
	return nil
}

func (c *MetadataClient) ProposeDelete(key string) error {

	req := &proto.ProposeDeleteRequest{Key: key}

	_, err := c.client.ProposeDelete(c.ctx, req)

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			logger.Errorf("ProposeDelete 调用失败: code=%s, message='%s'", st.Code(), st.Message())
		} else {
			logger.Errorf("ProposeDelete 调用失败: %v", err)
		}
	}
	return nil
}
