package handlers

import (
	"log/slog"
	"net/http"
	"task-api/database"
	"task-api/models"
	"task-api/validations"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		slog.Warn("CreateTask: invalid request body", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fieldsErr, err := validations.ValidateCreateTask(&task)
	if err != nil {
		slog.Warn("CreateTask: validation failed", "error", err.Error())
		if fieldsErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"validation_errors": fieldsErr})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	if task.Status == "" {
		task.Status = models.StatusTodo
	}

	if result := database.DB.Create(&task); result.Error != nil {
		slog.Error("CreateTask: database error", "error", result.Error.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	slog.Info("CreateTask: task created successfully", "task_id", task.ID, "title", task.Title)
	c.JSON(http.StatusCreated, task)
}

func GetAllTasks(c *gin.Context) {
	var tasks []models.Task

	if result := database.DB.Find(&tasks); result.Error != nil {
		slog.Error("GetAllTasks: database error", "error", result.Error.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	slog.Info("GetAllTasks: tasks retrieved", "count", len(tasks))
	c.JSON(http.StatusOK, tasks)
}

func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if result := database.DB.First(&task, id); result.Error != nil {
		slog.Warn("GetTaskByID: task not found", "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	slog.Info("GetTaskByID: task retrieved", "task_id", id)
	c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if result := database.DB.First(&task, id); result.Error != nil {
		slog.Warn("UpdateTask: task not found", "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		slog.Warn("UpdateTask: invalid request body", "id", id, "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fieldsErr, err := validations.ValidateUpdateTask(&input)
	if err != nil {
		slog.Warn("UpdateTask: validation failed", "id", id, "error", err.Error())
		if fieldsErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"validation_errors": fieldsErr})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	if input.Title != "" {
		task.Title = input.Title
	}
	if input.Description != "" {
		task.Description = input.Description
	}
	if input.Status != "" {
		task.Status = input.Status
	}

	if result := database.DB.Save(&task); result.Error != nil {
		slog.Error("UpdateTask: database error", "id", id, "error", result.Error.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	slog.Info("UpdateTask: task updated successfully", "task_id", id)
	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if result := database.DB.First(&task, id); result.Error != nil {
		slog.Warn("DeleteTask: task not found", "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if result := database.DB.Delete(&task); result.Error != nil {
		slog.Error("DeleteTask: database error", "id", id, "error", result.Error.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	slog.Info("DeleteTask: task deleted successfully", "task_id", id)
	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}
