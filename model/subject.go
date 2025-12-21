package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type Subject struct {
	ID          uuid.UUID   `gorm:"primaryKey;column:id"`
	UserID      uuid.UUID   `gorm:"column:user_id"`
	Name        string      `gorm:"column:name"`
	Description null.String `gorm:"nullable;column:description"`
	CreatedAt   time.Time   `gorm:"column:created_at"`
	UpdatedAt   time.Time   `gorm:"column:updated_at"`
	DeletedAt   null.Time   `gorm:"nullable;column:deleted_at"`
}
