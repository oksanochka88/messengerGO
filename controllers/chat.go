package controllers

import (
	"backMessage/models"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// CreateChatHandler создает новый чат
func CreateChatHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request struct {
		Name         string `json:"name"`
		Participants []int  `json:"participants"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Добавляем создателя чата в список участников
	intUserID, err := strconv.Atoi(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error converting user ID"})
		return
	}
	request.Participants = append(request.Participants, intUserID)

	// Соединение с базой данных
	db, err := sql.Open("postgres", "user=youruser password=yourpassword dbname=yourdb sslmode=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	chat := models.Chat{Name: request.Name, CreatedAt: time.Now()}
	if err := models.CreateChat(db, &chat, request.Participants); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chat": chat})
}

// GetChatsHandler возвращает все чаты пользователя
func GetChatsHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Соединение с базой данных
	db, err := sql.Open("postgres", "user=youruser password=yourpassword dbname=yourdb sslmode=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	chats, err := models.GetChatsByUserID(db, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chats": chats})
}

// SendMessageHandler отправляет сообщение в чат
func SendMessageHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	chatID := c.Param("chat_id")

	var request struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Соединение с базой данных
	db, err := sql.Open("postgres", "user=youruser password=yourpassword dbname=yourdb sslmode=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	chatIDInt, err := strconv.Atoi(chatID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	message := models.Message{
		ChatID:    chatIDInt,
		UserID:    userID.(string),
		Content:   request.Content,
		CreatedAt: time.Now(),
	}

	if err := models.CreateMessage(db, &message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}

// GetMessagesHandler возвращает все сообщения чата
func GetMessagesHandler(c *gin.Context) {
	chatID := c.Param("chat_id")

	// Соединение с базой данных
	db, err := sql.Open("postgres", "user=youruser password=yourpassword dbname=yourdb sslmode=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	messages, err := models.GetMessagesByChatID(db, chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}
