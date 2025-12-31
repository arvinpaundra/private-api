package model

import (
	"time"

	"github.com/google/uuid"
)

type SubmissionAnswer struct {
	ID           uuid.UUID `gorm:"primaryKey;column:id"`
	SubmissionID uuid.UUID `gorm:"column:submission_id"`
	QuestionSlug string    `gorm:"column:question_slug"`
	Question     string    `gorm:"column:question"`
	Answer       string    `gorm:"column:answer"`
	IsCorrect    bool      `gorm:"column:is_correct"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}
