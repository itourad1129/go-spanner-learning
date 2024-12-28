package controller

import (
	"cloud.google.com/go/spanner"
	"context"
	"github.com/gin-gonic/gin"
	"go-spanner-learning/domain"
	"net/http"
)

type UserRegisterController struct {
	SpannerClient       *spanner.Client
	UserRegisterUsecase domain.UserRegisterUsecase
}

func (urc *UserRegisterController) UserRegister(c *gin.Context) {

	var request domain.UserRegisterRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	_, err = urc.SpannerClient.ReadWriteTransaction(c, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		// ユーザー名の重複チェック
		checkUserInfo, err := urc.UserRegisterUsecase.GetUserByUserName(c, tx, request.Name)
		if err == nil {
			if checkUserInfo.Name != "" {
				// ユーザー名がすでに存在している場合、トランザクションをキャンセルし、レスポンスを返す
				return domain.ErrUserNameConflict
			}
		} else {
			return err
		}

		// ユーザー情報の作成
		userID, err := urc.UserRegisterUsecase.CreateUserInfo(c, tx, request.Name)
		if err != nil {
			return err
		}

		// ユーザー引き継ぎコードの作成
		transferCode, err := urc.UserRegisterUsecase.CreateUserTransfer(c, tx, userID)
		if err != nil {
			return err
		}

		userRegisterResponse := domain.UserRegisterResponse{
			UserID:       userID,
			TransferCode: transferCode,
		}

		c.Set("userRegisterResponse", userRegisterResponse)
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	if userRegisterResponse, exists := c.Get("userRegisterResponse"); exists {
		c.JSON(http.StatusOK, userRegisterResponse)
	}
}
