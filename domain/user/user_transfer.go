package user

import (
	"cloud.google.com/go/spanner"
	"context"
	"go-spanner-learning/domain/time"
)

const (
	TransferCode = "TransferCode"
)

type UserTransfer struct {
	UserID       int64  `spanner:"user_id"`
	TransferCode string `spanner:"transfer_code"`
	time.RecordCreateTime
}

type UserTransferRepository interface {
	Create(c context.Context, tx *spanner.ReadWriteTransaction, userID int64) (string, error)
	GetUserTransfer(c context.Context, tx *spanner.ReadWriteTransaction, userID string) (UserTransfer, error)
	GetTransferCode(c context.Context, transferCode string) (UserTransfer, error)
	Authenticate(ctx context.Context, userID int64, code string) (UserTransfer, error)
}
