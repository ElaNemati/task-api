package handlers

import (
	"log/slog"
	"net/http"
	"task-api/appError"
	"task-api/database"
	"task-api/dto"
	"task-api/models"
	"task-api/response"
	"task-api/validations"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	var req dto.CreateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("CreateTask: invalid request body", "error", err.Error())
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	fieldsErr, err := validations.ValidateCreateTask(&req)
	if err != nil {
		slog.Warn("CreateTask: validation failed", "error", err.Error())
		if fieldsErr != nil {
			response.ValidationError(c, http.StatusBadRequest, "validation failed", fieldsErr)
		} else {
			response.Error(c, http.StatusBadRequest, err.Error())
		}
		return
	}

	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}

	if task.Status == "" {
		task.Status = models.StatusTodo
	}

	if result := database.DB.Create(&task); result.Error != nil {
		appErr := appError.FromDBError(result.Error)
		slog.Error("CreateTask: database error", "error", result.Error.Error())
		response.FromAppError(c, appErr)
		return
	}

	slog.Info("CreateTask: task created successfully", "task_id", task.ID, "title", task.Title)
	response.Success(c, http.StatusCreated, "task created successfully", task)
}

func GetAllTasks(c *gin.Context) {
	var params models.TaskQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	} else if params.Limit > 100 {
		params.Limit = 100
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}

	if params.Status != "" && !isValidStatus(params.Status) {
		response.Error(c, http.StatusBadRequest, "status must be 'todo', 'in progress', or 'done'")
		return
	}

	var tasks []models.Task
	var totalCount int64

	db := database.DB.Model(&models.Task{})

	if params.Status != "" {
		db = db.Where("status = ?", params.Status)
	}

	db.Count(&totalCount)

	db = db.Order("created_at " + params.Sort)

	offset := (params.Page - 1) * params.Limit
	db = db.Limit(params.Limit).Offset(offset)

	if result := db.Find(&tasks); result.Error != nil {
		appErr := appError.FromDBError(result.Error)
		slog.Error("GetAllTasks: database error", "error", result.Error.Error())
		response.FromAppError(c, appErr)
		return
	}

	totalPages := int(totalCount) / params.Limit
	if int(totalCount)%params.Limit != 0 {
		totalPages++
	}

	result := models.PaginatedResponse{
		Tasks:      tasks,
		TotalCount: totalCount,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}

	slog.Info("GetAllTasks: tasks retrieved", "count", len(tasks), "total", totalCount)
	response.Success(c, http.StatusOK, "tasks retrieved successfully", result)
}

func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if result := database.DB.First(&task, id); result.Error != nil {
		appErr := appError.FromDBError(result.Error)
		slog.Warn("GetTaskByID: task not found", "id", id)
		response.FromAppError(c, appErr)
		return
	}

	slog.Info("GetTaskByID: task retrieved", "task_id", id)
	response.Success(c, http.StatusOK, "task retrieved successfully", task)
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if result := database.DB.First(&task, id); result.Error != nil {
		appErr := appError.FromDBError(result.Error)
		slog.Warn("UpdateTask: task not found", "id", id)
		response.FromAppError(c, appErr)
		return
	}

	var req dto.UpdateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("UpdateTask: invalid request body", "id", id, "error", err.Error())
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	fieldsErr, err := validations.ValidateUpdateTask(&req)
	if err != nil {
		slog.Warn("UpdateTask: validation failed", "id", id, "error", err.Error())
		if fieldsErr != nil {
			response.ValidationError(c, http.StatusBadRequest, "validation failed", fieldsErr)
		} else {
			response.Error(c, http.StatusBadRequest, err.Error())
		}
		return
	}

	if req.Title != "" {
		task.Title = req.Title
	}

	if req.Description != "" {
		task.Description = req.Description
	}

	if req.Status != "" {
		task.Status = req.Status
	}

	if result := database.DB.Save(&task); result.Error != nil {
		appErr := appError.FromDBError(result.Error)
		slog.Error("UpdateTask: database error", "id", id, "error", result.Error.Error())
		response.FromAppError(c, appErr)
		return
	}

	slog.Info("UpdateTask: task updated successfully", "task_id", id)
	response.Success(c, http.StatusOK, "task updated successfully", task)
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if result := database.DB.First(&task, id); result.Error != nil {
		appErr := appError.FromDBError(result.Error)
		slog.Warn("DeleteTask: task not found", "id", id)
		response.FromAppError(c, appErr)
		return
	}

	if result := database.DB.Delete(&task); result.Error != nil {
		appErr := appError.FromDBError(result.Error)
		slog.Error("DeleteTask: database error", "id", id, "error", result.Error.Error())
		response.FromAppError(c, appErr)
		return
	}

	slog.Info("DeleteTask: task deleted successfully", "task_id", id)
	response.Success(c, http.StatusOK, "task deleted successfully", nil)
}

func isValidStatus(status models.Status) bool {
	return status == models.StatusTodo ||
		status == models.StatusInProgress ||
		status == models.StatusDone
}
