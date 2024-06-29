package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Photo    []byte `json:"photo"`
	UniqueId string `json:"unique_id"`
	About    string `json:"about"`
}

func CreateUser(username, email, password, about string, photo []byte) error {
	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}

	// Строка подключения
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	// Подключение к базе данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return err
	}
	defer db.Close()

	fmt.Println("Successfully connected to PostgreSQL!")

	// Генерация уникального ID
	uniqueID := generateUniqueID()

	// Выполнение запроса на вставку данных пользователя
	_, err = db.Exec("INSERT INTO users (username, email, password, photo, unique_id, about, registration_date) VALUES ($1, $2, $3, $4, $5, $6, NOW())",
		username, email, hashedPassword, photo, uniqueID, about)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		return err
	}

	log.Printf("User %s created successfully", username)
	return nil
}

func AuthenticateUser(username, password string) (*User, error) {
	var user User

	connStr := "postgres://postgres:roman@localhost/postgres?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, err
	}
	defer db.Close()

	err = db.QueryRow("SELECT id, username, email, password, photo, unique_id, about FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Photo, &user.UniqueId, &user.About)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user found with username: %s", username)
			return nil, errors.New("invalid credentials")
		}
		log.Printf("Error querying user: %v", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Password mismatch for user: %s", username)
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

func GetUserByID(userID string) (*User, error) {
	var user User

	db, err := sql.Open("postgres", "your_connection_string")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.QueryRow("SELECT id, username, email, photo, unique_id, about FROM users WHERE id = $1", userID).
		Scan(&user.ID, &user.Username, &user.Email, &user.Photo, &user.UniqueId, &user.About)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUser(userID string, userUpdate UserUpdate) error {
	db, err := sql.Open("postgres", "your_connection_string")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE users SET username = $1, email = $2, age = $3 WHERE id = $4",
		userUpdate.Username, userUpdate.Email, userUpdate.Age, userID)
	return err
}

type UserUpdate struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

func generateUniqueID() string {
	return uuid.New().String()
}