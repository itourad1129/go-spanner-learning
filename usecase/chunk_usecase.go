package usecase

import (
	"context"
	"go-spanner-learning/domain"
	"go-spanner-learning/domain/master"
	"time"
)

type chunkUsecase struct {
	chunkVersionRepository master.ChunkVersionRepository
	contextTimeout         time.Duration
}

func NewChunkUsecase(chunkVersionRepository master.ChunkVersionRepository, timeout time.Duration) domain.ChunkUsecase {
	return &chunkUsecase{
		chunkVersionRepository: chunkVersionRepository,
		contextTimeout:         timeout,
	}
}

func (u chunkUsecase) GetChunkVersion(c context.Context, platformType int64) (master.ChunkVersion, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.chunkVersionRepository.GetChunkVersion(ctx, platformType)
}
