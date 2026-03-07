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
	Role        string    `gorm:"type:varchar(20);not null;default:'user'" json:"role" validate:"oneof=user admin guard"`
	Point       int       `gorm:"type:integer;not null;default:0;check:point>=0" json:"point"`
	Phone       string    `gorm:"type:varchar(20);default:null" json:"phone" validate:"omitempty,e164"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type UpdateProfileRequest struct {
	DisplayName string `json:"display_name" validate:"required,min=2,max=100"`
	Phone       string `json:"phone,omitempty" validate:"omitempty,e164,max=20"`
}
