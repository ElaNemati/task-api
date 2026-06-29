package handlers

import (
	"net/http"
	"task-api/validations"

	"task-api/database"
	"task-api/models"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validations.ValidateCreateTask(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if task.Status == "" {
		task.Status = models.StatusTodo
	}

	if result := database.DB.Create(&task); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func GetAllTasks(c *gin.Context) {
	var tasks []models.Task

	if result := database.DB.Find(&tasks); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if result := database.DB.First(&task, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if result := database.DB.First(&task, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validations.ValidateUpdateTask(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if result := database.DB.First(&task, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if result := database.DB.Delete(&task); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}
