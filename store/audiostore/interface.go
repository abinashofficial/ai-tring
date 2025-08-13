package audiostore

import(
	"aitring/model"
)

type Repository interface {
	saveMetadataToDisk()
	GetChunksByUser(userID string) []model.ChunkMeta 
	SaveChunk(chunkID string, data []byte, meta model.ChunkMeta)
	 GetMetadata(chunkID string) (model.ChunkMeta, bool)

}