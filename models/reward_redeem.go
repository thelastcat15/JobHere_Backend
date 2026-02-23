package models

import (
	"time"

	"github.com/google/uuid"
)

// RewardRedeem model
type RewardRedeem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	RewardID  uuid.UUID `gorm:"type:uuid;not null;index" json:"reward_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Reward *Reward `gorm:"foreignKey:RewardID;references:ID" json:"reward,omitempty"`
	User   *Auth   `gorm:"foreignKey:UserID;references:UID" json:"user,omitempty"`
}
