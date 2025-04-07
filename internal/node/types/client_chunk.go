package types

type DataChunk struct {
	FileName   string
	ChunkTotal int32
	ChunkId    int32
	Payload    []byte
	FileSize   int64
}
