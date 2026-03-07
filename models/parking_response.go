package models

import (
	"github.com/google/uuid"
)

type ParkingResponse struct {
	ID             uuid.UUID `json:"id"`
	Type           string    `json:"type"`
	Contact        string    `json:"contact"`
	Address        string    `json:"address"`
	Description    string    `json:"description"`
	CoordinateX    float64   `json:"coordinate_x"`
	CoordinateY    float64   `json:"coordinate_y"`
	AvailableSlots int64     `json:"available_slots" gorm:"column:available_slots"`
}

type ParkingDetailResponse struct {
	ID          uuid.UUID      `json:"id"`
	Type        string         `json:"type"`
	Contact     string         `json:"contact"`
	Address     string         `json:"address"`
	Description string         `json:"description"`
	CoordinateX float64        `json:"coordinate_x"`
	CoordinateY float64        `json:"coordinate_y"`
	Images      []PlaceImage   `json:"images"`
	Zones       []ZoneResponse `json:"zones"`
}
