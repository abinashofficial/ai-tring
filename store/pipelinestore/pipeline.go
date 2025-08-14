package pipelinestore

import(
	"aitring/model"

	"context"
		"log"
		"encoding/hex"
		"time"
			"crypto/sha256"
				"math/rand"
				"fmt"
				"runtime"
				"aitring/store/audiostore"
				"errors"
)




// trySend tries to send v to ch respecting policy and context.
func trySend[T any](ctx context.Context, policy model.BackpressurePolicy, ch chan T, v T, name string) bool {
	for {
		select {
		case <-ctx.Done():
			return false
		default:
		}
		select {
		case ch <- v:
			return true
		default:
			// channel full, apply policy
			switch policy {
			case model.RejectNew:
				model.MetricsRejected.Add(1)
				log.Printf("[%s] queue full -> reject new", name)
				return false
			case model.DropOldest:
				select {
				case <-ch: // drop one
					model.MetricsDropped.Add(1)
					log.Printf("[%s] queue full -> drop oldest", name)
				default:
					// race: became available, loop and retry
				}
			}
		}
	}
}

// ===================== Pipeline Config =====================

func DefaultConfig() model.PipelineConfig {
	return model.PipelineConfig{
		IngestionWorkers:      max(2, runtime.NumCPU()/2),
		ValidationWorkers:     max(2, runtime.NumCPU()/2),
		TransformationWorkers: max(2, runtime.NumCPU()),
		MetadataWorkers:       2,
		StorageWorkers:        2,
		IngestionQueue:        128,
		ValidationQueue:       128,
		TransformationQueue:   64,
		MetadataQueue:         64,
		StorageQueue:          64,
		Policy:                model.DropOldest,
	}
}

type Pipeline struct {
	cfg model.PipelineConfig

	ingestQ    chan model.RawChunk
	validateQ  chan model.ValidatedChunk
	transformQ chan model.TransformedChunk
	metadataQ  chan model.TransformedChunk
	storageQ   chan model.TransformedChunk

	store *audiostore.AudioStore
}

func New(cfg model.PipelineConfig, store *audiostore.AudioStore) Repository {
	p := &Pipeline{cfg: cfg, store: store}
	p.ingestQ = make(chan model.RawChunk, cfg.IngestionQueue)
	p.validateQ = make(chan model.ValidatedChunk, cfg.ValidationQueue)
	p.transformQ = make(chan model.TransformedChunk, cfg.TransformationQueue)
	p.metadataQ = make(chan model.TransformedChunk, cfg.MetadataQueue)
	p.storageQ = make(chan model.TransformedChunk, cfg.StorageQueue)
	return p
}

func (p *Pipeline) Start(ctx context.Context) {
	// Validation workers
	for i := 0; i < p.cfg.ValidationWorkers; i++ {
		go func(id int) {
			for {
				select {
				case <-ctx.Done():
					return
				case raw, ok := <-p.ingestQ:
					if !ok {
						return
					}
					valid := len(raw.Data) > 0 // simple check
					if !valid {
						continue
					}
					model.MetricsValidated.Add(1)
					vc := model.ValidatedChunk{RawChunk: raw, Valid: true}
					_ = trySend(ctx, p.cfg.Policy, p.validateQ, vc, "validate")
				}
			}
		}(i)
	}

	// Transformation workers
	for i := 0; i < p.cfg.TransformationWorkers; i++ {
		go func(id int) {
			for {
				select {
				case <-ctx.Done():
					return
				case vc, ok := <-p.validateQ:
					if !ok {
						return
					}
					// simulate heavy work (checksum + mock FFT + transcript)
					chk := sha256.Sum256(vc.Data)
					t := model.TransformedChunk{ValidatedChunk: vc, Checksum: hex.EncodeToString(chk[:])}
					// mock FFT features
					feat := make([]float64, 16)
					for i := range feat {
						feat[i] = rand.Float64()
					}
					t.FFT = feat
					t.Transcript = fmt.Sprintf("fake transcript (%d bytes)", len(vc.Data))
					time.Sleep(10 * time.Millisecond) // simulate CPU work
					model.MetricsTransformed.Add(1)
					_ = trySend(ctx, p.cfg.Policy, p.transformQ, t, "transform")
				}
			}
		}(i)
	}

	// Metadata extraction workers (here it's identity pass-through but place to enrich)
	for i := 0; i < p.cfg.MetadataWorkers; i++ {
		go func(id int) {
			for {
				select {
				case <-ctx.Done():
					return
				case t, ok := <-p.transformQ:
					if !ok {
						return
					}
					// could add more enrichment here
					_ = trySend(ctx, p.cfg.Policy, p.metadataQ, t, "metadata")
				}
			}
		}(i)
	}

	// Storage workers
	for i := 0; i < p.cfg.StorageWorkers; i++ {
		go func(id int) {
			for {
				select {
				case <-ctx.Done():
					return
				case t, ok := <-p.metadataQ:
					if !ok {
						return
					}
					p.store.SaveChunk(t)
					model.MetricsStored.Add(1)
				}
			}
		}(i)
	}
}


// Ingest pushes a new raw chunk into the pipeline with backpressure handling.
func (p *Pipeline) Ingest(ctx context.Context, raw model.RawChunk) (bool, error) {
	model.MetricsIngested.Add(1)
	ok := trySend(ctx, p.cfg.Policy, p.ingestQ, raw, "ingest")
	if !ok{
		return ok, errors.New("audio ingest failed")
	}
	return ok, nil
}

func (s *Pipeline) GetChunksByUser(userID string) ([]model.ChunkMeta, error)  {
	return s.store.GetChunksByUser(userID)
}

func (s *Pipeline) GetMetadata(chunkID string) (model.ChunkMeta, error) {
	return s.store.GetMetadata(chunkID)
}