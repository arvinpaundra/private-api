package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type Session struct {
	ID           uuid.UUID `gorm:"primaryKey;column:id"`
	UserId       uuid.UUID `gorm:"column:user_id"`
	AccessToken  string    `gorm:"column:access_token"`
	RefreshToken string    `gorm:"column:refresh_token"`
	RevokedAt    null.Time `gorm:"nullable;column:revoked_at"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}
