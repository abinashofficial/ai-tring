package audios

import (
	"aitring/model"
	"aitring/store"
	AudioStore "aitring/store/audiostore"
	"aitring/store/pipelinestore"
	"context"
)

func New(Store store.Store) Service {
	return &audioService{
		audioRepo: Store.AudioStore,
		audioStore: Store.PipelineStore,
	}
}

type audioService struct {
	audioRepo AudioStore.Repository
	audioStore pipelinestore.Repository
}


	func (s audioService) UploadAudio(ctx context.Context, raw model.RawChunk) (bool, error){
		return s.audioStore.Ingest(ctx, raw)
	}
	
	
	func (s audioService)GetAudioChunks(userID string) ([]model.ChunkMeta, error){
		return s.audioStore.GetChunksByUser(userID)

	}
	func (s audioService)GetAudioMetadata(chunkID string) (model.ChunkMeta, error){
		return  s.audioStore.GetMetadata(chunkID)

	}