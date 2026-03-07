package handlers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"jodhere.backend/config"
	"jodhere.backend/models"
	"jodhere.backend/utils"
)

// GetBookings
func GetBookings(c *fiber.Ctx) error {
	var rows []models.BookingRow

	err := config.DB.Table("bookings").
		Select(`
			bookings.id,
			bookings.status,
			bookings.booked_time_start,
			bookings.booked_time_end,
			bookings.checkin_time,
			bookings.checkout_time,

			parkings.id as parking_id,
			parkings.name as parking_name,

			parking_zones.id as zone_id,
			parking_zones.name as zone_name,

			parking_slots.id as slot_id,
			parking_slots.name as slot_name
		`).
		Joins("JOIN parkings ON parkings.id = bookings.parking_id").
		Joins("JOIN parking_zones ON parking_zones.id = bookings.zone_id").
		Joins("JOIN parking_slots ON parking_slots.id = bookings.slot_id").
		Order("bookings.created_at DESC").
		Scan(&rows).Error

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve bookings", err.Error())
	}

	var response []models.BookingResponse

	for _, r := range rows {
		response = append(response, models.MapBookingRowToResponse(r))
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Bookings retrieved successfully", response)
}

func GetBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	bookingID, err := uuid.Parse(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID format", nil)
	}

	userIDStr := c.Locals("user_id").(string)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user id", nil)
	}

	var row models.BookingRow

	err = config.DB.Table("bookings").
		Select(`
			bookings.id,
			bookings.status,
			bookings.booked_time_start,
			bookings.booked_time_end,
			bookings.checkin_time,
			bookings.checkout_time,

			parkings.id as parking_id,
			parkings.name as parking_name,

			parking_zones.id as zone_id,
			parking_zones.name as zone_name,

			parking_slots.id as slot_id,
			parking_slots.name as slot_name
		`).
		Joins("JOIN parkings ON parkings.id = bookings.parking_id").
		Joins("JOIN parking_zones ON parking_zones.id = bookings.zone_id").
		Joins("JOIN parking_slots ON parking_slots.id = bookings.slot_id").
		Where("bookings.id = ? AND bookings.user_id = ?", bookingID, userID).
		Scan(&row).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Booking not found", nil)
		}

		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve booking", err.Error())
	}

	response := models.MapBookingRowToResponse(row)

	return utils.SuccessResponse(c, fiber.StatusOK, "Booking retrieved successfully", response)
}

func CreateBooking(c *fiber.Ctx) error {

	var req models.BookingRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// get user id from JWT middleware
	userIDStr := c.Locals("user_id").(string)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user id", nil)
	}

	tx := config.DB.Begin()

	// check slot
	var slot models.ParkingSlot
	if err := tx.First(&slot, "id = ?", req.SlotID).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Slot not found", nil)
	}

	if slot.Status != "available" {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Slot already reserved", nil)
	}

	booking := models.Booking{
		UserID:          userID,
		ParkingID:       req.ParkingID,
		ZoneID:          req.ZoneID,
		SlotID:          req.SlotID,
		Status:          models.BookingPending,
		BookedTimeStart: req.BookedTimeStart,
	}

	if err := tx.Create(&booking).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create booking", err.Error())
	}

	if err := tx.Model(&models.ParkingSlot{}).
		Where("id = ?", req.SlotID).
		Update("status", "reserved").Error; err != nil {

		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update slot status", err.Error())
	}

	tx.Commit()

	return utils.SuccessResponse(c, fiber.StatusCreated, "Booking created", booking)
}

func UpdateBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	bookingID, err := uuid.Parse(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID format", nil)
	}

	userIDStr := c.Locals("user_id").(string)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user id", nil)
	}

	var booking models.Booking

	if err := config.DB.First(&booking, "id = ? AND user_id = ?", bookingID, userID).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Booking not found", nil)
		}

		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to find booking", err.Error())
	}

	var input models.Booking

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if input.Status == models.BookingCancelled || input.Status == models.BookingCompleted {

		err := config.DB.Model(&models.ParkingSlot{}).
			Where("id = ?", booking.SlotID).
			Update("status", "available").Error

		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update slot", err.Error())
		}
	}

	if err := config.DB.Model(&booking).Updates(input).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update booking", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Booking updated successfully", booking)
}

func DeleteBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	bookingID, err := uuid.Parse(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID format", nil)
	}

	userIDStr := c.Locals("user_id").(string)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user id", nil)
	}

	result := config.DB.Delete(&models.Booking{}, "id = ? AND user_id = ?", bookingID, userID)

	if result.Error != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete booking", result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Booking not found", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Booking deleted successfully", nil)
}

// ======================

func CheckinBooking(c *fiber.Ctx) error {

	id := c.Params("id")

	bookingID, err := uuid.Parse(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID", nil)
	}
	var booking models.Booking

	err = config.DB.First(&booking, "id = ?", bookingID).Error
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Booking not found", nil)
	}

	// check status
	if booking.Status != models.BookingConfirmed {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Booking not ready for checkin", nil)
	}

	now := time.Now()

	// check grace time
	expireTime := booking.BookedTimeStart.Add(time.Duration(booking.GraceMinutes) * time.Minute)

	if now.After(expireTime) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Booking expired", nil)
	}

	booking.Status = models.BookingCheckedIn
	booking.CheckinTime = &now

	if err := config.DB.Save(&booking).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Checkin failed", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Checkin successful", booking)
}

func CheckoutBooking(c *fiber.Ctx) error {

	id := c.Params("id")

	bookingID, err := uuid.Parse(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID", nil)
	}

	var booking models.Booking

	err = config.DB.First(&booking, "id = ?", bookingID).Error
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Booking not found", nil)
	}

	if booking.Status != models.BookingCheckedIn {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "User not checked in", nil)
	}

	now := time.Now()

	booking.Status = models.BookingCompleted
	booking.CheckoutTime = &now

	tx := config.DB.Begin()

	if err := tx.Save(&booking).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Checkout failed", err.Error())
	}

	// release slot
	err = tx.Model(&models.ParkingSlot{}).
		Where("id = ?", booking.SlotID).
		Update("status", "available").Error

	if err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Slot update failed", err.Error())
	}

	tx.Commit()

	return utils.SuccessResponse(c, fiber.StatusOK, "Checkout successful", booking)
}

func CancelBooking(c *fiber.Ctx) error {

	id := c.Params("id")

	bookingID, err := uuid.Parse(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid booking UUID", nil)
	}

	// get user id from JWT
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user id", nil)
	}

	var booking models.Booking

	err = config.DB.First(&booking, "id = ?", bookingID).Error
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Booking not found", nil)
		}

		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error", err.Error())
	}

	// verify booking owner
	if booking.UserID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You cannot cancel this booking", nil)
	}

	// check cancellable state
	if booking.Status != models.BookingConfirmed {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Booking cannot be cancelled", nil)
	}

	tx := config.DB.Begin()

	// update booking status
	err = tx.Model(&booking).Update("status", models.BookingCancelled).Error
	if err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to cancel booking", err.Error())
	}

	// release slot
	err = tx.Model(&models.ParkingSlot{}).
		Where("id = ?", booking.SlotID).
		Update("status", "available").Error

	if err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update slot", err.Error())
	}

	tx.Commit()

	return utils.SuccessResponse(c, fiber.StatusOK, "Booking cancelled successfully", nil)
}

func GetUserBookingHistory(c *fiber.Ctx) error {

	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user id", nil)
	}

	var rows []models.BookingRow

	err = config.DB.
		Table("bookings").
		Select(`
			bookings.id,
			bookings.status,
			bookings.booked_time_start,
			bookings.booked_time_end,
			bookings.checkin_time,
			bookings.checkout_time,

			parkings.id as parking_id,
			parkings.name as parking_name,

			parking_zones.id as zone_id,
			parking_zones.name as zone_name,

			parking_slots.id as slot_id,
			parking_slots.name as slot_name
		`).
		Joins("JOIN parkings ON parkings.id = bookings.parking_id").
		Joins("JOIN parking_zones ON parking_zones.id = bookings.zone_id").
		Joins("JOIN parking_slots ON parking_slots.id = bookings.slot_id").
		Where("bookings.user_id = ?", userID).
		Order("bookings.created_at DESC").
		Scan(&rows).Error

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch booking history", err.Error())
	}

	responses := make([]models.BookingResponse, 0, len(rows))

	for _, r := range rows {
		responses = append(responses, models.MapBookingRowToResponse(r))
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Booking history retrieved", responses)
}
