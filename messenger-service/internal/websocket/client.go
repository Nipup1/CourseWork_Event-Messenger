package webs

import (
	"context"
	"encoding/json"
	"go/course_work_6/internal/models"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

func (c *Client) ReadPump(h *Hub) {
    defer func() {
        h.unregister <- c
        c.Conn.Close()
        log.Info().Int("userID", c.UserID).Msg("Client disconnected")
    }()

    c.Conn.SetReadLimit(512)
    c.Conn.SetPongHandler(func(string) error {
        log.Debug().Int("userID", c.UserID).Msg("Received pong")
        return nil
    })

    for {
        _, message, err := c.Conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Error().Err(err).Int("userID", c.UserID).Msg("Unexpected close error")
            } else {
                log.Info().Err(err).Int("userID", c.UserID).Msg("Client disconnected or error reading message")
            }
            break
        }

        // Проверка на пустое сообщение
        if len(message) == 0 {
            log.Warn().Int("userID", c.UserID).Msg("Received empty message")
            continue
        }

        var msg models.Message
        if err := json.Unmarshal(message, &msg); err != nil {
            log.Error().Err(err).Int("userID", c.UserID).Msg("Error unmarshaling message")
            continue
        }

        if msg.ChatID == 0 || msg.Content == "" {
            log.Warn().Int("userID", c.UserID).Msg("Invalid message: ChatID or Content is empty")
            continue
        }
        msg.SenderID = c.UserID

        savedMsg, err := h.chatService.CreateMessage(context.Background(), msg.ChatID, msg.SenderID, msg.Content)
        if err != nil {
            log.Error().Err(err).Int("userID", c.UserID).Int("chatID", msg.ChatID).Msg("Error saving message")
            continue
        }

        msgBytes, err := json.Marshal(savedMsg)
        if err != nil {
            log.Error().Err(err).Int("userID", c.UserID).Msg("Error marshaling message for Redis")
            continue
        }
        if err := h.redis.Publish(context.Background(), "messages", msgBytes).Err(); err != nil {
            log.Error().Err(err).Int("userID", c.UserID).Msg("Error publishing to Redis")
        }
    }
}

func (c *Client) WritePump() {
    defer func() {
        c.Conn.Close()
        log.Info().Int("userID", c.UserID).Msg("WritePump stopped")
    }()

    ticker := time.NewTicker(30 * time.Second) // Отправляем ping каждые 30 секунд
    defer ticker.Stop()

    for {
        select {
        case message, ok := <-c.Send:
            if !ok {
                log.Warn().Int("userID", c.UserID).Msg("Send channel closed")
                return
            }
            msgBytes := messageBytes(message) // Исправлено: используем messageBytes
            if msgBytes == nil {
                log.Warn().Int("userID", c.UserID).Msg("Failed to serialize message, skipping")
                continue
            }
            if err := c.Conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
                log.Error().Err(err).Int("userID", c.UserID).Msg("Error writing message")
                return
            }
        case <-ticker.C:
            if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                log.Error().Err(err).Int("userID", c.UserID).Msg("Error sending ping")
                return
            }
            log.Debug().Int("userID", c.UserID).Msg("Sent ping")
        }
    }
}

// Вспомогательная функция для сериализации сообщения
func messageBytes(message models.Message) []byte {
    msgBytes, err := json.Marshal(message)
    if err != nil {
        log.Error().Err(err).Msg("Error marshaling message")
        return nil
    }
    return msgBytes
}