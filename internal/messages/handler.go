package messages

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultLimit  = 50
	maxLimit      = 100
	defaultOffset = 0
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service}
}

func (h *Handler) SendMessageHandler(ctx *gin.Context) {
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

	var input SendMessageInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	msg, err := h.service.SendMessage(ctx, userID, input)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "user is not a member of the chat"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, msg)
}

func (h *Handler) GetMessagesHandler(ctx *gin.Context) {
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

	chatIDParam := ctx.Query("chat_id")
	if chatIDParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "miss chat id param"})
		return
	}

	chatID, err := strconv.ParseInt(chatIDParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat id"})
		return
	}

	limit := defaultLimit
	offset := defaultOffset

	if limitParam := ctx.Query("limit"); limitParam != "" {
		l, err := strconv.Atoi(limitParam)
		if err != nil || l <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
			return
		}
		if l > maxLimit {
			l = maxLimit
		}
		limit = l
	}

	if offsetParam := ctx.Query("offset"); offsetParam != "" {
		o, err := strconv.Atoi(offsetParam)
		if err != nil || o < 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
			return
		}
		offset = o
	}

	msgs, err := h.service.GetMessages(ctx, chatID, userID, limit, offset)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "user is not a member of the chat"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, msgs)
}

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	chats := r.Group("/messages")
	{
		chats.POST("/send", h.SendMessageHandler)
		chats.GET("/get", h.GetMessagesHandler)
	}
}
