package api

import (
	"context"
	"net/http"
	"strings"
)

type key string
const(
	ContextEmailKey key = "EmailKey"
)

func AuthMiddleware(next http.Handler, authService AuthService) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
            return
        }
		
        userID, err := authService.ValidateToken(parts[1])
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), ContextEmailKey, userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func CORS(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
        // Разрешаем методы
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        // Разрешаем заголовки
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Обрабатываем preflight-запросы (OPTIONS)
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        // Передаем запрос следующему обработчику
        next.ServeHTTP(w, r)
	})
}