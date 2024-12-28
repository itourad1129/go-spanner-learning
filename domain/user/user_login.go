package user

import (
	"cloud.google.com/go/spanner"
	"context"
	"time"
)

const (
	LastLogin      = "LastLogin"
	TotalLoginDays = "TotalLoginDays"
)

type UserLogin struct {
	UserID         int64     `spanner:"user_id"`
	LastLogin      time.Time `spanner:"last_login"`
	TotalLoginDays int64     `spanner:"total_login_days"`
}

type UserLoginRepository interface {
	InsertOrUpdate(c context.Context, tx *spanner.ReadWriteTransaction, userId int64) (UserLogin, error)
	GetUserLogin(c context.Context, tx *spanner.ReadWriteTransaction, userId int64) (UserLogin, error)
}
