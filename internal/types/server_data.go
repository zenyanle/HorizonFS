package types

type SenderData struct {
	FileName   string
	ChunkTotal int32
	ChunkId    int32
	Payload    []byte
}
