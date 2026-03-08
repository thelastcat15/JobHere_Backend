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

// GetBookings godoc
// @Summary Get user bookings
// @Description Retrieve all bookings for the authenticated user
// @Tags Booking
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.BookingResponse}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/bookings [get]
func GetBookings(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user id", nil)
	}

	var rows []models.BookingRow

	err = config.DB.Table("bookings").
		Select(`
			bookings.id,
			bookings.status,
			bookings.booked_time_start,
			bookings.booked_time_end,
			bookings.hourly_rate,
			bookings.duration_hours,
			bookings.total_cost,

			parkings.id as parking_id,
			parkings.name as parking_name,

			parking_zones.id as zone_id,
			parking_zones.name as zone_name,
			parking_zones.hour_rate as zone_hour_rate,

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
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve bookings", err.Error())
	}

	var response []models.BookingResponse

	for _, r := range rows {
		response = append(response, models.MapBookingRowToResponse(r))
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Bookings retrieved successfully", response)
}

// GetBooking godoc
// @Summary Get booking
// @Description Retrieve booking by ID (owner only)
// @Tags Booking
// @Accept json
// @Produce json
// @Param id path string true "Booking UUID"
// @Success 200 {object} utils.Response{data=models.BookingResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/bookings/{id} [get]
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
			bookings.hourly_rate,
			bookings.duration_hours,
			bookings.total_cost,

			parkings.id as parking_id,
			parkings.name as parking_name,

			parking_zones.id as zone_id,
			parking_zones.name as zone_name,
			parking_zones.hour_rate as zone_hour_rate,

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

// CreateBooking godoc
// @Summary Create booking
// @Description Create a new parking slot booking (hourly rate charging starts immediately)
// @Tags Booking
// @Accept json
// @Produce json
// @Param booking body models.BookingRequest true "Booking request"
// @Success 201 {object} utils.Response{data=models.BookingResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/bookings [post]
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

	// Fetch zone to get hourly rate
	var zone models.ParkingZone
	if err := tx.First(&zone, "id = ?", req.ZoneID).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Zone not found", nil)
	}

	// Set BookedTimeStart to now (timer starts immediately)
	bookedTimeStart := time.Now()

	booking := models.Booking{
		UserID:          userID,
		ParkingID:       req.ParkingID,
		ZoneID:          req.ZoneID,
		SlotID:          req.SlotID,
		Status:          models.BookingPending,
		BookedTimeStart: bookedTimeStart,
		HourlyRate:      zone.HourRate,
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

	// Fetch parking info for response
	var parking models.Parking
	if err := config.DB.First(&parking, "id = ?", req.ParkingID).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve parking info", err.Error())
	}

	response := models.BookingResponse{
		ID:              booking.ID,
		Status:          booking.Status,
		BookedTimeStart: booking.BookedTimeStart,
		BookedTimeEnd:   booking.BookedTimeEnd,
		HourlyRate:      booking.HourlyRate,
		Parking: models.ParkingInfo{
			ID:   parking.ID,
			Name: parking.Name,
		},
		Zone: models.ZoneInfo{
			ID:       zone.ID,
			Name:     zone.Name,
			HourRate: zone.HourRate,
		},
		Slot: models.SlotInfo{
			ID:   slot.ID,
			Name: slot.Name,
		},
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Booking created - hourly charging started", response)
}

// DeleteBooking godoc
// @Summary Delete booking
// @Description Delete booking by ID (owner only)
// @Tags Booking
// @Accept json
// @Produce json
// @Param id path string true "Booking UUID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/bookings/{id} [delete]
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

// ConfirmUserAtParkingForSuccessBooking godoc
// @Summary Confirm user arrival at parking
// @Description Confirm that user has arrived at parking location and calculate parking cost
// @Tags Booking
// @Accept json
// @Produce json
// @Param id path string true "Booking UUID"
// @Success 200 {object} utils.Response{data=models.Booking}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/bookings/{id}/checkin [post]
func ConfirmUserAtParkingForSuccessBooking(c *fiber.Ctx) error {

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
	if booking.Status != models.BookingPending {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Booking status conflicting", nil)
	}

	now := time.Now()

	booking.Status = models.BookingArrived

	// Calculate duration and cost
	durationHours := now.Sub(booking.BookedTimeStart).Hours()
	totalCost := durationHours * booking.HourlyRate

	booking.DurationHours = durationHours
	booking.TotalCost = totalCost

	if err := config.DB.Save(&booking).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Checkin failed", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Checkin successful", booking)
}

// CompleteBookingPayment godoc
// @Summary Complete booking payment
// @Description Complete payment for arrived booking and mark as completed, release parking slot
// @Tags Booking
// @Accept json
// @Produce json
// @Param id path string true "Booking UUID"
// @Success 200 {object} utils.Response{data=models.Booking}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/bookings/{id}/complete [post]
func CompleteBookingPayment(c *fiber.Ctx) error {

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

	// check status - must be ARRIVED (waiting for payment)
	if booking.Status != models.BookingArrived {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Booking not ready for payment completion", nil)
	}

	tx := config.DB.Begin()

	// update booking status to completed
	booking.Status = models.BookingCompleted

	if err := tx.Save(&booking).Error; err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Payment completion failed", err.Error())
	}

	// release parking slot
	err = tx.Model(&models.ParkingSlot{}).
		Where("id = ?", booking.SlotID).
		Update("status", "available").Error

	if err != nil {
		tx.Rollback()
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to release parking slot", err.Error())
	}

	tx.Commit()

	return utils.SuccessResponse(c, fiber.StatusOK, "Payment completed successfully", booking)
}

// CancelBooking godoc
// @Summary Cancel booking
// @Description Cancel pending booking and release slot
// @Tags Booking
// @Accept json
// @Produce json
// @Param id path string true "Booking UUID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/bookings/{id}/cancel [post]
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
	if booking.Status != models.BookingPending {
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

// GetUserBookingHistory godoc
// @Summary Get user booking history
// @Description Retrieve booking history for current user
// @Tags Booking
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.BookingResponse}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/bookings/history [get]
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
			bookings.hourly_rate,
			bookings.duration_hours,
			bookings.total_cost,

			parkings.id as parking_id,
			parkings.name as parking_name,

			parking_zones.id as zone_id,
			parking_zones.name as zone_name,
			parking_zones.hour_rate as zone_hour_rate,

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
