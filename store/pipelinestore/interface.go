package pipelinestore

import(
	"aitring/model"
	"context"
	
)

type Repository interface {
	Ingest(ctx context.Context, raw model.RawChunk) (bool, error)
	 Start(ctx context.Context)
	 GetMetadata(chunkID string) (model.ChunkMeta, error)
	 GetChunksByUser(userID string) ([]model.ChunkMeta, error) 
}