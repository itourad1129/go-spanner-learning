package controller

import (
	"cloud.google.com/go/spanner"
	"github.com/gin-gonic/gin"
	"go-spanner-learning/domain"
	"net/http"
)

type ChunkController struct {
	SpannerClient *spanner.Client
	ChunkUsecase  domain.ChunkUsecase
}

func (cc *ChunkController) GetChunkVersion(c *gin.Context) {

	var request domain.GetChunkVersionRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	chunkVersion, err := cc.ChunkUsecase.GetChunkVersion(c, request.PlatformType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	chunkVersionResponse := domain.GetChunkVersionResponse{
		VersionID:      chunkVersion.VersionID,
		PlatformType:   chunkVersion.PlatformType,
		DeploymentName: chunkVersion.DeploymentName,
		ContentBuildID: chunkVersion.ContentBuildID,
	}

	c.JSON(http.StatusOK, chunkVersionResponse)
}
