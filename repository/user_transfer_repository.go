package repository

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go-spanner-learning/database"
	"go-spanner-learning/domain/time"
	"go-spanner-learning/domain/user"
)

type userTransferRepository struct {
	spannerClient *spanner.Client
	tableName     string
}

func NewUserTransferRepository(spannerClient *spanner.Client, tableName string) user.UserTransferRepository {
	return &userTransferRepository{
		spannerClient: spannerClient,
		tableName:     tableName,
	}
}

func (repo *userTransferRepository) GetUserTransfer(c context.Context, tx *spanner.ReadWriteTransaction, userID string) (user.UserTransfer, error) {
	var userTransfer user.UserTransfer
	columnNames, columns, err := database.GetSpannerColumns(user.UserTransfer{})
	if err != nil {
		return userTransfer, err
	}

	stmt := spanner.Statement{
		SQL:    fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1", columnNames, repo.tableName, columns[user.UserID]),
		Params: map[string]interface{}{"p1": userID},
	}

	iter := tx.Query(c, stmt)
	defer iter.Stop()

	if err := iter.Do(func(row *spanner.Row) error {
		return row.ToStruct(&userTransfer)
	}); err != nil {
		return userTransfer, err
	}
	return userTransfer, nil
}

func (repo *userTransferRepository) GetTransferCode(c context.Context, transferCode string) (user.UserTransfer, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *userTransferRepository) Create(ctx context.Context, tx *spanner.ReadWriteTransaction, userID int64) (string, error) {

	_, columns, err := database.GetSpannerColumns(user.UserTransfer{})
	if err != nil {
		return "", err
	}

	transferCode, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	// 挿入用のミューテーションを作成
	mutation := spanner.Insert(repo.tableName, []string{columns[user.UserID], columns[user.TransferCode], columns[time.CreateAt]},
		[]interface{}{userID, transferCode.String(), time.CommitTimeStamp()})

	// トランザクション内で挿入を行う
	if err := tx.BufferWrite([]*spanner.Mutation{mutation}); err != nil {
		return "", err
	}

	return transferCode.String(), nil
}

func (repo *userTransferRepository) Authenticate(ctx context.Context, userID int64, transferCode string) (user.UserTransfer, error) {
	var userTransfer user.UserTransfer
	columnNames, columns, err := database.GetSpannerColumns(user.UserTransfer{})
	if err != nil {
		return userTransfer, err
	}

	stmt := spanner.Statement{
		SQL:    fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1 AND %s = $2", columnNames, repo.tableName, columns[user.UserID], columns[user.TransferCode]),
		Params: map[string]interface{}{"p1": userID, "p2": transferCode},
	}

	iter := repo.spannerClient.Single().Query(ctx, stmt)
	defer iter.Stop()

	if err := iter.Do(func(row *spanner.Row) error {
		return row.ToStruct(&userTransfer)
	}); err != nil {
		return userTransfer, err
	}
	return userTransfer, nil
}
