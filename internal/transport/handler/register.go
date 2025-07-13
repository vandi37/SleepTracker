package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vandi37/SleepTracker/models"
	"github.com/vandi37/SleepTracker/pkg/logger"
	"go.uber.org/zap"
)

func RegisterHandler(allowed []string, handler *Handler, l *zap.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.CustomRecovery(func(c *gin.Context, e any) {
		err := fmt.Errorf("%v", e)
		logger.Error(c, "got a panic in handler", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Internal(err))
	}))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowed,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(AddLogger(l))
	users := r.Group("/users")
	users.POST("/register", handler.Register)
	users.POST("/login", handler.Login)
	users.POST("/refresh", handler.Refresh)
	users.GET("/", handler.LoginMiddleware, handler.GetSelf)
	users.GET("/:id", handler.LoginMiddleware, handler.GetUser)
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
