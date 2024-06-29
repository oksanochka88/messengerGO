package auth

import (
	"backMessage/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//func Login(c *gin.Context) {
//	var loginRequest LoginRequest
//	if err := c.ShouldBindJSON(&loginRequest); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	user, err := models.AuthenticateUser(loginRequest.Username, loginRequest.Password)
//	if err != nil {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
//		return
//	}
//
//	token, err := CreateJWT(user.ID)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"token": token})
//}

func Login(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.AuthenticateUser(loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	tokens, err := models.GetTokensByUserID(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tokens"})
		return
	}

	var validToken string
	for _, t := range tokens {
		valid, err := models.IsTokenValid(t.Token)
		if err == nil && valid {
			validToken = t.Token
			break
		}
	}

	if validToken == "" {
		token, err := CreateJWT(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		err = models.SaveToken(user.ID, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
			return
		}
		validToken = token
	}

	c.JSON(http.StatusOK, gin.H{"token": validToken})
}
