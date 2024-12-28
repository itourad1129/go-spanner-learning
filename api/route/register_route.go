package route

import (
	"cloud.google.com/go/spanner"
	"github.com/gin-gonic/gin"
	"go-spanner-learning/api/controller"
	"go-spanner-learning/repository"
	"go-spanner-learning/usecase"
	"time"
)

func UserRegisterRouter(timeout time.Duration, group *gin.RouterGroup, spannerClient *spanner.Client) {
	uir := repository.NewUserInfoRepository(spannerClient, "t_user_info")
	utr := repository.NewUserTransferRepository(spannerClient, "t_user_transfer")
	urc := controller.UserRegisterController{
		SpannerClient:       spannerClient,
		UserRegisterUsecase: usecase.NewUserRegisterUsecase(uir, utr, timeout),
	}
	group.POST("/userRegister", urc.UserRegister)
}
