package service

import (
	"HorizonFS/internal/node/metadata"
	"HorizonFS/pkg/logger"
	"HorizonFS/pkg/proto"
	"context"
	json "github.com/bytedance/sonic"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MetadataService struct {
	proto.UnimplementedMetadataServiceServer

	Store *metadata.Metastore
}

/*// NewMetadataService 创建 MetadataService 的 gRPC 服务器实例
// 参数需要传入你的元数据存储/Raft交互实例
func NewMetadataService(store *metadata.Metastore) *MetadataService {

	return &MetadataService{
		Store: store, // 将传入的 metastore 实例赋值给服务器
	}
}*/

func (s *MetadataService) ProposeSet(ctx context.Context, req *proto.ProposeSetRequest) (*proto.ProposeSetResponse, error) {
	// 1. 输入参数校验
	if req.Key == "" {
		logger.Warn("ProposeSet 请求空了")
		return nil, status.Error(codes.InvalidArgument, "元数据键 (key) 不能为空")
	}
	if req.Value == nil {
		logger.Warn("ProposeSet 请求缺少 value")
		return nil, status.Error(codes.InvalidArgument, "元数据值 (value) 不能为空")
	}

	chunkInfo := metadata.ChunkInfo{
		ChunkId:       req.Value.ChunkId,
		ChunkLocation: req.Value.ChunkLocation,
		ChunkStatus:   req.Value.ChunkStatus,
	}

	valueBytes, err := json.Marshal(chunkInfo) // 将 proto 消息序列化为 JSON
	if err != nil {
		return nil, status.Errorf(codes.Internal, "序列化元数据值失败: %v", err)
	}

	err = s.Store.Propose(ctx, "SET", req.Key, string(valueBytes))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "超时: %v", err)
	}
	return &proto.ProposeSetResponse{Success: true}, nil
}

func (s *MetadataService) GetMetadata(ctx context.Context, req *proto.GetMetadataRequest) (*proto.GetMetadataResponse, error) {
	if req.Key == "" {
		logger.Warn("GetMetadata 请求缺少 key")
		return nil, status.Error(codes.InvalidArgument, "元数据键 (key) 不能为空")
	}

	logger.Info("收到 GetMetadata 请求")

	internalValue, found := s.Store.Lookup(req.Key)

	if !found {
		logger.Info("GetMetadata 未找到 key")
		return &proto.GetMetadataResponse{Found: false}, nil
	}

	protoValue := convertInternalToProtoChunkInfo(internalValue)
	if protoValue == nil {
		logger.Error("内部 ChunkInfo 转换到 proto 失败")
		return nil, status.Errorf(codes.Internal, "内部数据转换失败")
	}

	logger.Info("成功获取元数据")
	return &proto.GetMetadataResponse{Value: protoValue, Found: true}, nil
}

func (s *MetadataService) ProposeDelete(ctx context.Context, req *proto.ProposeDeleteRequest) (*proto.ProposeDeleteResponse, error) {

	if req.Key == "" {
		logger.Warn("ProposeDelete 请求缺少 key")
		return nil, status.Error(codes.InvalidArgument, "元数据键 (key) 不能为空")
	}

	logger.Info("收到 ProposeDelete 请求")

	err := s.Store.Propose(ctx, "DEL", req.Key, "")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "超时: %v", err)
	}
	return &proto.ProposeDeleteResponse{Success: true}, nil
}

func convertInternalToProtoChunkInfo(internal metadata.ChunkInfo) *proto.ChunkInfo {
	return &proto.ChunkInfo{
		ChunkId:       internal.ChunkId,
		ChunkLocation: internal.ChunkLocation,
		ChunkStatus:   internal.ChunkStatus,
	}
}

// --- Metastore 需要实现的方法 (示例) ---
/*
// 在你的 kvstore 或 metastore 包中:

type Metastore struct {
	metadataStore map[string]ChunkInfo // 存储实际元数据
	proposeC      chan<- string        // 用于向 Raft 节点提议
	mu            sync.RWMutex         // 保护 metadataStore 的读写
	// ... 可能还有 snapshotter 等 Raft 相关的东西
}

// ProposeC 返回用于提交提案的通道
func (m *Metastore) ProposeC() chan<- string {
	return m.proposeC
}

// Lookup 从内存中查找元数据 (需要加读锁)
func (m *Metastore) Lookup(key string) (ChunkInfo, bool) {
	m.mu.RLock() // 获取读锁
	defer m.mu.RUnlock() // 保证锁被释放
	val, ok := m.metadataStore[key]
	return val, ok
}

// ApplySet 在 Raft 状态机应用日志时调用 (需要加写锁)
func (m *Metastore) ApplySet(key string, value ChunkInfo) {
	m.mu.Lock() // 获取写锁
	defer m.mu.Unlock() // 保证锁被释放
	m.metadataStore[key] = value
}

// ApplyDelete 在 Raft 状态机应用日志时调用 (需要加写锁)
func (m *Metastore) ApplyDelete(key string) {
	m.mu.Lock() // 获取写锁
	defer m.mu.Unlock() // 保证锁被释放
	delete(m.metadataStore, key)
}

// ... 其他需要的方法，如快照相关的 LoadFromSnapshot, SaveSnapshot 等 ...
*/

// --- kvstore.ChunkInfo (示例定义) ---
/*
// 在你的 kvstore 包中:
package kvstore

// ChunkInfo 代表内部存储的元数据结构
type ChunkInfo struct {
	ChunkId       string `json:"chunkId"`       // 与 proto 对应
	ChunkLocation string `json:"chunkLocation"` // 与 proto 对应
	ChunkStatus   string `json:"chunkStatus"`   // 与 proto 对应
	// ... 其他你需要的内部字段
}
*/
