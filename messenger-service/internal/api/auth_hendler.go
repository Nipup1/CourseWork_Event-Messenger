package api

import (
	"context"
	"encoding/json"
	"go/course_work_6/internal/chat"
	"net/http"
	"strconv"
)

type AuthHandler struct {
	authService *chat.AuthService
}

func NewAuhtHandler(authService *chat.AuthService, router *http.ServeMux){
	handler := &AuthHandler{
		authService: authService,
	}

	router.HandleFunc("POST /api/register", handler.Register)
	router.HandleFunc("POST /api/login", handler.Login)
	router.HandleFunc("GET /api/users", handler.Users)
	router.HandleFunc("GET /api/user/{id}", handler.User)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request){
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := h.authService.Register(context.Background(), req.Email, req.Password, req.FullName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := &RegisterResponse{
		UserId: int(userID),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request){
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(context.Background(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := &LoginResponse{
		Token: token,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Users(w http.ResponseWriter, r *http.Request){
	users, err := h.authService.Users(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := &UsersResponse{
		Users: users,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) User(w http.ResponseWriter, r *http.Request){
	userId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.authService.User(context.Background(), int64(userId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := &UserResponse{
		UserID: user.GetUserId(),
		Email: user.GetEmail(),
		FullName: user.GetFullName(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}


