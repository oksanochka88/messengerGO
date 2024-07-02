package routes

import (
	"backMessage/controllers"
	"backMessage/controllers/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Маршруты для аутентификации и авторизации
	r.POST("/login", auth.Login)
	r.POST("/register", auth.Register)

	// Проверка доступности сервиса
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	authUsers := r.Group("/")
	authUsers.Use(auth.AuthMiddleware(os.Getenv("JWT_SECRET")))

	authUsers.POST("/chats", controllers.CreateChatHandler)
	authUsers.GET("/chats", controllers.GetChatsHandler)
	authUsers.POST("/chats/:chat_id/messages", controllers.SendMessageHandler)
	authUsers.GET("/chats/:chat_id/messages", controllers.GetMessagesHandler)

	authUsers.GET("/profile", controllers.GetUserProfile)

	// Маршруты для пользователя
	userRoutes := r.Group("/user")
	{
		//userRoutes.GET("/profile", controllers.GetUserProfile)
		userRoutes.PUT("/:user_id", controllers.UpdateUserProfile)
	}

	return r
}
