package models

import "gorm.io/gorm"

// Enums for status fields
type ReportState string
type BookingStatus string

const (
	ReportPending  ReportState = "PENDING"
	ReportApproved ReportState = "APPROVED"
	ReportRejected ReportState = "REJECTED"

	BookingPending   BookingStatus = "PENDING"   // จองอยู่ (Booking in progress)
	BookingArrived   BookingStatus = "ARRIVED"   // ถึงแล้วรอจ่ายตัง (Arrived and waiting to pay)
	BookingCompleted BookingStatus = "COMPLETED" // เสร็จสิ้น (Completed)
	BookingCancelled BookingStatus = "CANCELLED" // ยกเลิก (Cancelled)
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
