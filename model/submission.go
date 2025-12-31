package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type SubmissionStatus string

const (
	InProgress SubmissionStatus = "inprogress"
	Submitted  SubmissionStatus = "submitted"
	Canceled   SubmissionStatus = "canceled"
)

type Submission struct {
	ID             uuid.UUID        `gorm:"primaryKey;column:id"`
	ModuleID       uuid.UUID        `gorm:"column:module_id"`
	Code           string           `gorm:"column:code"`
	StudentName    string           `gorm:"column:student_name"`
	Status         SubmissionStatus `gorm:"type:submission_status;column:status"`
	TotalQuestions int              `gorm:"column:total_questions"`
	SubmittedAt    null.Time        `gorm:"column:submitted_at"`
	CreatedAt      time.Time        `gorm:"column:created_at"`
	UpdatedAt      time.Time        `gorm:"column:updated_at"`

	Module  *Module             `gorm:"foreignKey:ModuleID;references:ID"`
	Answers []*SubmissionAnswer `gorm:"foreignKey:SubmissionID;references:ID"`
}
