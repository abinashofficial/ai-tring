package audios

import(
	"aitring/model"
)

type Service interface {
	// AudioStore is the interface for audio management	
	UploadAudio(chunkID string, data []byte, meta model.ChunkMeta) error
	GetAudioChunks(userID string) ([]model.ChunkMeta, error)
	GetAudioMetadata(chunkID string) (model.ChunkMeta, error)
}