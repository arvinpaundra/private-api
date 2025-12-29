package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type QuestionChoice struct {
	ID              uuid.UUID `gorm:"primaryKey;column:id"`
	QuestionID      uuid.UUID `gorm:"column:question_id"`
	Content         string    `gorm:"column:content"`
	IsCorrectAnswer bool      `gorm:"column:is_correct_answer"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
	DeletedAt       null.Time `gorm:"nullable;column:deleted_at"`
}
