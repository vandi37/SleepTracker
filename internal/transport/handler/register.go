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
	users.PUT("/", handler.LoginMiddleware, handler.UpdateUser)
	users.PATCH("/password", handler.LoginMiddleware, handler.UpdatePassword)
	users.DELETE("/", handler.LoginMiddleware, handler.DeleteUser)
	friends := r.Group("/friends", handler.LoginMiddleware)
	friends.POST("/request", handler.Request)
	friends.POST("/accept", handler.Accept)
	friends.GET("/", handler.GetFriendships)
	friends.DELETE("/:id", handler.DeleteFriendship)

	friends.GET("/:id/history/:page", handler.GetFriendSleeps)
	friends.GET("/:id/table/:page", handler.GetFriendScores)

	sleep := r.Group("/sleep", handler.LoginMiddleware)
	sleep.POST("/", handler.EnterSleep)
	sleep.PUT("/:id", handler.UpdateSleep)
	sleep.DELETE("/:id", handler.DeleteSleep)

	r.GET("/history/:page", handler.GetSleeps)
	r.GET("/table/:page", handler.GetScores)

	return r
}
