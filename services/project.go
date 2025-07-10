package services

import (
	"taskive/models"
	"time"

	"gorm.io/gorm"
)

type ProjectService struct {
	db *gorm.DB
}

func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{db: db}
}

type CreateProjectInput struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type UpdateProjectInput struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

func (s *ProjectService) Create(userID uint, input CreateProjectInput) (*models.Project, error) {
	project := &models.Project{
		Name:        input.Name,
		Description: input.Description,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		OwnerID:     userID,
	}

	tx := s.db.Begin()
	if err := tx.Create(project).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	member := &models.Member{
		ProjectID: project.ID,
		UserID:    userID,
		Role:      models.MemberRoleOwner,
	}

	if err := tx.Create(member).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) Update(projectID uint, input UpdateProjectInput) (*models.Project, error) {
	var project models.Project
	if err := s.db.First(&project, projectID).Error; err != nil {
		return nil, err
	}

	if input.Name != "" {
		project.Name = input.Name
	}
	project.Description = input.Description
	if !input.StartDate.IsZero() {
		project.StartDate = input.StartDate
	}
	if !input.EndDate.IsZero() {
		project.EndDate = input.EndDate
	}

	if err := s.db.Save(&project).Error; err != nil {
		return nil, err
	}

	return &project, nil
}

func (s *ProjectService) Delete(projectID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("project_id = ?", projectID).Delete(&models.Member{}).Error; err != nil {
			return err
		}
		if err := tx.Where("project_id = ?", projectID).Delete(&models.Task{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.Project{}, projectID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *ProjectService) GetUserProjects(userID uint) ([]models.Project, error) {
	var projects []models.Project
	err := s.db.Joins("JOIN project_members ON projects.id = project_members.project_id").
		Where("project_members.user_id = ?", userID).
		Find(&projects).Error
	return projects, err
}

func (s *ProjectService) GetByID(projectID uint) (*models.Project, error) {
	var project models.Project
	if err := s.db.First(&project, projectID).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (s *ProjectService) AddMember(projectID, userID uint, role models.MemberRole) error {
	member := &models.Member{
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
	}
	return s.db.Create(member).Error
} 