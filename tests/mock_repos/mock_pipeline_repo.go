package mock_repos

import(
		"aitring/model"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockPipeRepo struct {
	mock.Mock
}



func (m *MockPipeRepo)Start(ctx context.Context) {
	m.Called(ctx)
}

func (m *MockPipeRepo) Ingest(ctx context.Context, raw model.RawChunk) (bool, error) {
	args := m.Called(ctx, raw)
	if val, ok := args.Get(0).(bool); ok {
		return val, args.Error(1)
	}
	return false, args.Error(1)
}

func (m *MockPipeRepo) GetMetadata(chunkID string) (model.ChunkMeta, error) {
	args := m.Called(chunkID)
	if val, ok := args.Get(0).(model.ChunkMeta); ok {
		return val, args.Error(1)
	}
	return model.ChunkMeta{}, args.Error(1)
}

func (m *MockPipeRepo) GetChunksByUser(userID string) ([]model.ChunkMeta, error) {
	args := m.Called(userID)
	if val, ok := args.Get(0).([]model.ChunkMeta); ok {
		return val, args.Error(1)
	}
	return nil, args.Error(1)
}
