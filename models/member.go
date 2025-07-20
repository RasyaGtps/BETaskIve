package models

import (
	"time"

	"gorm.io/gorm"
)

type MemberRole string
type MemberStatus string

const (
	MemberRoleOwner  MemberRole = "OWNER"
	MemberRoleEditor MemberRole = "EDITOR"
	MemberRoleViewer MemberRole = "VIEWER"
)

const (
	MemberStatusPending  MemberStatus = "PENDING"
	MemberStatusAccepted MemberStatus = "ACCEPTED"
	MemberStatusRejected MemberStatus = "REJECTED"
)

type Member struct {
	ProjectID uint         `gorm:"primarykey;column:project_id" json:"project_id"`
	Project   Project      `gorm:"foreignKey:ProjectID" json:"-"`
	UserID    uint         `gorm:"primarykey;column:user_id" json:"user_id"`
	User      User         `gorm:"foreignKey:UserID" json:"user"`
	Role      MemberRole   `gorm:"type:varchar(20);not null;column:role" json:"role" validate:"required,oneof=OWNER EDITOR VIEWER"`
	Status    MemberStatus `gorm:"type:varchar(20);not null;default:PENDING;column:status" json:"status"`
	CreatedAt time.Time    `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time    `gorm:"column:updated_at" json:"updated_at"`
}

func (m *Member) BeforeCreate(tx *gorm.DB) error {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = time.Now()
	}
	return nil
}

func (Member) TableName() string {
	return "project_members"
}