package models

import (
	"time"

	"github.com/google/uuid"
)

// ParkingSlot model
type ParkingSlot struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ZoneID    uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_zone_slot"`
	Name      string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_zone_slot" json:"name" validate:"required,max=50"`
	Status    string    `gorm:"type:varchar(20);not null;default:'available'" json:"status" validate:"required,oneof=available occupied reserved"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	ParkingZone *ParkingZone `gorm:"foreignKey:ZoneID;constraint:OnDelete:CASCADE;" json:"parking_zone,omitempty"`
	Sensors     []Sensor     `gorm:"foreignKey:SlotID;constraint:OnDelete:CASCADE;" json:"sensors,omitempty"`
	Bookings    []Booking    `gorm:"foreignKey:SlotID;constraint:OnDelete:CASCADE;" json:"bookings,omitempty"`
	Reports     []Report     `gorm:"foreignKey:SlotID;constraint:OnDelete:CASCADE;" json:"reports,omitempty"`
}

type ParkingSlotResponse struct {
	ID     uuid.UUID `json:"id"`
	ZoneID uuid.UUID `json:"zone_id"`
	Name   string    `json:"name"`
	Status string    `json:"status"`
}
