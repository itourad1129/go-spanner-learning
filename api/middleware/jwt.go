package middleware

import (
	"cloud.google.com/go/spanner"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go-spanner-learning/domain"
	"go-spanner-learning/domain/user"
	"go-spanner-learning/repository"
	"log"
	"time"
)

var IdentityKey = "userID"
var LoginUserID int64 = 0

func HandlerJwtMiddleWare(authMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return func(context *gin.Context) {
		err := authMiddleware.MiddlewareInit()
		if err != nil {
			log.Fatal("authMiddleware.MiddlewareInit() Error:" + err.Error())
		}
	}
}

func NewJwtMiddleware(spannerClient *spanner.Client) (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(InitParams(spannerClient))
	return authMiddleware, err
}

func InitParams(spannerClient *spanner.Client) *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: IdentityKey,
		PayloadFunc: PayloadFunc,

		IdentityHandler: IdentityHandler(),
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals domain.UserLoginRequest
			if err := c.ShouldBindJSON(&loginVals); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			userID := loginVals.UserID
			LoginUserID = userID
			transferCode := loginVals.TransferCode

			ctx := c.Request.Context()
			userTransferRepo := repository.NewUserTransferRepository(spannerClient, "t_user_transfer")

			user, err := userTransferRepo.Authenticate(ctx, userID, transferCode)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			if user.TransferCode == "" {
				return nil, jwt.ErrFailedAuthentication
			}

			return user, nil
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}
}

func PayloadFunc(data interface{}) jwt.MapClaims {
	if user, ok := data.(user.UserTransfer); ok {
		return jwt.MapClaims{
			"userID": user.UserID,
		}
	}
	return jwt.MapClaims{}
}

func IdentityHandler() func(c *gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		return claims["userID"]
	}
}
