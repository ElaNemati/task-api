package main

import (
	"log"
	"log/slog"
	"task-api/database"
	"task-api/handlers"
	"task-api/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	closeLogs, err := logger.Init()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer closeLogs()

	database.Connect()
	slog.Info("Server starting", "port", "8080")

	r := gin.Default()

	r.POST("/tasks", handlers.CreateTask)
	r.GET("/tasks", handlers.GetAllTasks)
	r.GET("/tasks/:id", handlers.GetTaskByID)
	r.PUT("/tasks/:id", handlers.UpdateTask)
	r.DELETE("/tasks/:id", handlers.DeleteTask)

	r.Run(":8080")
}
