package models

import (
	"time"

	"github.com/google/uuid"
)

// Profile model
type Profile struct {
	UID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"uid"`
	Email       string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email" validate:"required,email"`
	DisplayName string    `gorm:"type:varchar(100);not null" json:"display_name" validate:"required"`
	Point       int       `gorm:"type:integer;not null;default:0;check:point>=0" json:"point"`
	Phone       string    `gorm:"type:varchar(20);default:null" json:"phone" validate:"omitempty,e164"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
