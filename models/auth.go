package models

import (
	"time"

	"github.com/google/uuid"
)

// Auth model
type Auth struct {
	UID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"uid"`
	DisplayName string    `gorm:"type:varchar(255);not null" json:"display_name" validate:"required,min=2,max=255"`
	Email       string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"email" validate:"required,email"`
	Phone       string    `gorm:"type:varchar(20);uniqueIndex;default:null" json:"phone" validate:"omitempty,e164"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Profile       *Profile       `gorm:"foreignKey:UID;references:UID" json:"profile,omitempty"`
	Bookings      []Booking      `gorm:"foreignKey:UserID" json:"bookings,omitempty"`
	Reports       []Report       `gorm:"foreignKey:UserID" json:"reports,omitempty"`
	CodeRedeems   []CodeRedeem   `gorm:"foreignKey:UserID" json:"code_redeems,omitempty"`
	Rewards       []Reward       `gorm:"foreignKey:UserID" json:"rewards,omitempty"`
	RewardRedeems []RewardRedeem `gorm:"foreignKey:UserID" json:"reward_redeems,omitempty"`
}
