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
