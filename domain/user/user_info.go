package user

import (
	"cloud.google.com/go/spanner"
	"context"
	"go-spanner-learning/domain/time"
)

const (
	UserID = "UserID"
	Name   = "Name"
)

type UserInfo struct {
	UserID int64  `spanner:"user_id"`
	Name   string `spanner:"name"`
	time.RecordTime
}

type UserInfoRepository interface {
	Create(c context.Context, tx *spanner.ReadWriteTransaction, userName string) (int64, error)
	Fetch(c context.Context) ([]UserInfo, error)
	GetUserID(c context.Context, userId string) (UserInfo, error)
	GetUserName(c context.Context, tx *spanner.ReadWriteTransaction, id string) (UserInfo, error)
}
