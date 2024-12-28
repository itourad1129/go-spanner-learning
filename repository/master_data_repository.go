package repository

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"go-spanner-learning/database"
	"go-spanner-learning/domain/master"
)

type masterDataVersionRepository struct {
	spannerClient *spanner.Client
	tableName     string
}

func NewMasterDataVersionRepository(spannerClient *spanner.Client, tableName string) master.MasterDataVersionRepository {
	return masterDataVersionRepository{
		spannerClient: spannerClient,
		tableName:     tableName,
	}
}

func (repo masterDataVersionRepository) GetMasterDataVersion(c context.Context) ([]master.MasterDataVersion, error) {
	var masterDataVersions []master.MasterDataVersion
	columnNames, _, err := database.GetSpannerColumns(master.MasterDataVersion{})
	if err != nil {
		return nil, err
	}

	selectStmt := spanner.Statement{
		SQL: fmt.Sprintf("SELECT %s FROM %s", columnNames, repo.tableName),
	}

	iter := repo.spannerClient.Single().Query(c, selectStmt)
	defer iter.Stop()

	if err := iter.Do(func(row *spanner.Row) error {
		var masterDataVersion master.MasterDataVersion
		if err := row.ToStruct(&masterDataVersion); err != nil {
			return err // 行のデコードエラー
		}
		masterDataVersions = append(masterDataVersions, masterDataVersion)
		return nil
	}); err != nil {
		return nil, err
	}
	return masterDataVersions, nil
}
