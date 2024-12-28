package repository

import (
	"cloud.google.com/go/spanner"
	"context"
	"errors"
	"fmt"
	"go-spanner-learning/database"
	mytime "go-spanner-learning/domain/time"
	"go-spanner-learning/domain/user"
	"google.golang.org/api/iterator"
	"time"
)

type userLoginRepository struct {
	spannerClient *spanner.Client
	tableName     string
}

func NewUserLoginRepository(spannerClient *spanner.Client, tableName string) user.UserLoginRepository {
	return &userLoginRepository{
		spannerClient: spannerClient,
		tableName:     tableName,
	}
}

func (repo userLoginRepository) InsertOrUpdate(ctx context.Context, tx *spanner.ReadWriteTransaction, userID int64) (user.UserLogin, error) {
	var userLogin user.UserLogin
	_, columns, err := database.GetSpannerColumns(user.UserLogin{})
	if err != nil {
		return userLogin, err
	}

	selectStmt := spanner.Statement{
		SQL:    fmt.Sprintf("SELECT %s, %s FROM %s WHERE %s = $1", columns[user.LastLogin], columns[user.TotalLoginDays], repo.tableName, columns[user.UserID]),
		Params: map[string]interface{}{"p1": userID},
	}

	var lastLogin time.Time
	var totalLoginDays int64

	iter := tx.Query(ctx, selectStmt)
	defer iter.Stop()

	row, err := iter.Next()
	if row == nil {
		lastLogin = time.Time{}
		totalLoginDays = 0
	} else if err != nil && !errors.Is(err, iterator.Done) {
		return userLogin, err
	} else if err := row.Columns(&lastLogin, &totalLoginDays); err != nil {
		return userLogin, err
	}

	currentTime := mytime.Now()

	if lastLogin.IsZero() || lastLogin.Before(currentTime.Truncate(24*time.Hour)) {
		totalLoginDays += 1
	}

	mutation := spanner.InsertOrUpdate(
		repo.tableName,
		[]string{
			columns[user.UserID],
			columns[user.LastLogin],
			columns[user.TotalLoginDays],
		},
		[]interface{}{
			userID,
			mytime.CommitTimeStamp(),
			totalLoginDays,
		},
	)

	if err := tx.BufferWrite([]*spanner.Mutation{mutation}); err != nil {
		return userLogin, err
	}

	userLogin = user.UserLogin{
		UserID:         userID,
		LastLogin:      mytime.CommitTimeStamp(),
		TotalLoginDays: totalLoginDays,
	}

	return userLogin, nil
}

func (u userLoginRepository) GetUserLogin(c context.Context, tx *spanner.ReadWriteTransaction, userId int64) (user.UserLogin, error) {
	//TODO implement me
	panic("implement me")
}
