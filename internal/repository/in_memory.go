package repository

import (
	"maps"
	"sync"

	"github.com/Danil-Zaichkin/rest-api-stream-handler/internal/entity"
	opt "github.com/moznion/go-optional"
)

type InMemoryRepository struct {
	mu               sync.RWMutex
	calculateResults map[string]*entity.StreamContext
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		calculateResults: map[string]*entity.StreamContext{},
	}
}

func (imr *InMemoryRepository) GetStreamContext(streamID string) opt.Option[*entity.StreamContext] {
	imr.mu.RLock()
	defer imr.mu.RUnlock()
	streamCtx, ok := imr.calculateResults[streamID]
	if !ok {
		return nil
	}

	return opt.Some(streamCtx)
}

func (imr *InMemoryRepository) GetStreamsContexts() map[string]*entity.StreamContext {
	imr.mu.Lock()
	defer imr.mu.Unlock()

	return maps.Clone(imr.calculateResults)
}

func (imr *InMemoryRepository) InitAndGetStreamContext(streamID string) *entity.StreamContext {
	imr.mu.RLock()
	streamCtx, ok := imr.calculateResults[streamID]
	imr.mu.RUnlock()

	if !ok {
		imr.mu.Lock()

		streamCtx, ok = imr.calculateResults[streamID]
		if !ok {
			streamCtx = &entity.StreamContext{
				StreamID: streamID,
			}

			imr.calculateResults[streamID] = streamCtx
		}

		imr.mu.Unlock()
	}

	return streamCtx
}

func (imr *InMemoryRepository) SaveStreamContext(streamID string, streamCtx *entity.StreamContext) {
	imr.mu.Lock()
	imr.calculateResults[streamID] = streamCtx
	imr.mu.Unlock()
}
