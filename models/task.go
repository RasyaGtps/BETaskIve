package models

import (
	"time"

	"gorm.io/gorm"
)

type TaskStatus string
type TaskPriority string

const (
	TaskStatusTodo       TaskStatus = "TODO"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusDone       TaskStatus = "DONE"

	TaskPriorityLow    TaskPriority = "LOW"
	TaskPriorityMedium TaskPriority = "MEDIUM"
	TaskPriorityHigh   TaskPriority = "HIGH"
)

type Task struct {
	ID          uint         `gorm:"primarykey" json:"id"`
	ProjectID   uint         `json:"project_id"`
	Project     Project      `gorm:"foreignKey:ProjectID" json:"-"`
	Title       string       `gorm:"not null" json:"title" validate:"required"`
	Description string       `json:"description"`
	Status      TaskStatus   `gorm:"type:varchar(20);default:'TODO'" json:"status"`
	Priority    TaskPriority `gorm:"type:varchar(20);default:'MEDIUM'" json:"priority"`
	DueDate     time.Time    `json:"due_date"`
	AssigneeID  *uint        `json:"assignee_id"`
	Assignee    *User        `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`
	Comments    []Comment    `json:"comments,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	if t.Status == "" {
		t.Status = TaskStatusTodo
	}
	if t.Priority == "" {
		t.Priority = TaskPriorityMedium
	}
	return nil
} 