package domain

import (
	"context"
	"go-spanner-learning/domain/master"
)

type ChunkUsecase interface {
	GetChunkVersion(c context.Context, platformType int64) (master.ChunkVersion, error)
}

type GetChunkVersionRequest struct {
	PlatformType int64 `form:"platformType" binding:"required"`
}

type GetChunkVersionResponse struct {
	VersionID      int64  `json:"versionID"`
	PlatformType   int64  `json:"platformType"`
	DeploymentName string `json:"deploymentName"`
	ContentBuildID string `json:"contentBuildID"`
}
