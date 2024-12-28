package route

import (
	"cloud.google.com/go/spanner"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go-spanner-learning/api/controller"
	"go-spanner-learning/repository"
	"go-spanner-learning/usecase"
	"net/http"
	"time"
)

func UserLoginRouter(timeout time.Duration, handleJwtMiddleware *jwt.GinJWTMiddleware, group *gin.RouterGroup, spannerClient *spanner.Client) {
	group.POST("/login", func(c *gin.Context) {
		handleJwtMiddleware.LoginHandler(c)
		if c.Writer.Status() == http.StatusOK {
			ulr := repository.NewUserLoginRepository(spannerClient, "t_user_login")
			utr := repository.NewUserTransferRepository(spannerClient, "t_user_transfer")
			ulc := controller.UserLoginController{
				SpannerClient:    spannerClient,
				UserLoginUsecase: usecase.NewUserLoginUsecase(ulr, utr, timeout),
			}
			ulc.UserLogin(c)
		}
	})
}
