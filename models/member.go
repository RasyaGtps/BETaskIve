package models

type MemberRole string

const (
	MemberRoleOwner  MemberRole = "OWNER"
	MemberRoleEditor MemberRole = "EDITOR"
	MemberRoleViewer MemberRole = "VIEWER"
)

type Member struct {
	ProjectID uint       `gorm:"primarykey" json:"project_id"`
	Project   Project    `gorm:"foreignKey:ProjectID" json:"-"`
	UserID    uint       `gorm:"primarykey" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID" json:"user"`
	Role      MemberRole `gorm:"type:varchar(20);not null" json:"role" validate:"required,oneof=OWNER EDITOR VIEWER"`
}

func (Member) TableName() string {
	return "project_members"
} 