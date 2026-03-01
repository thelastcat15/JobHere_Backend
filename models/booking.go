package models

import (
	"time"

	"github.com/google/uuid"
)

// Booking model
type Booking struct {
	ID              uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID     `gorm:"type:uuid;not null;index" json:"user_id"`
	PlaceID         uuid.UUID     `gorm:"type:uuid;not null;index" json:"place_id"`
	Status          BookingStatus `gorm:"type:varchar(50);not null;default:'CONFIRMED'" json:"status" validate:"required,oneof=CONFIRMED COMPLETED CANCELLED"`
	BookedTimeStart time.Time     `gorm:"not null" json:"booked_time_start" validate:"required"`
	BookedTimeEnd   time.Time     `gorm:"not null" json:"booked_time_end" validate:"required,gtfield=BookedTimeStart"`
	ZoneID          uuid.UUID     `gorm:"type:uuid;not null;index" json:"zone_id"`
	SlotID          uuid.UUID     `gorm:"type:uuid;not null;index" json:"slot_id"`
	CreatedAt       time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time     `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User        *Profile     `gorm:"foreignKey:UserID;references:UID" json:"user,omitempty"`
	Place       *Place       `gorm:"foreignKey:PlaceID;references:ID" json:"place,omitempty"`
	ParkingZone *ParkingZone `gorm:"foreignKey:ZoneID;references:ID" json:"parking_zone,omitempty"`
	ParkingSlot *ParkingSlot `gorm:"foreignKey:SlotID;references:ID" json:"parking_slot,omitempty"`
}
