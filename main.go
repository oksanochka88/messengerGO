package main

import (
	"backMessage/database"
	"backMessage/routes"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	database.InitDB(database.ConnStr)

	//Загрузка переменных окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	router := routes.SetupRouter()

	// Запускает сервер на порту 8080
	router.Run(":8080")
}

//package main
//
//import (
//	"backMessage/controllers"
//	"github.com/gin-gonic/gin"
//	"github.com/gorilla/websocket"
//	_ "github.com/lib/pq"
//	"log"
//	"net/http"
//)

//func main() {
//	// Строка подключения к вашей базе данных PostgreSQL
//	connStr := "postgres://postgres:roman@localhost/message?sslmode=disable"
//
//	// Открытие соединения с базой данных
//	db, err := sql.Open("postgres", connStr)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer db.Close()
//
//	// Проверка соединения с базой данных
//	err = db.Ping()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println("Successfully connected to PostgreSQL!")
//
//	//// Вызов функции для вставки нового пользователя
//	//newUser := User{
//	//	Username:         "johndoe",
//	//	Email:            "johndoe@example.com",
//	//	Password:         "securepassword",
//	//	Gender:           "male",
//	//	Photo:            []byte{}, // Пример пустого массива байт для фото
//	//	UniqueID:         "123e4567-e89b-12d3-a456-426614174000",
//	//	RegistrationDate: "2024-06-22 12:00:00", // Пример даты регистрации
//	//}
//	//
//	//if err := insertUser(db, newUser); err != nil {
//	//	log.Fatal(err)
//	//}
//	//
//	//fmt.Println("User inserted successfully!")
//}
//
//// Структура для представления данных пользователя
//type User struct {
//	Username         string
//	Email            string
//	Password         string
//	Gender           string
//	Photo            []byte
//	UniqueID         string
//	RegistrationDate string
//}
//
//// Функция для вставки пользователя в базу данных
//func insertUser(db *sql.DB, user User) error {
//	// SQL-запрос на вставку данных
//	query := `INSERT INTO users (username, email, password, gender, photo, unique_id, registration_date)
//              VALUES ($1, $2, $3, $4, $5, $6, $7)`
//
//	// Выполнение SQL-запроса с данными пользователя
//	_, err := db.Exec(query, user.Username, user.Email, user.Password, user.Gender, user.Photo, user.UniqueID, user.RegistrationDate)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

//package main
//
//import (
//	"github.com/gin-gonic/gin"
//	"github.com/golang-jwt/jwt/v4"
//	"net/http"
//	"time"
//)
//
//var jwtKey = []byte("your_secret_key")
//
//type Claims struct {
//	UserID int `json:"user_id"`
//	jwt.RegisteredClaims
//}
//
//func GenerateToken(userID int) (string, error) {
//	expirationTime := time.Now().Add(24 * time.Hour)
//
//	claims := &Claims{
//		UserID: userID,
//		RegisteredClaims: jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(expirationTime),
//		},
//	}
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	tokenString, err := token.SignedString(jwtKey)
//	if err != nil {
//		return "", err
//	}
//
//	return tokenString, nil
//}
//
//func main() {
//	r := gin.Default()
//	r.POST("/login", func(c *gin.Context) {
//		// Аутентификация пользователя
//
//		// После успешной аутентификации генерируем токен
//		token, err := GenerateToken(123) // Пример userID
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
//			return
//		}
//
//		c.JSON(http.StatusOK, gin.H{"token": token})
//	})
//	r.Run()
//}

//package main
//
//import (
//	"backMessage/controllers"
//	"github.com/gin-gonic/gin"
//	"github.com/gorilla/websocket"
//	"log"
//	"net/http"
//)
//
//var upgrader = websocket.Upgrader{
//	CheckOrigin: func(r *http.Request) bool {
//		return true
//	},
//}
//
//func handleConnections(c *gin.Context) {
//	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
//	if err != nil {
//		log.Fatalf("Failed to upgrade to WebSocket: %v", err)
//	}
//	defer ws.Close()
//
//	for {
//		// Read message from browser
//		_, msg, err := ws.ReadMessage()
//		if err != nil {
//			log.Printf("Error reading message: %v", err)
//			break
//		}
//		log.Printf("Received: %s", msg)
//
//		// Write message back to browser
//		err = ws.WriteMessage(websocket.TextMessage, msg)
//		if err != nil {
//			log.Printf("Error writing message: %v", err)
//			break
//		}
//	}
//}
//
//func SetupRouter() *gin.Engine {
//	r := gin.Default()
//
//	r.POST("/register", controllers.Register)
//
//	userRoutes := r.Group("/user")
//	{
//		userRoutes.GET("/:user_id", controllers.GetUserProfile)
//		userRoutes.PUT("/:user_id", controllers.UpdateUserProfile)
//	}
//
//	r.GET("/ws", handleConnections)
//
//	return r
//}
//
//func main() {
//	r := SetupRouter()
//	log.Println("Server started on :8080")
//	err := r.Run(":8080")
//	if err != nil {
//		log.Fatalf("Failed to start server: %v", err)
//	}
//}
