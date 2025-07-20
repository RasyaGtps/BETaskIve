package services

import (
	"fmt"
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
		Status:    models.MemberStatusAccepted, // Owner langsung accepted
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
		Where("project_members.user_id = ? AND project_members.status = ?", userID, models.MemberStatusAccepted).
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
		Status:    models.MemberStatusPending, // Member baru status pending
	}
	return s.db.Create(member).Error
}

type ProjectInvitation struct {
	ID          uint      `json:"id"`
	ProjectID   uint      `json:"project_id"`
	ProjectName string    `json:"project_name"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}

func (s *ProjectService) GetUserInvitations(userID uint) ([]ProjectInvitation, error) {
	var invitations []ProjectInvitation

	// Debug: Print the SQL query
	query := s.db.Table("project_members").
		Select(`
			project_members.project_id as id,
			project_members.project_id as project_id,
			projects.name as project_name,
			project_members.role as role,
			project_members.created_at as created_at
		`).
		Joins("LEFT JOIN projects ON projects.id = project_members.project_id").
		Where("project_members.user_id = ? AND project_members.status = ?", userID, models.MemberStatusPending)

	fmt.Printf("SQL Query: %v\n", query.Statement.SQL.String())
	fmt.Printf("Query Values: userID=%d, status=%s\n", userID, models.MemberStatusPending)

	err := query.Find(&invitations).Error
	if err != nil {
		return nil, fmt.Errorf("error finding invitations: %w", err)
	}

	fmt.Printf("Found invitations: %+v\n", invitations)
	return invitations, nil
}

func (s *ProjectService) AcceptInvitation(projectID, userID uint) error {
	return s.db.Model(&models.Member{}).
		Where("project_id = ? AND user_id = ? AND status = ?", projectID, userID, models.MemberStatusPending).
		Update("status", models.MemberStatusAccepted).Error
}

func (s *ProjectService) RejectInvitation(projectID, userID uint) error {
	return s.db.Model(&models.Member{}).
		Where("project_id = ? AND user_id = ? AND status = ?", projectID, userID, models.MemberStatusPending).
		Update("status", models.MemberStatusRejected).Error
} 