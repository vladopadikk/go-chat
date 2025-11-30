package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service}
}

func (h *Handler) RegisterHandler(ctx *gin.Context) {
	var userIn UserInput

	if err := ctx.ShouldBindJSON(&userIn); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.service.Register(ctx.Request.Context(), userIn)
	if err != nil {
		if errors.Is(err, ErrEmailExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, UserResponse{
		ID:       id,
		Username: userIn.Username,
		Email:    userIn.Email,
	})
}

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/register", handler.RegisterHandler)
}
