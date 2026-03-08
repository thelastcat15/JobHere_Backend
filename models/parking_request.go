package models

type CreateParkingRequest struct {
	Name        string  `json:"name" validate:"required,max=150"`
	Type        string  `json:"type" validate:"required,max=50"`
	Contact     string  `json:"contact" validate:"omitempty,max=100"`
	Address     string  `json:"address" validate:"required,min=5,max=500"`
	Description string  `json:"description"`
	CoordinateX float64 `json:"coordinate_x" validate:"required,latitude"`
	CoordinateY float64 `json:"coordinate_y" validate:"required,longitude"`
}

type UpdateParkingRequest struct {
	Name        *string  `json:"name" validate:"omitempty,max=150"`
	Type        *string  `json:"type" validate:"omitempty,max=50"`
	Contact     *string  `json:"contact" validate:"omitempty,max=100"`
	Address     *string  `json:"address" validate:"omitempty,min=5,max=500"`
	Description *string  `json:"description"`
	CoordinateX *float64 `json:"coordinate_x" validate:"omitempty,latitude"`
	CoordinateY *float64 `json:"coordinate_y" validate:"omitempty,longitude"`
}
