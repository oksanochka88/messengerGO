package models

import (
	"database/sql"
	"time"
)

// ChatParticipant represents a participant in a chat
type ChatParticipant struct {
	ChatID   int       `json:"chat_id"`
	UserID   int       `json:"user_id"`
	JoinedAt time.Time `json:"joined_at"`
}

func GetChatParticipants(db *sql.DB, chatID int) ([]ChatParticipant, error) {
	rows, err := db.Query("SELECT chat_id, user_id, joined_at FROM chatparticipants WHERE chat_id = $1", chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []ChatParticipant
	for rows.Next() {
		var participant ChatParticipant
		if err := rows.Scan(&participant.ChatID, &participant.UserID, &participant.JoinedAt); err != nil {
			return nil, err
		}
		participants = append(participants, participant)
	}
	return participants, nil
}
