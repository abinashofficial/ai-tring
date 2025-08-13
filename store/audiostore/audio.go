package audiostore

import (

"aitring/model"
"sync"
"os"
"encoding/json"
)


func NewAudioStore(metadataFile string) *AudioStore {
	as := &AudioStore{
		chunks:       make(map[string][]byte),
		metadata:     make(map[string]model.ChunkMeta),
		metadataFile: metadataFile,
	}
	return as
}

func New() Repository {
	return NewAudioStore("metadata.json")
}	


type AudioStore struct {
	mu           sync.RWMutex
	chunks       map[string][]byte
	metadata     map[string]model.ChunkMeta
	metadataFile string
}

func (s *AudioStore) SaveChunk(chunkID string, data []byte, meta model.ChunkMeta) {
	s.mu.Lock()
	s.chunks[chunkID] = data
	s.metadata[chunkID] = meta
	s.mu.Unlock()
	s.saveMetadataToDisk()
}

func (s *AudioStore) GetMetadata(chunkID string) (model.ChunkMeta, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	meta, exists := s.metadata[chunkID]
	return meta, exists
}

func (s *AudioStore) GetChunksByUser(userID string) []model.ChunkMeta {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var results []model.ChunkMeta
	for _, m := range s.metadata {
		if m.UserID == userID {
			results = append(results, m)
		}
	}
	return results
}

func (s *AudioStore) saveMetadataToDisk() {
	s.mu.RLock()
	data, _ := json.MarshalIndent(s.metadata, "", "  ")
	s.mu.RUnlock()
	_ = os.WriteFile(s.metadataFile, data, 0644)
}