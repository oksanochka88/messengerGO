package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type Chat struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func GetUserIDsByNicknames(db *sql.DB, nicknames []string) ([]int, error) {
	if len(nicknames) == 0 {
		log.Println("No nicknames provided")
		return nil, errors.New("no nicknames provided")
	}

	// Построение запроса с использованием подстановок $1, $2 и т.д.
	queryPlaceholders := make([]string, len(nicknames))
	for i := range nicknames {
		queryPlaceholders[i] = fmt.Sprintf("$%d", i+1)
	}
	query := fmt.Sprintf("SELECT id FROM users WHERE username IN (%s)", strings.Join(queryPlaceholders, ","))
	log.Printf("Generated query: %s", query)

	args := make([]interface{}, len(nicknames))
	for i, nickname := range nicknames {
		args[i] = nickname
	}
	log.Printf("Arguments for query: %v", args)

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		userIDs = append(userIDs, id)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error with rows: %v", err)
		return nil, err
	}

	log.Printf("Retrieved user IDs: %v", userIDs)
	return userIDs, nil
}

// CreateChat создает новый чат и добавляет участников
func CreateChat(db *sql.DB, chat *Chat, participants []int) error {
	log.Println("Starting CreateChat function")

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return err
	}
	log.Println("Transaction started")

	var chatID int
	err = tx.QueryRow("INSERT INTO chats (name, created_at) VALUES ($1, $2) RETURNING id", chat.Name, chat.CreatedAt).Scan(&chatID)
	if err != nil {
		log.Printf("Error inserting into chats table: %v", err)
		tx.Rollback()
		return err
	}
	log.Printf("Chat ID: %d", chatID)

	stmt, err := tx.Prepare("INSERT INTO chatparticipants (chat_id, user_id, joined_at) VALUES ($1, $2, $3)")
	if err != nil {
		log.Printf("Error preparing statement for chatparticipants: %v", err)
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	log.Println("Statement prepared for chatparticipants")

	for _, userID := range participants {
		log.Printf("Adding participant userID: %d", userID)
		if _, err := stmt.Exec(chatID, userID, time.Now()); err != nil {
			log.Printf("Error inserting into chatparticipants: %v", err)
			tx.Rollback()
			return err
		}
	}
	log.Println("Participants added")

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}
	log.Println("Transaction committed successfully")

	return nil
}

// GetChatsByUserID получает все чаты для конкретного пользователя
func GetChatsByUserID(db *sql.DB, userID string) ([]Chat, error) {
	query := `
        SELECT c.id, c.name, c.created_at
        FROM chats c
        JOIN chatparticipants cp ON c.id = cp.chat_id
        WHERE cp.user_id = $1`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []Chat
	for rows.Next() {
		var chat Chat
		if err := rows.Scan(&chat.ID, &chat.Name, &chat.CreatedAt); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}
