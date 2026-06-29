package validations

import (
	"errors"
	"strings"

	"task-api/models"
)

func ValidateCreateTask(task *models.Task) error {

	if strings.TrimSpace(task.Title) == "" {
		return errors.New("title is required")
	}

	if task.Status != "" && !isValidStatus(task.Status) {
		return errors.New("status must be one of: todo, in progress, done")
	}

	return nil
}

func ValidateUpdateTask(task *models.Task) error {

	if task.Title != "" && strings.TrimSpace(task.Title) == "" {
		return errors.New("title cannot be empty")
	}

	if task.Status != "" && !isValidStatus(task.Status) {
		return errors.New("status must be one of: todo, in progress, done")
	}

	return nil
}

func isValidStatus(status models.Status) bool {

	switch status {

	case models.StatusTodo,
		models.StatusInProgress,
		models.StatusDone:
		return true

	default:
		return false
	}
}
