package master

import (
	"context"
)

const (
	VersionID    = "VersionID"
	PlatformType = "PlatformType"
)

type ChunkVersion struct {
	VersionID      int64  `spanner:"version_id"`
	PlatformType   int64  `spanner:"platform_type"`
	DeploymentName string `spanner:"deployment_name"`
	ContentBuildID string `spanner:"content_build_id"`
}

type ChunkVersionRepository interface {
	GetChunkVersion(c context.Context, platformType int64) (ChunkVersion, error)
}
