package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/vladopadikk/go-chat/internal/chat"
	"github.com/vladopadikk/go-chat/internal/messages"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Handler struct {
	hub            *Hub
	chatService    *chat.Service
	messageService *messages.Service
}

func NewHandler(hub *Hub, chatService *chat.Service, msgService *messages.Service) *Handler {
	return &Handler{hub, chatService, msgService}
}

func (h *Handler) ServeWS(ctx *gin.Context) {
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

	chatList, err := h.chatService.GetChatsList(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get chats"})
		return
	}

	var chatIDs []int64
	for _, ch := range chatList.Chats {
		chatIDs = append(chatIDs, ch.ID)
	}

	if len(chatIDs) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user has no chats"})
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := NewClient(h.hub, conn, userID, chatIDs, h.messageService)
	h.hub.register <- client

	go client.WritePump()
	go client.ReadPump()

	log.Printf("WebSocket connection established: userID=%d", userID)
}

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.GET("/ws", h.ServeWS)
}
