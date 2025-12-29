package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type Question struct {
	ID        uuid.UUID `gorm:"primaryKey;column:id"`
	ModuleID  uuid.UUID `gorm:"column:module_id"`
	Content   string    `gorm:"column:content"`
	Slug      string    `gorm:"column:slug"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt null.Time `gorm:"nullable;column:deleted_at"`

	Choices []*QuestionChoice `gorm:"foreignKey:QuestionID;references:ID"`
}
