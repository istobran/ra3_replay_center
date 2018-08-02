package utils

type BodyChunk struct {
	TimeCode  uint32
	ChunkType byte
	ChunkSize uint32
	Data      []byte
	Zero      uint32
}

type ReplayBody struct {
	Size   int         // Chunk count
	Chunks []BodyChunk // Chunk list
}
