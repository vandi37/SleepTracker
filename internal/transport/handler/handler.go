package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vandi37/SleepTracker/internal/service"
	"github.com/vandi37/SleepTracker/models"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(ctx *gin.Context) {
	var req models.UserReq
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.JsonError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}
	res, err := h.service.Register(ctx.Request.Context(), req.Username, req.Nickname, req.Password, req.Birth)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code(), err.JsonError())
		return
	}
	ctx.JSON(http.StatusOK, res)
}
func (h *Handler) Login(ctx *gin.Context) {
	var req models.UserReq
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.JsonError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}
	res, err := h.service.Login(ctx, req.Username, req.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code(), err.JsonError())
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *Handler) Refresh(ctx *gin.Context) {
	var s string
	if err := json.NewDecoder(ctx.Request.Body).Decode(&s); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.JsonError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}
	res, err := h.service.Refresh(ctx, s)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code(), err.JsonError())
		return
	}
	ctx.JSON(http.StatusOK, res)
}
func (h *Handler) GetSelf(ctx *gin.Context) {
	id := ctx.GetInt64(ID_KEY)
	if id <= 0 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Internal(errors.New("got no id")))
	}
	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		ctx.AbortWithError(err.Code(), err.JsonError())
		return
	}
	ctx.JSON(http.StatusOK, user)
}
func (h *Handler) GetUser(ctx *gin.Context) {
	sid := ctx.GetString("id")
	id, parseErr := strconv.ParseInt(sid, 10, 64)
	if parseErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Invalid("id", sid))
		return
	}
	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		ctx.AbortWithError(err.Code(), err.JsonError())
		return
	}
	ctx.JSON(http.StatusOK, user)
}
