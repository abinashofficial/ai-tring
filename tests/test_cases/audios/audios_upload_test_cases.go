package audios

import (
	"aitring/model"
	"aitring/tests"
	"errors"
	"time"

	"github.com/stretchr/testify/mock"
	"aitring/tests/mock_repos"

)




type AudioIngestTestCase struct 
{
        Name           string
        Input          model.RawChunk
        MockSetup       func(mockPipe *mock_repos.MockPipeRepo)
        ExpectedOutput bool
        ExpectedErr    error
    }



var AudioIngestTestCases = []AudioIngestTestCase{
    {
        Name: "Valid Upload",
        Input: model.RawChunk{
            ChunkID: "123",
            UserID:  "u1",
            Data:    []byte("data"),
        },
        MockSetup: func(mockPipe *mock_repos.MockPipeRepo) {
            tests.MockPieRepo.On("Ingest", mock.Anything, model.RawChunk{
            ChunkID: "123",
            UserID:  "u1",
            Data:    []byte("data"),
        }).
                Return(true, nil)
        },
        ExpectedOutput: true,
        ExpectedErr:    nil,
    },
}








var chunk = model.ChunkMeta{
    
        ChunkID:    "test-chunk-id",
        UserID:     "test-user-id",
        SessionID:  "test-session-id",
        Checksum:   "test-checksum",
        Transcript: "test transcript",
        Timestamp:  time.Date(2025, 8, 13, 12, 0, 0, 0, time.UTC),
        Size:       98767,
    }

var chunks = []model.ChunkMeta{
    {
        ChunkID:    "test-chunk-id",
        UserID:     "test-user-id",
        SessionID:  "test-session-id",
        Checksum:   "test-checksum",
        Transcript: "test transcript",
        Timestamp:  time.Date(2025, 8, 13, 12, 0, 0, 0, time.UTC),
        Size:       98767,
    },
}

type AudioGetChunksByUserTestCase struct {
    Name           string
    Input          string
    MockSetup       func(mockPipe *mock_repos.MockPipeRepo)
    ExpectedOutput []model.ChunkMeta
    ExpectedErr    error
}
var AudioGetChunksByUserTestCases = []AudioGetChunksByUserTestCase{
    {
        Name:  "Valid Case",
        Input: "42",
        MockSetup: func(mockPipe *mock_repos.MockPipeRepo) {
            mockPipe.On("GetChunksByUser", "42").Return(chunks, nil)
        },
        ExpectedOutput: chunks,
        ExpectedErr:    nil,
    },
    {
        Name:  "Error Scenario",
        Input: "",
        MockSetup: func(mockPipe *mock_repos.MockPipeRepo) {
            mockPipe.On("GetChunksByUser", "").Return([]model.ChunkMeta{}, errors.New("no audio chunks found for user"))
        },
        ExpectedOutput: []model.ChunkMeta{},
        ExpectedErr:    errors.New("no audio chunks found for user"),
    },
}



type AudioGetMetadataTestCase struct {
    Name           string
    Input          string
    MockSetup      func(mockPipe *mock_repos.MockPipeRepo)
    ExpectedOutput model.ChunkMeta
    ExpectedErr    error
}


var AudioGetMetadataTestCases = []AudioGetMetadataTestCase{
    {
        Name:  "Valid Case",
        Input: "42",
        MockSetup: func(mockPipe *mock_repos.MockPipeRepo) {
            mockPipe.
                On("GetMetadata", "42").
                Return(chunk, nil)
        },
        ExpectedOutput: chunk,
        ExpectedErr:    nil,
    },
    {
        Name:  "Error Scenario",
        Input: "",
        MockSetup: func(mockPipe *mock_repos.MockPipeRepo) {
            mockPipe.
                On("GetMetadata", "").
                Return([]model.ChunkMeta{}, errors.New("not found"))
        },
        ExpectedOutput: model.ChunkMeta{},
        ExpectedErr:    errors.New("not found"),
    },
}



