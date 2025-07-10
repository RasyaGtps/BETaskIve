package services

import (
	"taskive/models"

	"gorm.io/gorm"
)

type CommentService struct {
	db *gorm.DB
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{db: db}
}

type CreateCommentInput struct {
	Text string `json:"text" validate:"required"`
}

func (s *CommentService) Create(taskID, userID uint, input CreateCommentInput) (*models.Comment, error) {
	comment := &models.Comment{
		TaskID: taskID,
		UserID: userID,
		Text:   input.Text,
	}

	if err := s.db.Create(comment).Error; err != nil {
		return nil, err
	}

	if err := s.db.Preload("User").First(comment, comment.ID).Error; err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) GetTaskComments(taskID uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := s.db.Where("task_id = ?", taskID).
		Preload("User").
		Order("created_at DESC").
		Find(&comments).Error
	return comments, err
}

func (s *CommentService) Delete(commentID uint) error {
	return s.db.Delete(&models.Comment{}, commentID).Error
} 