package audiostore

import (

"aitring/model"
"sync"
"os"
"encoding/json"
"fmt"
"path/filepath"
"errors"
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

func (s *AudioStore) SaveChunk(t model.TransformedChunk) {
	s.mu.Lock()
	s.chunks[t.ChunkID] = t.Data
	m := model.ChunkMeta{ChunkID: t.ChunkID, SessionID: t.SessionID, UserID: t.UserID, Timestamp: t.Received, Size: len(t.Data), Checksum: t.Checksum, Transcript: t.Transcript}
	s.metadata[t.ChunkID] = m
	s.mu.Unlock()
	s.SaveMetadataToDisk()
}

func (s *AudioStore) GetMetadata(chunkID string) (model.ChunkMeta, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	meta, exists := s.metadata[chunkID]
		if !exists {
		return model.ChunkMeta{}, errors.New("audio metadata not found")
	}
	return meta, nil
}

func (s *AudioStore) GetChunksByUser(userID string) ([]model.ChunkMeta, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]model.ChunkMeta, 0)
	for _, m := range s.metadata {
		if m.UserID == userID {
			out = append(out, m)
		}
	}
		if len(out) == 0 {
		return nil, errors.New("no audio chunks found for user")
	}
	return out, nil
}

func (s *AudioStore) SaveMetadataToDisk() {
	// Ensure Data directory exists
	dataDir := "Data"
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	// Prepare the file path inside Data folder
	filePath := filepath.Join(dataDir, s.metadataFile)

	s.mu.RLock()
	data, err := json.MarshalIndent(s.metadata, "", "  ")
	s.mu.RUnlock()
	if err != nil {
		fmt.Println("Error marshaling metadata:", err)
		return
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		fmt.Println("Error writing metadata file:", err)
	}
}

