package audios

import (
	"aitring/store"
	audioStore "aitring/store/audiostore"
	"aitring/store/pipelinestore"
	audioTests "aitring/tests/test_cases/audios"
	"context"
	"fmt"
	"os"
	"testing"
	// "aitring/tests"

	"github.com/stretchr/testify/assert"
		"aitring/tests/mock_repos"
		"github.com/stretchr/testify/mock"


)

var TestAudioService Service

var mockRepos store.Store


var ctx context.Context

func TestMain(m *testing.M) {
    cfg := pipelinestore.DefaultConfig()
    audStore := audioStore.NewAudioStore("metadata.json")

    pipeline := pipelinestore.New(cfg, audStore)
    if pipeline == nil {
        panic("PipelineStore is nil – check pipelinestore.New()")
    }

    mockRepos = store.Store{
        AudioStore:    audStore,
        PipelineStore: pipeline,
    }

    TestAudioService = New(mockRepos)

    fmt.Println("Starting Audios Related Test Cases")
    code := m.Run()
    fmt.Println("Done with Audios Service Relates Test Cases")

    os.Exit(code)
}


func TestUploadAudio(t *testing.T) {
	  mockPipeline := new(mock_repos.MockPipeRepo)
    mockAudio := new(mock_repos.MockAudioRepo)

	    // Setup mock to return what UploadAudio expects
    mockPipeline.
        On("Ingest", mock.Anything, mock.Anything).
        Return(true, nil)

    TestAudioService := New(store.Store{
        AudioStore:    audioStore.Repository(mockAudio),
        PipelineStore: pipelinestore.Repository(mockPipeline),
    })
 for _, tc := range audioTests.AudioIngestTestCases {
        t.Run(tc.Name, func(t *testing.T) {
            tc.MockSetup(mockPipeline)

            output, err := TestAudioService.UploadAudio(ctx, tc.Input)

            assert.Equal(t, tc.ExpectedOutput, output)
            assert.Equal(t, tc.ExpectedErr, err)
        })
    }
			mockPipeline.AssertExpectations(t)

}




func TestGetAudioChunks(t *testing.T) {
	  mockPipeline := new(mock_repos.MockPipeRepo)
    mockAudio := new(mock_repos.MockAudioRepo)

    TestAudioService := New(store.Store{
        AudioStore:    audioStore.Repository(mockAudio),
        PipelineStore: pipelinestore.Repository(mockPipeline),
    })
		for _, tc := range audioTests.AudioGetChunksByUserTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			 tc.MockSetup(mockPipeline)
			_, err := TestAudioService.GetAudioChunks(tc.Input)
			if tc.ExpectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tc.ExpectedErr.Error(), err.Error())
			}
			mockPipeline.AssertExpectations(t)
		})
	}
}

func TestGetAudioMetadata(t *testing.T) {
	  mockPipeline := new(mock_repos.MockPipeRepo)
    mockAudio := new(mock_repos.MockAudioRepo)

    TestAudioService := New(store.Store{
        AudioStore:    audioStore.Repository(mockAudio),
        PipelineStore: pipelinestore.Repository(mockPipeline),
    })
 for _, tc := range audioTests.AudioGetMetadataTestCases {
        t.Run(tc.Name, func(t *testing.T) {
            tc.MockSetup(mockPipeline)

            output, err := TestAudioService.GetAudioMetadata(tc.Input)

            assert.Equal(t, tc.ExpectedOutput, output)
            assert.Equal(t, tc.ExpectedErr, err)
        })
    }
			mockPipeline.AssertExpectations(t)

}
