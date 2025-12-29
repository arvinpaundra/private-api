package model

import (
	"time"

	"github.com/guregu/null/v6"
)

type Question struct {
	ID        string    `gorm:"primaryKey;column:id"`
	ModuleID  string    `gorm:"column:module_id"`
	Title     string    `gorm:"column:title"`
	Slug      string    `gorm:"column:slug"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt null.Time `gorm:"nullable;column:deleted_at"`

	Choices []*QuestionChoice `gorm:"foreignKey:QuestionID;references:ID"`
}
