package main

import (
	"context"
	"sync"
	"time"
)

func NewRepo() *RepositoryMap {
	return &RepositoryMap{
		db: make(map[string]*StreamContext),
	}
}

type RepositoryMap struct {
	db map[string]*StreamContext
	mu sync.RWMutex
}

func (r *RepositoryMap) SaveOperation(_ context.Context, p Package) (int, error) {

	r.mu.RLock()
	streamCtx, ok := r.db[p.StreamID]
	r.mu.RUnlock()

	if !ok {
		r.mu.Lock()

		streamCtx, ok = r.db[p.StreamID]
		if !ok {
			streamCtx = &StreamContext{
				StreamID: p.StreamID,
			}

			r.db[p.StreamID] = streamCtx
		}

		r.mu.Unlock()
	}

	streamCtx.mu.Lock()
	defer streamCtx.mu.Unlock()

	// имитируем сложные вычисления
	delay := getDelayByOperation(p)
	time.Sleep(delay)

	streamCtx, err := applyOpertaion(p, streamCtx)
	if err != nil {
		return 0, err
	}

	r.mu.Lock()
	r.db[p.StreamID] = streamCtx
	r.mu.Unlock()

	return streamCtx.Value, nil
}

func applyOpertaion(p Package, streamCtx *StreamContext) (*StreamContext, error) {
	switch p.Op {
	case Sum:
		streamCtx.Value += p.Value
	case Sub:
		streamCtx.Value -= p.Value
	case Mul:
		streamCtx.Value *= p.Value
	case Div:
		if p.Value == 0 {
			return nil, ErrDivByZero
		}
	}

	return streamCtx, nil
}

func getDelayByOperation(p Package) time.Duration {
	switch p.Op {
	case Sum:
		return time.Second / 4
	case Sub:
		return time.Second / 4
	case Mul:
		return time.Second / 2
	case Div:
		return time.Second
	default:
		return 0
	}
}
