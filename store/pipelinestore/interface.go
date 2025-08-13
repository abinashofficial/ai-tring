package pipelinestore

import(
	"aitring/model"
	"context"
	
)

type Repository interface {
	Ingest(ctx context.Context, raw model.RawChunk) bool
	 Start(ctx context.Context)
	 GetMetadata(chunkID string) (model.ChunkMeta, bool)
	 GetChunksByUser(userID string) []model.ChunkMeta
}