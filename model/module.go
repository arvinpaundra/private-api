package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type ModuleType string

const (
	MultipleChoice ModuleType = "multiple_choice"
	MatchingType   ModuleType = "matching_type"
)

type Module struct {
	ID          uuid.UUID   `gorm:"primaryKey;column:id"`
	UserID      uuid.UUID   `gorm:"column:user_id"`
	SubjectID   uuid.UUID   `gorm:"column:subject_id"`
	GradeID     uuid.UUID   `gorm:"column:grade_id"`
	Title       string      `gorm:"column:title"`
	Slug        string      `gorm:"column:slug"`
	Description null.String `gorm:"nullable;column:description"`
	Type        ModuleType  `gorm:"type:module_type;column:type"`
	IsPublished bool        `gorm:"column:is_published"`
	CreatedAt   time.Time   `gorm:"column:created_at"`
	UpdatedAt   time.Time   `gorm:"column:updated_at"`
	DeletedAt   null.Time   `gorm:"nullable;column:deleted_at"`

	Subject   *Subject    `gorm:"foreignKey:SubjectID;references:ID"`
	Grade     *Grade      `gorm:"foreignKey:GradeID;references:ID"`
	Questions []*Question `gorm:"foreignKey:ModuleID;references:ID"`
}
