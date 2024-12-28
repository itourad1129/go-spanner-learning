package controller

import (
	"cloud.google.com/go/spanner"
	"context"
	"github.com/gin-gonic/gin"
	"go-spanner-learning/api/middleware"
	"go-spanner-learning/domain"
	"net/http"
)

type UserLoginController struct {
	SpannerClient    *spanner.Client
	UserLoginUsecase domain.UserLoginUsecase
}

func (ulc *UserLoginController) UserLogin(c *gin.Context) {

	var userID = middleware.LoginUserID

	_, err := ulc.SpannerClient.ReadWriteTransaction(c, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {

		userLogin, err := ulc.UserLoginUsecase.InsertOrUpdate(c, tx, userID)
		if err != nil {
			return err
		}

		userLoginResponse := domain.UserLoginResponse{
			UserID:         userID,
			TotalLoginDays: userLogin.TotalLoginDays,
		}

		c.Set("userLoginResponse", userLoginResponse)
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	if userLoginResponse, exists := c.Get("userLoginResponse"); exists {
		c.JSON(http.StatusOK, userLoginResponse)
	}
}
