package models

import "time"

type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string    `json:"title" gorm:"unique;not null" validate:"required,min=4,max=100"`
	Description string    `json:"description" validate:"max=500"`
	Status      Status    `json:"status" gorm:"type:varchar(20);default:'todo'"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
