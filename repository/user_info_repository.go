package repository

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"go-spanner-learning/database"
	"go-spanner-learning/domain/time"
	"go-spanner-learning/domain/user"
)

type userInfoRepository struct {
	spannerClient *spanner.Client
	tableName     string
}

func NewUserInfoRepository(spannerClient *spanner.Client, tableName string) user.UserInfoRepository {
	return &userInfoRepository{
		spannerClient: spannerClient,
		tableName:     tableName,
	}
}

func (repo *userInfoRepository) Create(ctx context.Context, tx *spanner.ReadWriteTransaction, userName string) (int64, error) {

	_, columns, err := database.GetSpannerColumns(user.UserInfo{})
	if err != nil {
		return 0, err
	}

	stmt := spanner.Statement{SQL: "SELECT nextval('user_id_sequence')"}
	iter := tx.Query(ctx, stmt)
	defer iter.Stop()

	var userID int64
	row, err := iter.Next()
	if err != nil {
		return 0, err
	}
	if err := row.Columns(&userID); err != nil {
		return 0, err
	}

	// 挿入用のミューテーションを作成
	mutation := spanner.Insert(repo.tableName, []string{columns[user.UserID], columns[user.Name], columns[time.CreateAt], columns[time.UpdateAt]},
		[]interface{}{userID, userName, time.CommitTimeStamp(), time.CommitTimeStamp()})

	// トランザクション内で挿入を行う
	if err := tx.BufferWrite([]*spanner.Mutation{mutation}); err != nil {
		return 0, err
	}
	return userID, nil
}

func (repo *userInfoRepository) Fetch(c context.Context) ([]user.UserInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *userInfoRepository) GetUserID(c context.Context, userID string) (user.UserInfo, error) {

	var userInfo user.UserInfo
	columnNames, columns, err := database.GetSpannerColumns(user.UserInfo{})
	if err != nil {
		return userInfo, err
	}

	stmt := spanner.Statement{
		SQL:    fmt.Sprintf("SELECT %s FROM %s WHERE %s = @userID", columnNames, repo.tableName, columns[user.UserID]),
		Params: map[string]interface{}{"userID": userID},
	}

	iter := repo.spannerClient.Single().Query(c, stmt)
	defer iter.Stop()

	if err := iter.Do(func(row *spanner.Row) error {
		return row.ToStruct(&userInfo)
	}); err != nil {
		return userInfo, err
	}
	return userInfo, nil
}

func (repo *userInfoRepository) GetUserName(c context.Context, tx *spanner.ReadWriteTransaction, userName string) (user.UserInfo, error) {
	var userInfo user.UserInfo
	columnNames, columns, err := database.GetSpannerColumns(user.UserInfo{})
	if err != nil {
		return userInfo, err
	}

	stmt := spanner.Statement{
		SQL:    fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1", columnNames, repo.tableName, columns[user.Name]),
		Params: map[string]interface{}{"p1": userName},
	}

	iter := tx.Query(c, stmt)
	defer iter.Stop()

	if err := iter.Do(func(row *spanner.Row) error {
		return row.ToStruct(&userInfo)
	}); err != nil {
		return userInfo, err
	}
	return userInfo, nil
}
