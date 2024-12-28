package master

import "context"

type MasterDataVersion struct {
	MasterDataID int64 `spanner:"master_data_id"`
	Version      int64 `spanner:"version"`
	ChunkID      int64 `spanner:"chunk_id"`
}

type MasterDataVersionRepository interface {
	GetMasterDataVersion(c context.Context) ([]MasterDataVersion, error)
}
