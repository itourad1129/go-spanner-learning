package usecase

import (
	"cloud.google.com/go/spanner"
	"context"
	"go-spanner-learning/domain"
	"go-spanner-learning/domain/user"
	"time"
)

type userLoginUsecase struct {
	userLoginRepository user.UserLoginRepository
	contextTimeout      time.Duration
}

func NewUserLoginUsecase(userLoginRepository user.UserLoginRepository, userTransferRepository user.UserTransferRepository, timeout time.Duration) domain.UserLoginUsecase {
	return &userLoginUsecase{
		userLoginRepository: userLoginRepository,
		contextTimeout:      timeout,
	}
}

func (u userLoginUsecase) InsertOrUpdate(c context.Context, tx *spanner.ReadWriteTransaction, userID int64) (user.UserLogin, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.userLoginRepository.InsertOrUpdate(ctx, tx, userID)
}
