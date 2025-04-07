package metadata

import "encoding/json"

type ChunkInfo struct {
	ChunkId       string `json:"chunkId"`
	ChunkLocation string `json:"chunkLocation"`
	ChunkStatus   string `json:"chunkStatus"`
}

// 序列化函数：将 ChunkInfo 转换为 JSON 字符串
func serializeChunkInfo(c ChunkInfo) (string, error) {
	// 使用 json.Marshal 将结构体转换为 JSON 字符串
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// 反序列化函数：将 JSON 字符串转换为 ChunkInfo 结构体
func deserializeChunkInfo(jsonData string) (ChunkInfo, error) {
	var c ChunkInfo
	// 使用 json.Unmarshal 将 JSON 字符串转换为结构体
	err := json.Unmarshal([]byte(jsonData), &c)
	if err != nil {
		return ChunkInfo{}, err
	}
	return c, nil
}
