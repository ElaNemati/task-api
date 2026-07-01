package models

import (
	"gorm.io/gorm"

	"time"

)

type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string         `json:"title" gorm:"unique;not null"`
	Description string         `json:"description"`
	Status      Status         `json:"status" gorm:"type:varchar(20);default:'todo'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}