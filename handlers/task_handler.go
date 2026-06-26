package handlers

import (
	"net/http"

	"task-api/database"
	"task-api/models"

	"github.com/gin-gonic/gin"
)

func isValidStatus(s models.Status) bool {
	return s == models.StatusTodo || s == models.StatusInProgress || s == models.StatusDone
}

// POST /tasks - ساخت تسک جدید
func CreateTask(c *gin.Context) {
	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if task.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	if task.Status == "" {
		task.Status = models.StatusTodo
	} else if !isValidStatus(task.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be 'todo', 'in progress', or 'done'"})
		return
	}

	if result := database.DB.Create(&task); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GET /tasks - گرفتن همه تسک‌ها
func GetAllTasks(c *gin.Context) {
	var tasks []models.Task

	if result := database.DB.Find(&tasks); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GET /tasks/:id - گرفتن یک تسک با ID
func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if result := database.DB.First(&task, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// PUT /tasks/:id - آپدیت تسک
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

	if input.Status != "" && !isValidStatus(input.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be 'todo', 'in progress', or 'done'"})
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

	database.DB.Save(&task)
	c.JSON(http.StatusOK, task)
}

// DELETE /tasks/:id - حذف تسک
func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if result := database.DB.First(&task, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	database.DB.Delete(&task)
	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}
