package chat

import (
	"context"
	"fmt"
	"go/course_work_6/internal/models"
)

func (s *ChatService) CreateMessage(ctx context.Context, chatID, senderID int, content string) (models.Message, error) {
	var msg models.Message
	err := s.db.QueryRow(ctx, `
		INSERT INTO messages (chat_id, sender_id, content, created_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		RETURNING id, chat_id, sender_id, content, created_at
	`, chatID, senderID, content).Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Content, &msg.CreatedAt)
	if err != nil {
		return models.Message{}, fmt.Errorf("failed to create message: %w", err)
	}

	return msg, nil
}

//TODO: Доделать нормальную пагинацию
func (s *ChatService) GetMessages(ctx context.Context, chatID int) ([]models.Message, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, chat_id, sender_id, content, created_at
		FROM messages
		WHERE chat_id = $1
		ORDER BY created_at ASC
		LIMIT 50
	`, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// func (s *ChatService) CheckChatMembership(ctx context.Context, chatID, userID int) error {
// 	var exists bool
// 	err := s.db.QueryRow(ctx, `
// 		SELECT EXISTS (
// 			SELECT 1 FROM chat_members
// 			WHERE chat_id = $1 AND user_id = $2
// 		)
// 	`, chatID, userID).Scan(&exists)
// 	if err != nil {
// 		return fmt.Errorf("failed to check membership: %w", err)
// 	}
// 	if !exists {
// 		return fmt.Errorf("user %d is not a member of chat %d", userID, chatID)
// 	}
// 	return nil
// }