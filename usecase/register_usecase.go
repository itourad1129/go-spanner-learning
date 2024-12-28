package usecase

import (
	"cloud.google.com/go/spanner"
	"context"
	"go-spanner-learning/domain"
	"go-spanner-learning/domain/user"
	"time"
)

type userRegisterUsecase struct {
	userInfoRepository     user.UserInfoRepository
	userTransferRepository user.UserTransferRepository
	contextTimeout         time.Duration
}

func NewUserRegisterUsecase(userInfoRepository user.UserInfoRepository, userTransferRepository user.UserTransferRepository, timeout time.Duration) domain.UserRegisterUsecase {
	return &userRegisterUsecase{
		userInfoRepository:     userInfoRepository,
		userTransferRepository: userTransferRepository,
		contextTimeout:         timeout,
	}
}

func (u *userRegisterUsecase) CreateUserInfo(c context.Context, tx *spanner.ReadWriteTransaction, userName string) (int64, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.userInfoRepository.Create(ctx, tx, userName)
}

func (u *userRegisterUsecase) CreateUserTransfer(c context.Context, tx *spanner.ReadWriteTransaction, userID int64) (string, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.userTransferRepository.Create(ctx, tx, userID)
}

func (u *userRegisterUsecase) GetUserByUserName(c context.Context, tx *spanner.ReadWriteTransaction, email string) (user.UserInfo, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.userInfoRepository.GetUserName(ctx, tx, email)
}
