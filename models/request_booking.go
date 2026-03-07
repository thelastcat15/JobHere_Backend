package models

import (
	"time"

	"github.com/google/uuid"
)

type BookingRequest struct {
	ParkingID uuid.UUID `json:"parking_id" validate:"required"`
	ZoneID    uuid.UUID `json:"zone_id" validate:"required"`
	SlotID    uuid.UUID `json:"slot_id" validate:"required"`

	BookedTimeStart time.Time `json:"booked_time_start" validate:"required"`
}
