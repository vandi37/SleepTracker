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
	users.PATCH("/password")
	users.DELETE("/")
	friends := r.Group("/friends")
	friends.POST("/request")
	friends.POST("/accept")
	friends.GET("/")
	friends.DELETE("/:id")

	friends.GET("/:id/history/:page")
	friends.GET("/:id/table/:page")

	sleep := r.Group("/sleep")
	sleep.POST("/sleep")
	sleep.PUT("/:id")
	sleep.DELETE("/:id")

	r.GET("/history/:page")
	r.GET("/table/:page")

	return r
}
