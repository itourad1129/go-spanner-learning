package domain

import (
	"cloud.google.com/go/spanner"
	"context"
	"go-spanner-learning/domain/user"
)

type UserLoginRequest struct {
	UserID       int64  `form:"userID" json:"userID" binding:"required"`
	TransferCode string `form:"transferCode" json:"transferCode" binding:"required"`
}

type UserLoginResponse struct {
	UserID         int64 `json:"userID"`
	TotalLoginDays int64 `json:"totalLoginDays"`
}

type UserLoginUsecase interface {
	InsertOrUpdate(c context.Context, tx *spanner.ReadWriteTransaction, userId int64) (user.UserLogin, error)
}
