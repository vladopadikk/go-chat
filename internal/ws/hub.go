package ws

type Broadcast struct {
	ChatID int64
	Data   []byte
}

type Hub struct {
	clients    map[int64]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Broadcast
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int64]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Broadcast),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			for _, chatID := range client.chats {
				if h.clients[chatID] == nil {
					h.clients[chatID] = make(map[*Client]bool)
				}
				h.clients[chatID][client] = true
			}

		case client := <-h.unregister:
			for _, chatID := range client.chats {
				if clients, ok := h.clients[chatID]; ok {
					delete(clients, client)
					if len(clients) == 0 {
						delete(h.clients, chatID)
					}
				}
			}
			close(client.send)

		case msg := <-h.broadcast:
			if clients, ok := h.clients[msg.ChatID]; ok {
				for c := range clients {
					select {
					case c.send <- msg.Data:
					default:
						close(c.send)
						delete(clients, c)
					}
				}
			}
		}
	}
}
