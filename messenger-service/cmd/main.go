package main

import (
	"context"
	"go/course_work_6/internal/api"
	"go/course_work_6/internal/chat"
	"go/course_work_6/internal/config"
	"go/course_work_6/internal/storage"
    "go/course_work_6/internal/websocket"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
    cfg := config.MustLoad()

    pgConn := storage.NewPostgresConn(context.Background(), cfg.SroragePath)
    redisClient := storage.NewRedisClient(cfg.RedisPath, "")
    defer pgConn.Close(context.Background())

    chatService := chat.NewChatService(pgConn)
    hub := webs.NewHub(redisClient.Client, chatService)
    authService := chat.NewAuthService(cfg.Secret, cfg.GRPCAddres)

    go hub.Run()

    router := http.NewServeMux()
    api.NewMessengerHandler(chatService, hub, authService, router)
    api.NewAuhtHandler(authService, router)

    log.Info().Msgf("Starting server on %s", cfg.Hostname)
    server := http.Server{
		Addr: cfg.Hostname,
		Handler: api.CORS(router),
	}

    stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

    go func (stop chan os.Signal)  {
        if err := server.ListenAndServe(); err != nil {
            log.Fatal().Msgf("Server failed: %v", err)
        }

        stop <- syscall.SIGINT
    }(stop)
	
    <-stop

    server.Shutdown(context.Background())
	log.Info().Msg("server stoped")
}