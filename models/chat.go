package models

import (
	"backMessage/database"
	"log"
	"time"
)

type Chat struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateChat(name string) error {
	query := `INSERT INTO chats (name, created_at) VALUES ($1, $2)`
	_, err := database.DB.Exec(query, name, time.Now())
	if err != nil {
		log.Printf("Error inserting chat into database: %v", err)
		return err
	}

	log.Printf("Chat %s created successfully", name)
	return nil
}
