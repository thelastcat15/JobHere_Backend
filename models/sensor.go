package models

import (
	"time"

	"github.com/google/uuid"
)

// Sensor model
type Sensor struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name" validate:"required,max=100"`
	SlotID    uuid.UUID `gorm:"type:uuid;not null;index" json:"slot_id"`
	URL       string    `gorm:"type:varchar(500);not null;uniqueIndex" json:"url" validate:"required,url,max=500"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationship
	ParkingSlot *ParkingSlot `gorm:"foreignKey:SlotID;references:SlotID" json:"parking_slot,omitempty"`
}
