package models

import (
	"time"

	"github.com/google/uuid"
)

// ParkingZone model
type ParkingZone struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ParkingID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_place_zone_name" json:"parking_id"`
	HourRate  float64   `gorm:"type:numeric(10,2);not null;check:hour_rate > 0" json:"hour_rate" validate:"required,gt=0"`
	Name      string    `gorm:"type:varchar(100);not null;uniqueIndex:idx_place_zone_name" json:"name" validate:"required,max=100"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Parking      *Parking      `gorm:"foreignKey:ParkingID;references:ID;constraint:OnDelete:CASCADE;" json:"parking,omitempty"`
	ParkingSlots []ParkingSlot `gorm:"foreignKey:ZoneID;constraint:OnDelete:CASCADE;" json:"parking_slots,omitempty"`
	Bookings     []Booking     `gorm:"foreignKey:ZoneID" json:"bookings,omitempty"`
	Reports      []Report      `gorm:"foreignKey:ZoneID" json:"reports,omitempty"`
}

type ZoneResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
