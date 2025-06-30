package handler

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(allowed []string) *gin.Engine {
	r := gin.New()
	r.Use(gin.CustomRecovery(func(c *gin.Context, err any) {

	}))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowed,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	users := r.Group("/users")
	users.POST("/register")
	users.POST("/login")
	users.POST("/refresh")
	users.GET("/")
	users.GET("/:id")
	users.PUT("/")
	users.DELETE("/")
	friends := r.Group("/friends")
	friends.POST("/request")
	friends.POST("/accept")
	friends.GET("/")
	friends.DELETE("/:id")

	friends.GET("/:id/history/:page")
	sleep := r.Group("/sleep")
	sleep.POST("/")
	sleep.GET("/history/:page")

	return r
}
