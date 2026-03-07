package models

import (
	"time"

	"github.com/google/uuid"
)

// Booking model (Database)
type Booking struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`

	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	ParkingID uuid.UUID `gorm:"type:uuid;not null;index"`
	ZoneID    uuid.UUID `gorm:"type:uuid;not null;index"`
	SlotID    uuid.UUID `gorm:"type:uuid;not null;index"`

	Status BookingStatus `gorm:"type:varchar(30);not null;default:'PENDING';index;check:status IN ('PENDING','CONFIRMED','CHECKED_IN','COMPLETED','CANCELLED','EXPIRED')"`

	BookedTimeStart time.Time `gorm:"not null"`
	BookedTimeEnd   time.Time `gorm:"not null"`

	CheckinTime  *time.Time
	CheckoutTime *time.Time

	GraceMinutes int `gorm:"default:60"`

	CreatedAt time.Time
	UpdatedAt time.Time

	// Relationships
	User        *Profile     `gorm:"foreignKey:UserID;references:UID;constraint:OnDelete:CASCADE;" json:"user,omitempty"`
	Parking     *Parking     `gorm:"foreignKey:ParkingID;references:ID;constraint:OnDelete:CASCADE;" json:"parking,omitempty"`
	ParkingZone *ParkingZone `gorm:"foreignKey:ZoneID;references:ID;constraint:OnDelete:CASCADE;" json:"parking_zone,omitempty"`
	ParkingSlot *ParkingSlot `gorm:"foreignKey:SlotID;references:ID;constraint:OnDelete:CASCADE;" json:"parking_slot,omitempty"`
}
