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
		return
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

func (h *Handler) UpdateUser(ctx *gin.Context) {
	id := ctx.GetInt64(ID_KEY)
	if id <= 0 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Internal(errors.New("got no id")))
		return
	}
	var req models.User
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.JsonError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}
	req.Id = id
	if err := h.service.UpdateUser(ctx, req); err != nil {
		ctx.AbortWithStatusJSON(err.Code(), err.JsonError())
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *Handler) UpdatePassword(ctx *gin.Context) {
	id := ctx.GetInt64(ID_KEY)
	if id <= 0 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Internal(errors.New("got no id")))
		return
	}
	var req models.Password
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.JsonError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}
	if err := h.service.UpdatePassword(ctx, id, req.Password); err != nil {
		ctx.AbortWithStatusJSON(err.Code(), err.JsonError())
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *Handler) DeleteUser(ctx *gin.Context) {
	id := ctx.GetInt64(ID_KEY)
	if id <= 0 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Internal(errors.New("got no id")))
		return
	}
	if err := h.service.DeleteUser(ctx, id); err != nil {
		ctx.AbortWithStatusJSON(err.Code(), err.JsonError())
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *Handler) Request(ctx *gin.Context) {
	id := ctx.GetInt64(ID_KEY)
	if id <= 0 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Internal(errors.New("got no id")))
		return
	}
	var req models.Id
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.JsonError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}
	fid, err := h.service.Request(ctx, id, req.Id)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code(), err.JsonError())
		return
	}
	ctx.JSON(http.StatusOK, models.Id{Id: fid})
}

func (h *Handler) Accept(ctx *gin.Context) {
	id := ctx.GetInt64(ID_KEY)
	if id <= 0 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Internal(errors.New("got no id")))
		return
	}
	var req models.Id
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.JsonError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}
	if err := h.service.Accept(ctx, req.Id, id); err != nil {
		ctx.AbortWithStatusJSON(err.Code(), err.JsonError())
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *Handler) GetFriendships(ctx *gin.Context) {
	id := ctx.GetInt64(ID_KEY)
	if id <= 0 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Internal(errors.New("got no id")))
		return
	}
	var req models.GetFriendships
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.JsonError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}
	req.UserId = id
	friendships, err := h.service.GetFriendships(ctx, req)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code(), err.JsonError())
		return
	}
	ctx.JSON(http.StatusOK, models.Friendships{Friendships: friendships})
}

func (h *Handler) DeleteFriendship(ctx *gin.Context) {
	id := ctx.GetInt64(ID_KEY)
	if id <= 0 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Internal(errors.New("got no id")))
		return
	}
	sid := ctx.GetString("id")
	friendshipId, parseErr := strconv.ParseInt(sid, 10, 64)
	if parseErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Invalid("id", sid))
		return
	}
	if err := h.service.DeleteFriendship(ctx, friendshipId, id); err != nil {
		ctx.AbortWithStatusJSON(err.Code(), err.JsonError())
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *Handler) EnterSleep(ctx *gin.Context) {
	id := ctx.GetInt64(ID_KEY)
	if id <= 0 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Internal(errors.New("got no id")))
		return
	}
	var req models.Sleep
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.JsonError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}
	req.UserId = id
	sid, err := h.service.EnterSleep(ctx, req)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code(), err.JsonError())
		return
	}
	ctx.JSON(http.StatusOK, models.Id{Id: sid})
}

func (h *Handler) UpdateSleep(ctx *gin.Context) {
	id := ctx.GetInt64(ID_KEY)
	if id <= 0 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Internal(errors.New("got no id")))
		return
	}
	var req models.Sleep
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.JsonError{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}
	req.UserId = id
	if err := h.service.UpdateSleep(ctx, req); err != nil {
		ctx.AbortWithStatusJSON(err.Code(), err.JsonError())
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *Handler) DeleteSleep(ctx *gin.Context) {
	id := ctx.GetInt64(ID_KEY)
	if id <= 0 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Internal(errors.New("got no id")))
		return
	}
	sid := ctx.GetString("id")
	sleepId, parseErr := strconv.ParseInt(sid, 10, 64)
	if parseErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Invalid("id", sid))
		return
	}
	if err := h.service.DeleteSleep(ctx, sleepId, id); err != nil {
		ctx.AbortWithStatusJSON(err.Code(), err.JsonError())
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *Handler) GetSleeps(ctx *gin.Context) {
	id := ctx.GetInt64(ID_KEY)
	if id <= 0 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Internal(errors.New("got no id")))
		return
	}
	sPage := ctx.GetString("page")
	page, parseErr := strconv.Atoi(sPage)
	if parseErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Invalid("page", page))
		return
	}
	sleeps, err := h.service.GetSleeps(ctx, id, page)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code(), err.JsonError())
		return
	}
	ctx.JSON(http.StatusOK, models.Sleeps{Sleeps: sleeps})
}

func (h *Handler) GetScores(ctx *gin.Context) {
	id := ctx.GetInt64(ID_KEY)
	if id <= 0 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Internal(errors.New("got no id")))
		return
	}
	sPage := ctx.GetString("page")
	page, parseErr := strconv.Atoi(sPage)
	if parseErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Invalid("page", page))
		return
	}
	scores, err := h.service.GetScores(ctx, id, page)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code(), err.JsonError())
		return
	}
	ctx.JSON(http.StatusOK, models.Scores{Scores: scores})
}

func (h *Handler) GetFriendSleeps(ctx *gin.Context) {
	id := ctx.GetInt64(ID_KEY)
	if id <= 0 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Internal(errors.New("got no id")))
		return
	}
	sid := ctx.GetString("id")
	friendshipId, parseErr := strconv.ParseInt(sid, 10, 64)
	if parseErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Invalid("id", sid))
		return
	}
	sPage := ctx.GetString("page")
	page, parseErr := strconv.Atoi(sPage)
	if parseErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Invalid("page", page))
		return
	}
	second, err := h.service.GetSecond(ctx, friendshipId, id)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code(), err.JsonError())
		return
	}
	sleeps, err := h.service.GetSleeps(ctx, second, page)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code(), err.JsonError())
		return
	}
	ctx.JSON(http.StatusOK, models.Sleeps{Sleeps: sleeps})
}
func (h *Handler) GetFriendScores(ctx *gin.Context) {
	id := ctx.GetInt64(ID_KEY)
	if id <= 0 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Internal(errors.New("got no id")))
		return
	}
	sid := ctx.GetString("id")
	friendshipId, parseErr := strconv.ParseInt(sid, 10, 64)
	if parseErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Invalid("id", sid))
		return
	}
	sPage := ctx.GetString("page")
	page, parseErr := strconv.Atoi(sPage)
	if parseErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Invalid("page", page))
		return
	}
	second, err := h.service.GetSecond(ctx, friendshipId, id)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code(), err.JsonError())
		return
	}
	scores, err := h.service.GetScores(ctx, second, page)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code(), err.JsonError())
		return
	}
	ctx.JSON(http.StatusOK, models.Scores{Scores: scores})
}
