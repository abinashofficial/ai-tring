package model


import (
	"time"
		"sync/atomic"

)

type ChunkMeta struct {
	ChunkID    string    `json:"chunk_id"`
	SessionID  string    `json:"session_id"`
	UserID     string    `json:"user_id"`
	Timestamp  time.Time `json:"timestamp"`
	Size       int       `json:"size_bytes"`
	Checksum   string    `json:"checksum"`
	Transcript string    `json:"transcript,omitempty"`
}

type RawChunk struct {
	ChunkID   string
	SessionID string
	UserID    string
	Data      []byte
	Received  time.Time
	    	AckCh     chan ChunkMeta // optional: notify caller when processed

}

type ValidatedChunk struct {
	RawChunk
	Valid bool
}

type TransformedChunk struct {
	ValidatedChunk
	Checksum   string
	FFT        []float64 // mocked features
	Transcript string    // mocked transcript
}

type BackpressurePolicy int

const (
	RejectNew  BackpressurePolicy = iota // if queue full: reject new item
	DropOldest                           // if queue full: drop oldest to make room
)

type PipelineConfig struct {
	IngestionWorkers      int
	ValidationWorkers     int
	TransformationWorkers int
	MetadataWorkers       int
	StorageWorkers        int

	IngestionQueue      int
	ValidationQueue     int
	TransformationQueue int
	MetadataQueue       int
	StorageQueue        int

	Policy BackpressurePolicy
}

// ===================== Metrics =====================

type counter uint64

func (c *counter) Add(n uint64) { atomic.AddUint64((*uint64)(c), n) }
func (c *counter) Get() uint64  { return atomic.LoadUint64((*uint64)(c)) }

var (
	MetricsIngested    counter
	MetricsValidated   counter
	MetricsTransformed counter
	MetricsStored      counter
	MetricsRejected    counter
	MetricsDropped     counter
)