package models

import (
	"time"

	"github.com/google/uuid"
)

// CodeRedeem model
type CodeRedeem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Type      uuid.UUID `gorm:"type:uuid;not null" json:"type"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationship
	User *Profile `gorm:"foreignKey:UserID;references:UID" json:"user,omitempty"`
}
