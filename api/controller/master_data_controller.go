package controller

import (
	"cloud.google.com/go/spanner"
	"github.com/gin-gonic/gin"
	"go-spanner-learning/domain"
	"net/http"
)

type MasterDataVersionController struct {
	SpannerClient            *spanner.Client
	MasterDataVersionUsecase domain.MasterDataVersionUsecase
}

func (mdc *MasterDataVersionController) GetMasterDataVersion(c *gin.Context) {

	masterDataVersions, err := mdc.MasterDataVersionUsecase.GetMasterDataVersion(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	var masterDataVersionResponse []domain.GetMasterDataVersionResponse
	for _, masterDataVersion := range masterDataVersions {
		response := domain.GetMasterDataVersionResponse{
			MasterDataID: masterDataVersion.MasterDataID,
			Version:      masterDataVersion.Version,
			ChunkID:      masterDataVersion.ChunkID,
		}
		masterDataVersionResponse = append(masterDataVersionResponse, response)
	}

	c.JSON(http.StatusOK, masterDataVersionResponse)
}
