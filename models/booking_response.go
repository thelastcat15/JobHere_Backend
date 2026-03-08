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

	HourlyRate    float64 `json:"hourly_rate"`
	DurationHours float64 `json:"duration_hours,omitempty"`
	TotalCost     float64 `json:"total_cost,omitempty"`

	Parking ParkingInfo `json:"parking"`
	Zone    ZoneInfo    `json:"zone"`
	Slot    SlotInfo    `json:"slot"`
}

type ParkingInfo struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ZoneInfo struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	HourRate float64   `json:"hour_rate"`
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

	ParkingID   uuid.UUID
	ParkingName string

	ZoneID       uuid.UUID
	ZoneName     string
	ZoneHourRate float64

	SlotID   uuid.UUID
	SlotName string

	HourlyRate    float64
	DurationHours float64
	TotalCost     float64
}

func MapBookingRowToResponse(r BookingRow) BookingResponse {
	return BookingResponse{
		ID:              r.ID,
		Status:          r.Status,
		BookedTimeStart: r.BookedTimeStart,
		BookedTimeEnd:   r.BookedTimeEnd,

		Parking: ParkingInfo{
			ID:   r.ParkingID,
			Name: r.ParkingName,
		},

		Zone: ZoneInfo{
			ID:       r.ZoneID,
			Name:     r.ZoneName,
			HourRate: r.ZoneHourRate,
		},

		Slot: SlotInfo{
			ID:   r.SlotID,
			Name: r.SlotName,
		},

		HourlyRate:    r.HourlyRate,
		DurationHours: r.DurationHours,
		TotalCost:     r.TotalCost,
	}
}
