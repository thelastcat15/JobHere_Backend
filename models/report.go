package models

import (
	"time"

	"github.com/google/uuid"
)

// Report model
type Report struct {
	ID        uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID   `gorm:"type:uuid;not null;index" json:"user_id"`
	Image     string      `gorm:"type:varchar(500);default:null" json:"image" validate:"omitempty,url,max=500"`
	Content   string      `gorm:"type:text;not null" json:"content" validate:"required,min=10,max=2000"`
	Timestamp time.Time   `gorm:"not null;default:CURRENT_TIMESTAMP" json:"timestamp"`
	State     ReportState `gorm:"type:varchar(50);not null;default:'PENDING'" json:"state" validate:"required,oneof=PENDING APPROVED REJECTED"`
	PlaceID   uuid.UUID   `gorm:"type:uuid;not null;index" json:"place_id"`
	ZoneID    uuid.UUID   `gorm:"type:uuid;index;default:null" json:"zone_id"`
	SlotID    uuid.UUID   `gorm:"type:uuid;index;default:null" json:"slot_id"`
	CreatedAt time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time   `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User        *Auth        `gorm:"foreignKey:UserID;references:UID" json:"user,omitempty"`
	Place       *Place       `gorm:"foreignKey:PlaceID;references:ID" json:"place,omitempty"`
	ParkingZone *ParkingZone `gorm:"foreignKey:ZoneID;references:ID" json:"parking_zone,omitempty"`
	ParkingSlot *ParkingSlot `gorm:"foreignKey:SlotID;references:ID" json:"parking_slot,omitempty"`
}
