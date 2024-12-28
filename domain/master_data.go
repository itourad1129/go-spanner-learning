package domain

import (
	"context"
	"go-spanner-learning/domain/master"
)

type MasterDataVersionUsecase interface {
	GetMasterDataVersion(c context.Context) ([]master.MasterDataVersion, error)
}

type GetMasterDataVersionResponse struct {
	MasterDataID int64 `json:"masterDataid"`
	Version      int64 `json:"version"`
	ChunkID      int64 `json:"chunkID"`
}
