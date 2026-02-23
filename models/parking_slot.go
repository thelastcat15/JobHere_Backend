package models

import (
	"time"

	"github.com/google/uuid"
)

// ParkingSlot model
type ParkingSlot struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ZoneID    uuid.UUID `gorm:"type:uuid;not null;index" json:"zone_id"`
	Name      string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_zone_slot" json:"name" validate:"required,max=50"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	ParkingZone *ParkingZone `gorm:"foreignKey:ZoneID;references:ID" json:"parking_zone,omitempty"`
	Sensors     []Sensor     `gorm:"foreignKey:SlotID;references:ID" json:"sensors,omitempty"`
	Bookings    []Booking    `gorm:"foreignKey:SlotID;references:ID" json:"bookings,omitempty"`
	Reports     []Report     `gorm:"foreignKey:SlotID;references:ID" json:"reports,omitempty"`
}
