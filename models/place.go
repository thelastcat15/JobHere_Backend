package models

import (
	"time"

	"github.com/google/uuid"
)

// Place model
type Place struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Type        string    `gorm:"type:varchar(50);not null" json:"type" validate:"required,max=50"`
	Contact     string    `gorm:"type:varchar(100);default:null" json:"contact" validate:"omitempty,max=100"`
	Address     string    `gorm:"type:text;not null" json:"address" validate:"required,min=5,max=500"`
	Description string    `gorm:"type:text;default:null" json:"description"`
	CoordinateX float64   `gorm:"type:numeric(10,6);not null" json:"coordinate_x" validate:"required,latitude"`
	CoordinateY float64   `gorm:"type:numeric(10,6);not null" json:"coordinate_y" validate:"required,longitude"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	ParkingZones []ParkingZone `gorm:"foreignKey:PlaceID;references:ID" json:"parking_zones,omitempty"`
	Bookings     []Booking     `gorm:"foreignKey:PlaceID;references:ID" json:"bookings,omitempty"`
	Reports      []Report      `gorm:"foreignKey:PlaceID;references:ID" json:"reports,omitempty"`
	Images       []PlaceImage  `gorm:"foreignKey:PlaceID;references:ID" json:"images,omitempty"`
}
