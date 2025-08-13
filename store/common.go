package store

import (
	"aitring/store/audiostore"
	"aitring/store/pipelinestore"
)

type Store struct {
	AudioStore  audiostore.Repository
	PipelineStore pipelinestore.Repository
}