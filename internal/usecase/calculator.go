package usecase

import (
	"context"

	"github.com/Danil-Zaichkin/rest-api-stream-handler/internal/entity"
	opt "github.com/moznion/go-optional"
)

type DBRepository interface {
	SaveStreamsContexts(ctx context.Context, streamsCtx map[string]*entity.StreamContext) error
	GetStreamContext(ctx context.Context, streamID string) (opt.Option[*entity.StreamContext], error)
}

type InMemoryRepo interface {
	GetStreamContext(streamID string) opt.Option[*entity.StreamContext]
	SaveStreamContext(streamID string, streamCtx *entity.StreamContext)
	InitAndGetStreamContext(streamID string) *entity.StreamContext
}

type CalculatorUsecase struct {
	dbRepo     DBRepository
	memoryRepo InMemoryRepo
}

func NewCalculatorUsecase(dbRepo DBRepository, memoryRepo InMemoryRepo) *CalculatorUsecase {
	return &CalculatorUsecase{
		dbRepo:     dbRepo,
		memoryRepo: memoryRepo,
	}
}

func (cu *CalculatorUsecase) ApplyOperation(ctx context.Context, p entity.Package) (int, error) {
	streamCtxOpt := cu.memoryRepo.GetStreamContext(p.StreamID)

	streamCtx, err := streamCtxOpt.Take()
	if err != nil {
		streamCtxOpt, err = cu.dbRepo.GetStreamContext(ctx, p.StreamID)
		if err != nil {
			return 0, err
		}

		streamCtx, err = streamCtxOpt.Take()
		if err != nil {
			streamCtx = cu.memoryRepo.InitAndGetStreamContext(p.StreamID)
		}
	}

	err = streamCtx.ApplyOperation(p)
	if err != nil {
		return 0, err
	}

	cu.memoryRepo.SaveStreamContext(p.StreamID, streamCtx)

	return streamCtx.Value, nil
}
