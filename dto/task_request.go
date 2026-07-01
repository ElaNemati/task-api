package dto

import (
	"task-api/models"
)

type CreateTaskRequest struct {
	Title       string        `json:"title" validate:"required,min=4,max=100"`
	Description string        `json:"description" validate:"max=500"`
	Status      models.Status `json:"status"`
}

type UpdateTaskRequest struct {
	Title       string        `json:"title" validate:"omitempty,min=4,max=100"`
	Description string        `json:"description" validate:"omitempty,max=500"`
	Status      models.Status `json:"status"`
}
