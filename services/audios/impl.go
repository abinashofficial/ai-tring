package audios

import (
	"aitring/model"
	"aitring/store"
	AudioStore "aitring/store/audiostore"
	"aitring/store/pipelinestore"
	"errors"
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


	func (s audioService) UploadAudio(ctx context.Context, raw model.RawChunk) bool{
		return s.audioStore.Ingest(ctx, raw)

	}
	
	
	func (s audioService)GetAudioChunks(userID string) ([]model.ChunkMeta, error){
		meta := s.audioStore.GetChunksByUser(userID)
	if len(meta) == 0 {
		return nil, errors.New("no audio chunks found for user")
	}
		return meta, nil

	}
	func (s audioService)GetAudioMetadata(chunkID string) (model.ChunkMeta, error){
			meta, found := s.audioStore.GetMetadata(chunkID)
	if !found {
		return model.ChunkMeta{}, errors.New("audio metadata not found")
	}
		return meta, nil

	}