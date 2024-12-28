package route

import (
	"cloud.google.com/go/spanner"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func Setup(gin *gin.Engine, handleJwtMiddleware *jwt.GinJWTMiddleware, timeout time.Duration, spannerClient *spanner.Client) {
	publicRouter := gin.Group("")
	ChunkRouter(timeout, publicRouter, spannerClient)
	MasterDataRouter(timeout, publicRouter, spannerClient)
	UserLoginRouter(timeout, handleJwtMiddleware, publicRouter, spannerClient)
	UserRegisterRouter(timeout, publicRouter, spannerClient)
	gin.NoRoute(handleJwtMiddleware.MiddlewareFunc(), handleNoRoute())

	authRouter := gin.Group("auth", handleJwtMiddleware.MiddlewareFunc())
	authRouter.GET("/hello", helloHandler)
}

func handleNoRoute() func(c *gin.Context) {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	}
}

func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID := claims["userID"]
	c.JSON(http.StatusOK, gin.H{
		"userID": userID,
	})
	log.Printf("Hello claims: %#v\n", claims["userID"])
}
