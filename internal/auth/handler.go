package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service}
}

func (h *Handler) LoginHandler(ctx *gin.Context) {
	var loginIn LoginInput

	if err := ctx.ShouldBindJSON(&loginIn); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	tokens, err := h.service.Login(ctx.Request.Context(), loginIn)
	if err != nil {
		switch err {
		case ErrUserNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case ErrInvalidPassword:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/login", handler.LoginHandler)
}
