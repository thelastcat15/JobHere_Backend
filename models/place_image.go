package models

import (
	"time"

	"github.com/google/uuid"
)

// PlaceImage model
type PlaceImage struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PlaceID   uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_place_path" json:"place_id"`
	Path      string    `gorm:"type:varchar(500);not null;uniqueIndex:idx_place_path" json:"path" validate:"required,url,max=500"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationship
	Place *Place `gorm:"foreignKey:PlaceID;references:PlaceID" json:"place,omitempty"`
}
