package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"jobhere.backend/config"
	"jobhere.backend/models"
	"jobhere.backend/utils"
)

// CreateParkingZone godoc
// @Summary Create Parking Zone
// @Description Create a new parking zone
// @Tags ParkingZone
// @Accept json
// @Produce json
// @Param zone body models.ParkingZone true "ParkingZone payload"
// @Success 201 {object} models.ParkingZone
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/parking-zones [post]
func CreateParkingZone(c *fiber.Ctx) error {
	var zone models.ParkingZone

	if err := c.BodyParser(&zone); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if validationErrs := utils.ValidateStruct(zone); len(validationErrs) > 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", validationErrs)
	}

	if zone.ID == uuid.Nil {
		zone.ID = uuid.New()
	}

	result := config.DB.Create(&zone)
	if result.Error != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create parking zone", result.Error.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Parking zone created successfully", zone)
}

// GetParkingZone godoc
// @Summary Get parking zone by ID
// @Description Retrieve a parking zone by UUID
// @Tags ParkingZone
// @Accept json
// @Produce json
// @Param id path string true "ParkingZone UUID"
// @Success 200 {object} models.ParkingZone
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/parking-zones/{id} [get]
func GetParkingZone(c *fiber.Ctx) error {
	id := c.Params("id")

	zoneID, err := uuid.Parse(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID format", nil)
	}

	var zone models.ParkingZone

	result := config.DB.
		Preload("Place").
		Preload("ParkingSlots").
		First(&zone, "id = ?", zoneID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Parking zone not found", nil)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve parking zone", result.Error.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Parking zone retrieved successfully", zone)
}

// ListParkingZones godoc
// @Summary List parking zones
// @Description List parking zones with optional filtering and pagination
// @Tags ParkingZone
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param place_id query string false "Place UUID"
// @Success 200 {array} models.ParkingZone
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/parking-zones [get]
func ListParkingZones(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	placeID := c.Query("place_id", "")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var zones []models.ParkingZone
	var total int64

	query := config.DB.Model(&models.ParkingZone{})

	if placeID != "" {
		query = query.Where("place_id = ?", placeID)
	}

	query.Count(&total)

	result := query.
		Preload("Place").
		Preload("ParkingSlots").
		Offset(offset).
		Limit(pageSize).
		Find(&zones)

	if result.Error != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve parking zones", result.Error.Error())
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	paginatedData := utils.PaginatedData{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		Items:      zones,
	}

	return utils.PaginatedResponse(c, fiber.StatusOK, "Parking zones retrieved successfully", paginatedData)
}

// UpdateParkingZone godoc
// @Summary Update a parking zone
// @Description Update a parking zone by UUID
// @Tags ParkingZone
// @Accept json
// @Produce json
// @Param id path string true "ParkingZone UUID"
// @Param zone body models.ParkingZone true "ParkingZone payload"
// @Success 200 {object} models.ParkingZone
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/parking-zones/{id} [put]
func UpdateParkingZone(c *fiber.Ctx) error {
	id := c.Params("id")

	zoneID, err := uuid.Parse(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID format", nil)
	}

	var updateData models.ParkingZone

	if err := c.BodyParser(&updateData); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	var zone models.ParkingZone

	result := config.DB.First(&zone, "id = ?", zoneID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Parking zone not found", nil)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve parking zone", result.Error.Error())
	}

	if err := config.DB.Model(&zone).Updates(updateData).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update parking zone", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Parking zone updated successfully", zone)
}

// DeleteParkingZone godoc
// @Summary Delete a parking zone
// @Description Delete a parking zone by UUID
// @Tags ParkingZone
// @Accept json
// @Produce json
// @Param id path string true "ParkingZone UUID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/parking-zones/{id} [delete]
func DeleteParkingZone(c *fiber.Ctx) error {
	id := c.Params("id")

	zoneID, err := uuid.Parse(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID format", nil)
	}

	result := config.DB.Delete(&models.ParkingZone{}, "id = ?", zoneID)
	if result.Error != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete parking zone", result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Parking zone not found", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Parking zone deleted successfully", nil)
}
