package audios

import (
	"aitring/model"
	"aitring/store"
	AudioStore "aitring/store/audiostore"
	"errors"
)

func New(Store store.Store) Service {
	return &audioService{
		audioRepo: Store.AudioStore,
	}
}

type audioService struct {
	audioRepo AudioStore.Repository
}


	func (s audioService) UploadAudio(chunkID string, data []byte, meta model.ChunkMeta) error{
		 s.audioRepo.SaveChunk(chunkID, data, meta)

		return nil

	}
	
	
	func (s audioService)GetAudioChunks(userID string) ([]model.ChunkMeta, error){
		meta := s.audioRepo.GetChunksByUser(userID)
	if len(meta) == 0 {
		return nil, errors.New("no audio chunks found for user")
	}
		return meta, nil

	}
	func (s audioService)GetAudioMetadata(chunkID string) (model.ChunkMeta, error){
			meta, found := s.audioRepo.GetMetadata(chunkID)
	if !found {
		return model.ChunkMeta{}, errors.New("audio metadata not found")
	}
		return meta, nil

	}