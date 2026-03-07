package models

import "gorm.io/gorm"

// Enums for status fields
type ReportState string
type BookingStatus string

const (
	ReportPending  ReportState = "PENDING"
	ReportApproved ReportState = "APPROVED"
	ReportRejected ReportState = "REJECTED"

	BookingPending   BookingStatus = "PENDING"
	BookingConfirmed BookingStatus = "CONFIRMED"
	BookingCheckedIn BookingStatus = "CHECKED_IN"
	BookingCompleted BookingStatus = "COMPLETED"
	BookingCancelled BookingStatus = "CANCELLED"
	BookingExpired   BookingStatus = "EXPIRED"
)

// AutoMigrate will migrate all models in the package
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Profile{},
		&Parking{},
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
