package audiostore

import(
	"aitring/model"
)

type Repository interface {
	saveMetadataToDisk()
	GetChunksByUser(userID string) []model.ChunkMeta 
	SaveChunk(t model.TransformedChunk)
	 GetMetadata(chunkID string) (model.ChunkMeta, bool)
	 
}