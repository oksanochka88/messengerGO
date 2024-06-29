package routes

import (
	"backMessage/controllers"
	"backMessage/controllers/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Маршруты для аутентификации и авторизации
	r.POST("/login", auth.Login)
	r.POST("/register", auth.Register)

	// Маршруты для чатов и сообщений
	//r.POST("/chats", controllers.CreateChat)
	//r.GET("/chats", controllers.GetChats)
	//r.POST("/chats/:chat_id/messages", controllers.SendMessage)
	//r.GET("/chats/:chat_id/messages", controllers.GetMessages)

	// Проверка доступности сервиса
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Маршруты для пользователя
	userRoutes := r.Group("/user")
	{
		userRoutes.GET("/:user_id", controllers.GetUserProfile)
		userRoutes.PUT("/:user_id", controllers.UpdateUserProfile)
	}

	return r
}
