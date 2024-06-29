package auth

import (
	"backMessage/models"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	About    string `json:"about"`
}

func Register(c *gin.Context) {
	var registerRequest RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userExists(registerRequest.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	err := models.CreateUser(registerRequest.Username, registerRequest.Email, registerRequest.Password, registerRequest.About, nil) // nil для photo в данном примере
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func userExists(username string) bool {
	connStr := "postgres://postgres:roman@localhost/postgres?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return false
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if user exists: %v", err)
		return false
	}

	return exists
}
