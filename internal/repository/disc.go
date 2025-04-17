package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Danil-Zaichkin/rest-api-stream-handler/internal/entity"
	opt "github.com/moznion/go-optional"
	"github.com/redis/go-redis/v9"
)

type DBRepository struct {
	rdb *redis.Client
}

func NewDBRepo(client *redis.Client) *DBRepository {
	return &DBRepository{
		rdb: client,
	}
}

func (r *DBRepository) SaveStreamsContexts(ctx context.Context, streamsCtx map[string]*entity.StreamContext) error {
	for streamID, streamCtx := range streamsCtx {
		streamCtxJSON, err := json.Marshal(streamCtx)
		if err != nil {
			return fmt.Errorf("can't masrshal stream ctx: %w", err)
		}

		err = r.rdb.Set(
			ctx,
			fmt.Sprintf("stream:%s", streamID),
			string(streamCtxJSON),
			time.Hour,
		).Err()
		if err != nil {
			return fmt.Errorf("can't save ctx in redis: %w", err)
		}
	}

	return nil
}

func (r *DBRepository) GetStreamContext(ctx context.Context, streamID string) (opt.Option[*entity.StreamContext], error) {
	val, err := r.rdb.Get(ctx, fmt.Sprintf("stream:%s", streamID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("can't get ctx from redis: %w", err)
	}

	var streamCtx entity.StreamContext
	err = json.Unmarshal([]byte(val), &streamCtx)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal ctx value: %w", err)
	}

	return opt.Some(&streamCtx), nil
}
