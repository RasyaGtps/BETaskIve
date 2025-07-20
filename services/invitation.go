package services

import (
	"taskive/models"
	"time"

	"gorm.io/gorm"
)

type InvitationService struct {
	db *gorm.DB
}

func NewInvitationService(db *gorm.DB) *InvitationService {
	return &InvitationService{db: db}
}

type InvitationResponse struct {
	ID          uint            `json:"id"`
	ProjectID   uint            `json:"project_id"`
	ProjectName string          `json:"project_name"`
	InviterID   uint            `json:"inviter_id"`
	InviterName string          `json:"inviter_name"`
	Role        models.MemberRole `json:"role"`
	Status      string          `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
}

func (s *InvitationService) GetUserInvitations(userID uint) ([]InvitationResponse, error) {
	var invitations []struct {
		models.Member
		ProjectName string
		InviterName string
	}

	err := s.db.Table("project_members").
		Select("project_members.*, projects.name as project_name, users.name as inviter_name").
		Joins("JOIN projects ON project_members.project_id = projects.id").
		Joins("JOIN users ON projects.owner_id = users.id").
		Where("project_members.user_id = ? AND project_members.status = ?", userID, "PENDING").
		Find(&invitations).Error

	if err != nil {
		return nil, err
	}

	var response []InvitationResponse
	for _, inv := range invitations {
		response = append(response, InvitationResponse{
			ID:          inv.ProjectID, // Using composite key as ID
			ProjectID:   inv.ProjectID,
			ProjectName: inv.ProjectName,
			InviterID:   inv.UserID,
			InviterName: inv.InviterName,
			Role:        inv.Role,
			Status:      "PENDING",
			CreatedAt:   inv.CreatedAt,
		})
	}

	return response, nil
}

func (s *InvitationService) RespondToInvitation(userID, projectID uint, accept bool) error {
	tx := s.db.Begin()

	var member models.Member
	if err := tx.Where("project_id = ? AND user_id = ? AND status = ?", 
		projectID, userID, "PENDING").First(&member).Error; err != nil {
		tx.Rollback()
		return err
	}

	if accept {
		member.Status = "ACCEPTED"
		if err := tx.Save(&member).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if err := tx.Delete(&member).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
} 