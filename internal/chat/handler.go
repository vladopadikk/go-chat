package chat

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

func (h *Handler) CreatePrivateChatHandler(ctx *gin.Context) {
	userIDAny, exist := ctx.Get("userID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized"})
		return
	}
	userID, ok := userIDAny.(int64)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	var input CreatePrivateChatInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	chat, err := h.service.CreatePrivateChat(ctx.Request.Context(), userID, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, chat)
}

func (h *Handler) CreateGroupChatHandler(ctx *gin.Context) {
	userIDAny, exist := ctx.Get("userID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized"})
	}
	userID := userIDAny.(int64)

	var input CreateGroupChatInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	chat, err := h.service.CreateGroupChat(ctx.Request.Context(), userID, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, chat)
}

func (h *Handler) GetChatsHandler(ctx *gin.Context) {
	userIDAny, exist := ctx.Get("userID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized"})
	}
	userID := userIDAny.(int64)

	chatList, err := h.service.GetChatsList(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, chatList)
}

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	chats := r.Group("/chats")
	{
		chats.POST("/private", h.CreatePrivateChatHandler)
		chats.POST("/group", h.CreateGroupChatHandler)
		chats.GET("", h.GetChatsHandler)
	}
}
