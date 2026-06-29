package validations

import (
	"errors"
	"strings"
	"task-api/models"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func formatValidationError(err error) map[string]string {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return map[string]string{"error": err.Error()}
	}

	messages := make(map[string]string, len(validationErrors))
	for _, field := range validationErrors {
		switch field.Tag() {
		case "required":
			messages[field.Field()] = "is required"
		case "min":
			messages[field.Field()] = "must be at least " + field.Param() + " characters"
		case "max":
			messages[field.Field()] = "must be at most " + field.Param() + " characters"
		default:
			messages[field.Field()] = "is invalid"
		}
	}
	return messages
}

func ValidateCreateTask(task *models.Task) (map[string]string, error) {
	if err := validate.Struct(task); err != nil {
		return formatValidationError(err), err
	}

	if strings.TrimSpace(task.Title) == "" {
		return nil, errors.New("title is required")
	}

	if task.Status != "" && !isValidStatus(task.Status) {
		return nil, errors.New("status must be one of: todo, in progress, done")
	}

	return nil, nil
}

func ValidateUpdateTask(task *models.Task) (map[string]string, error) {
	if err := validate.StructPartial(task, "Title", "Description"); err != nil {
		return formatValidationError(err), err
	}

	if task.Title != "" && strings.TrimSpace(task.Title) == "" {
		return nil, errors.New("title cannot be empty")
	}

	if task.Status != "" && !isValidStatus(task.Status) {
		return nil, errors.New("status must be one of: todo, in progress, done")
	}

	return nil, nil
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
