package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"primaryKey;autoIncrement;column:id"`
	Email     string    `gorm:"column:email"`
	Password  string    `gorm:"column:password"`
	Fullname  string    `gorm:"column:fullname"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	Sessions []*Session `gorm:"foreignKey:UserId;references:ID"`
}
