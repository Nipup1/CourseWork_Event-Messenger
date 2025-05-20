package chat

import (
	"context"
    "fmt"
	"go/course_work_6/internal/models"

	"github.com/jackc/pgx/v5"
)

type ChatService struct {
    db *pgx.Conn
}

func NewChatService(db *pgx.Conn) *ChatService {
    return &ChatService{db: db}
}

func (s *ChatService) CreateChat(chatName string, memberIDs []int) (int, error) {
    ctx := context.Background()

    tx, err := s.db.Begin(ctx)
    if err != nil {
        return 0, err
    }
    defer tx.Rollback(ctx)

    var chatID int
    err = tx.QueryRow(ctx, `
        INSERT INTO chats (name, created_at)
        VALUES ($1, CURRENT_TIMESTAMP)
        RETURNING id
    `, chatName).Scan(&chatID)
    if err != nil {
        return 0, err
    }

    for _, memberID := range memberIDs {
        _, err = tx.Exec(ctx, `
            INSERT INTO chat_members (chat_id, user_id)
            VALUES ($1, $2)
        `, chatID, memberID)
        if err != nil {
            return 0, err
        }
    }

    if err := tx.Commit(ctx); err != nil {
        return 0, err
    }

    return chatID, nil
}

func (s *ChatService) GetUserChats(ctx context.Context, userID int) ([]models.Chat, error) {
    query := `
        SELECT c.id, c.name, c.created_at
        FROM chats c
        JOIN chat_members cm ON c.id = cm.chat_id
        WHERE cm.user_id = $1
        ORDER BY c.created_at DESC
    `
    rows, err := s.db.Query(ctx, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to query chats: %w", err)
    }
    defer rows.Close()

    var chats []models.Chat
    for rows.Next() {
        var chat models.Chat
        if err := rows.Scan(&chat.ID, &chat.Name, &chat.CreatedAt); err != nil {
            return nil, fmt.Errorf("failed to scan chat: %w", err)
        }
        chats = append(chats, chat)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating chats: %w", err)
    }

    return chats, nil
}

func (s *ChatService) GetUsersInChat (ctx context.Context, chatID int) (map[int]struct{}, []int, error){
    query := `
        SELECT user_id
        FROM chat_members
        WHERE chat_id = $1;
    `
    rows, err := s.db.Query(ctx, query, chatID)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to query chats: %w", err)
    }
    defer rows.Close()

    usersMap := make(map[int]struct{})
    usersSlice := make([]int, 0)
    for rows.Next() {
        var userID int
        if err := rows.Scan(&userID); err != nil {
            return nil, nil, fmt.Errorf("failed to scan chat: %w", err)
        }
        usersMap[userID] = struct{}{}
        usersSlice = append(usersSlice, userID)
    }

    if err := rows.Err(); err != nil {
        return nil, nil, fmt.Errorf("error iterating chats: %w", err)
    }

    return usersMap, usersSlice, nil
}