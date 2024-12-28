package repository

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"go-spanner-learning/database"
	"go-spanner-learning/domain/master"
)

type chunkVersionRepository struct {
	spannerClient *spanner.Client
	tableName     string
}

func NewChunkVersionRepository(spannerClient *spanner.Client, tableName string) master.ChunkVersionRepository {
	return chunkVersionRepository{
		spannerClient: spannerClient,
		tableName:     tableName,
	}
}

func (repo chunkVersionRepository) GetChunkVersion(c context.Context, platformType int64) (master.ChunkVersion, error) {
	var chunkVersion master.ChunkVersion
	columnNames, columns, err := database.GetSpannerColumns(master.ChunkVersion{})
	if err != nil {
		return chunkVersion, err
	}

	selectStmt := spanner.Statement{
		SQL:    fmt.Sprintf("SELECT %s FROM %s /*@ FORCE_INDEX = idx_version_id_platform_type_desc */ WHERE %s = $1 ORDER BY %s DESC LIMIT 1", columnNames, repo.tableName, columns[master.PlatformType], columns[master.VersionID]),
		Params: map[string]interface{}{"p1": platformType},
	}

	iter := repo.spannerClient.Single().Query(c, selectStmt)
	defer iter.Stop()

	if err := iter.Do(func(row *spanner.Row) error {
		return row.ToStruct(&chunkVersion)
	}); err != nil {
		return chunkVersion, err
	}
	return chunkVersion, nil
}
