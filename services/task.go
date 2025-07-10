package services

import (
	"taskive/models"
	"time"

	"gorm.io/gorm"
)

type TaskService struct {
	db *gorm.DB
}

func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{db: db}
}

type CreateTaskInput struct {
	Title       string            `json:"title" validate:"required"`
	Description string            `json:"description"`
	Status      models.TaskStatus `json:"status"`
	Priority    models.TaskPriority `json:"priority"`
	DueDate     time.Time         `json:"due_date"`
	AssigneeID  *uint            `json:"assignee_id"`
}

type UpdateTaskInput struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      models.TaskStatus `json:"status"`
	Priority    models.TaskPriority `json:"priority"`
	DueDate     time.Time         `json:"due_date"`
	AssigneeID  *uint            `json:"assignee_id"`
}

func (s *TaskService) Create(projectID uint, input CreateTaskInput) (*models.Task, error) {
	task := &models.Task{
		ProjectID:   projectID,
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Priority:    input.Priority,
		DueDate:     input.DueDate,
		AssigneeID:  input.AssigneeID,
	}

	if err := s.db.Create(task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) Update(taskID uint, input UpdateTaskInput) (*models.Task, error) {
	var task models.Task
	if err := s.db.First(&task, taskID).Error; err != nil {
		return nil, err
	}

	if input.Title != "" {
		task.Title = input.Title
	}
	task.Description = input.Description
	if input.Status != "" {
		task.Status = input.Status
	}
	if input.Priority != "" {
		task.Priority = input.Priority
	}
	if !input.DueDate.IsZero() {
		task.DueDate = input.DueDate
	}
	task.AssigneeID = input.AssigneeID

	if err := s.db.Save(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *TaskService) Delete(taskID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("task_id = ?", taskID).Delete(&models.Comment{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.Task{}, taskID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *TaskService) GetProjectTasks(projectID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := s.db.Where("project_id = ?", projectID).
		Preload("Assignee").
		Find(&tasks).Error
	return tasks, err
}

func (s *TaskService) GetByID(taskID uint) (*models.Task, error) {
	var task models.Task
	if err := s.db.Preload("Assignee").First(&task, taskID).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *TaskService) UpdateStatus(taskID uint, status models.TaskStatus) error {
	return s.db.Model(&models.Task{}).Where("id = ?", taskID).Update("status", status).Error
} 