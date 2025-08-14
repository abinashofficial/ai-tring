package audiostore

import(
	"aitring/model"
)

type Repository interface {
	SaveMetadataToDisk()
	GetChunksByUser(userID string) ([]model.ChunkMeta, error) 
	SaveChunk(t model.TransformedChunk)
	 GetMetadata(chunkID string) (model.ChunkMeta, error)
	 
}