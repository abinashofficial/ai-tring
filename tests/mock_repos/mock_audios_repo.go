package mock_repos

import (
	"aitring/model"

	"github.com/stretchr/testify/mock"
)

type MockAudioRepo struct {
	mock.Mock
}


	 func (m *MockAudioRepo) SaveMetadataToDisk(){
m.Called()
}

	 func (m *MockAudioRepo) SaveChunk(t model.TransformedChunk){
m.Called(t)
}

func (m *MockAudioRepo) GetMetadata(chunkID string) (model.ChunkMeta, error) {
	args := m.Called(chunkID)
	if val, ok := args.Get(0).(model.ChunkMeta); ok {
		return val, args.Error(1)
	}
	return model.ChunkMeta{}, args.Error(1)
}

func (m *MockAudioRepo) GetChunksByUser(userID string) ([]model.ChunkMeta, error) {
	args := m.Called(userID)
	if val, ok := args.Get(0).([]model.ChunkMeta); ok {
		return val, args.Error(1)
	}
	return nil, args.Error(1)
}
