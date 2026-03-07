package models

import (
	"time"

	"github.com/google/uuid"
)

type BookingResponse struct {
	ID uuid.UUID `json:"id"`

	Status BookingStatus `json:"status"`

	BookedTimeStart time.Time `json:"booked_time_start"`
	BookedTimeEnd   time.Time `json:"booked_time_end"`

	CheckinTime  *time.Time `json:"checkin_time,omitempty"`
	CheckoutTime *time.Time `json:"checkout_time,omitempty"`

	Parking ParkingInfo `json:"parking"`
	Zone    ZoneInfo    `json:"zone"`
	Slot    SlotInfo    `json:"slot"`
}

type ParkingInfo struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ZoneInfo struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type SlotInfo struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type BookingRow struct {
	ID uuid.UUID

	Status BookingStatus

	BookedTimeStart time.Time
	BookedTimeEnd   time.Time

	CheckinTime  *time.Time
	CheckoutTime *time.Time

	ParkingID   uuid.UUID
	ParkingName string

	ZoneID   uuid.UUID
	ZoneName string

	SlotID   uuid.UUID
	SlotName string
}

func MapBookingRowToResponse(r BookingRow) BookingResponse {
	return BookingResponse{
		ID:              r.ID,
		Status:          r.Status,
		BookedTimeStart: r.BookedTimeStart,
		BookedTimeEnd:   r.BookedTimeEnd,
		CheckinTime:     r.CheckinTime,
		CheckoutTime:    r.CheckoutTime,

		Parking: ParkingInfo{
			ID:   r.ParkingID,
			Name: r.ParkingName,
		},

		Zone: ZoneInfo{
			ID:   r.ZoneID,
			Name: r.ZoneName,
		},

		Slot: SlotInfo{
			ID:   r.SlotID,
			Name: r.SlotName,
		},
	}
}
