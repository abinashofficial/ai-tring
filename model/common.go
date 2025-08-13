package model


import (
	"time"
)

type ChunkMeta struct {
	ChunkID   string    `json:"chunk_id"`
	SessionID string    `json:"session_id"`
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
	Size      int       `json:"size_bytes"`
	Transcript string   `json:"transcript,omitempty"` // optional processed transcript
}

