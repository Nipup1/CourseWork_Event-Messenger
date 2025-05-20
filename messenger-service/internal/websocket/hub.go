package webs

import (
	"context"
	"encoding/json"
	"go/course_work_6/internal/chat"
	"go/course_work_6/internal/models"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Hub struct {
    clients    map[*Client]bool
    Broadcast  chan models.Message
    Register   chan *Client
    unregister chan *Client
    redis      *redis.Client
    chatService *chat.ChatService
    mu         sync.Mutex
}

type Client struct {
    UserID int
    Conn   *websocket.Conn
    Send   chan models.Message
}

func NewHub(redis *redis.Client, chatService *chat.ChatService) *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        Broadcast:  make(chan models.Message),
        Register:   make(chan *Client),
        unregister: make(chan *Client),
        redis:      redis,
        chatService: chatService,
    }
}

func (h *Hub) Run() {
    pubsub := h.redis.Subscribe(context.Background(), "messages")
    defer pubsub.Close()

    go func() {
        for msg := range pubsub.Channel() {
            var message models.Message
            json.Unmarshal([]byte(msg.Payload), &message)
            h.Broadcast <- message
        }
    }()

    for {
        select {
        case client := <-h.Register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()
        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                close(client.Send)
                delete(h.clients, client)
            }
            h.mu.Unlock()
        case message := <-h.Broadcast:
            h.mu.Lock()
            users := h.GetChatMembers(message.ChatID)
            for client := range h.clients {
                if _, in := users[client.UserID]; in{
                    select {
                    case client.Send <- message:
                    default:
                        close(client.Send)
                        delete(h.clients, client)
                    }
                }
            }
            h.mu.Unlock()
        }
    }
}

func (h *Hub) GetChatMembers(chatID int) map[int]struct{} {
    users, _, _:= h.chatService.GetUsersInChat(context.Background(), chatID)
    return users
}