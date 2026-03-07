package models

import (
	"time"

	"github.com/google/uuid"
)

// PlaceImage model
type PlaceImage struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ParkingID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_place_path" json:"parking_id"`
	Path      string    `gorm:"type:varchar(500);not null;uniqueIndex:idx_place_path" json:"path" validate:"required,url,max=500"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationship
	Parking *Parking `gorm:"foreignKey:ParkingID;references:ID" json:"parking,omitempty"`
}
