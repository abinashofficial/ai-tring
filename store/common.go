package store

import (
	"aitring/store/audiostore"
)

type Store struct {
	AudioStore  audiostore.Repository
}