package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"not null" json:"name" validate:"required"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	OwnerID     uint      `json:"owner_id"`
	Owner       User      `gorm:"foreignKey:OwnerID" json:"owner"`
	Tasks       []Task    `json:"tasks,omitempty"`
	Members     []Member  `json:"members,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	if p.StartDate.IsZero() {
		p.StartDate = time.Now()
	}
	return nil
} 