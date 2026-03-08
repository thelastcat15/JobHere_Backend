package models

import (
	"github.com/google/uuid"
)

// ParkingResponse represents parking list response with available slots count
type ParkingResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	Contact        string    `json:"contact"`
	Address        string    `json:"address"`
	Description    string    `json:"description"`
	CoordinateX    float64   `json:"coordinate_x"`
	CoordinateY    float64   `json:"coordinate_y"`
	AvailableSlots int       `json:"available_slots" gorm:"column:available_slots"`
}

// ParkingDetailResponse represents detailed parking information with zones and images
type ParkingDetailResponse struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Contact     string       `json:"contact"`
	Address     string       `json:"address"`
	Description string       `json:"description"`
	CoordinateX float64      `json:"coordinate_x"`
	CoordinateY float64      `json:"coordinate_y"`
	Images      []PlaceImage `json:"images"`
	Zones       []ZoneInfo   `json:"zones"`
}
