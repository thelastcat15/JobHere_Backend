package models

import "gorm.io/gorm"

// Enums for status fields
type ReportState string
type BookingStatus string

const (
	ReportPending  ReportState = "PENDING"
	ReportApproved ReportState = "APPROVED"
	ReportRejected ReportState = "REJECTED"

	BookingConfirmed BookingStatus = "CONFIRMED"
	BookingCompleted BookingStatus = "COMPLETED"
	BookingCancelled BookingStatus = "CANCELLED"
)

// AutoMigrate will migrate all models in the package
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Profile{},
		&Place{},
		&ParkingZone{},
		&ParkingSlot{},
		&Sensor{},
		&PlaceImage{},
		&Booking{},
		&Report{},
		&CodeRedeem{},
		&Reward{},
		&RewardRedeem{},
	)
}
