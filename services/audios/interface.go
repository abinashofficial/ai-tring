package audios

import (
	"aitring/model"
	"context"
)

type Service interface {
	// AudioStore is the interface for audio management	
	UploadAudio(ctx context.Context, raw model.RawChunk) (bool, error)
	GetAudioChunks(userID string) ([]model.ChunkMeta, error)
	GetAudioMetadata(chunkID string) (model.ChunkMeta, error)
}