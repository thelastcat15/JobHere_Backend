package models

import (
	"time"

	"github.com/google/uuid"
)

// Profile model
type Profile struct {
	UID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"uid"`
	Point     int       `gorm:"type:integer;not null;default:0;check:point>=0" json:"point"`
	Phone     string    `gorm:"type:varchar(20);default:null" json:"phone" validate:"omitempty,e164"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationship
	Auth *Auth `gorm:"foreignKey:UID;references:UID" json:"auth,omitempty"`
}
