package api

import (
	"encoding/json"
	"fmt"
	"go/course_work_6/internal/chat"
	"go/course_work_6/internal/models"
	webs "go/course_work_6/internal/websocket"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type MessengerHandler struct {
	chatService *chat.ChatService
	hub         *webs.Hub
	authService *chat.AuthService
}

func (h *MessengerHandler) CreateChat() any {
	panic("unimplemented")
}

type AuthService interface {
	ValidateToken(token string) (int, error)
}

func NewMessengerHandler(chatService *chat.ChatService, hub *webs.Hub, authService *chat.AuthService, router *http.ServeMux) {
	handler :=  &MessengerHandler{
		chatService: chatService,
		hub:         hub,
		authService: authService,
	}

	router.Handle("POST /api/chats", AuthMiddleware(http.HandlerFunc(handler.HandleChat), authService))
	router.Handle("GET /api/chats", AuthMiddleware(http.HandlerFunc(handler.HandleChatsList), authService))
    router.Handle("GET /api/chats/{id}/messages", AuthMiddleware(http.HandlerFunc(handler.HandleMessages), authService))
	router.HandleFunc("GET /api/chats/{id}/users", handler.HandleUserList)
    router.HandleFunc("GET /api/ws", handler.HandleWebSocket)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *MessengerHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	var req HandleChatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	chatID, err := h.chatService.CreateChat(req.Name, req.MemberIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := &HandleChatsResponse{
		ChatId: chatID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *MessengerHandler) HandleChatsList(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextEmailKey).(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	chats, err := h.chatService.GetUserChats(r.Context(), userID)
	if err != nil {
		log.Printf("Error fetching chats: %v", err)
		http.Error(w, "Failed to fetch chats", http.StatusInternalServerError)
		return
	}

	resp := &HandleChatsListResponse{
		Chats: chats,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *MessengerHandler) HandleMessages(w http.ResponseWriter, r *http.Request) {
	chatID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	messages, err := h.chatService.GetMessages(r.Context(), chatID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := &HandleMessagesResponse{
		Messages: messages,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *MessengerHandler) HandleUserList(w http.ResponseWriter, r *http.Request){
	chatID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	_, users, err := h.chatService.GetUsersInChat(r.Context(), chatID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := &HandleUserListResponse{
		UserIDs: users,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *MessengerHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("WebSocket connection attempt")
	token := r.URL.Query().Get("token")
	if token == "" {
		log.Error().Msg("Token required")
		http.Error(w, "Token required", http.StatusUnauthorized)
		return
	}

	userID, err := h.authService.ValidateToken(token)
	if err != nil {
		log.Error().Err(err).Msg("Invalid token")
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	log.Info().Int("userID", userID).Msg("Token validated")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error upgrading to WebSocket")
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}
	log.Info().Msg("WebSocket connection established")

	client := &webs.Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan models.Message),
	}

	h.hub.Register <- client
	log.Info().Int("userID", userID).Msg("Client registered")

	go client.WritePump()
	client.ReadPump(h.hub)
}
