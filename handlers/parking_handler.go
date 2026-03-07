package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"jodhere.backend/config"
	"jodhere.backend/models"
	"jodhere.backend/utils"
)

// CreateParking godoc
// @Summary Create Parking
// @Description Create a new parking
// @Tags Parking
// @Accept json
// @Produce json
// @Param parking body models.Parking true "Parking payload"
// @Success 201 {object} models.Parking
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/parking [post]
// func CreateParking(c *fiber.Ctx) error {
// 	var parking models.Parking

// 	if err := c.BodyParser(&parking); err != nil {
// 		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
// 	}

// 	if validationErrs := utils.ValidateStruct(parking); len(validationErrs) > 0 {
// 		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", validationErrs)
// 	}

// 	if parking.ID == uuid.Nil {
// 		parking.ID = uuid.New()
// 	}

// 	result := config.DB.Create(&parking)
// 	if result.Error != nil {
// 		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create parking", result.Error.Error())
// 	}

// 	return utils.SuccessResponse(c, fiber.StatusCreated, "Parking created successfully", parking)
// }

// GetParking godoc
// @Summary Get parking by ID
// @Description Retrieve a parking by UUID
// @Tags Parking
// @Accept json
// @Produce json
// @Param id path string true "Parking UUID"
// @Success 200 {object} models.ParkingDetailResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/parking/{parking_id} [get]
func GetParking(c *fiber.Ctx) error {
	parking_id := c.Params("parking_id")

	parkingID, err := uuid.Parse(parking_id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID format", nil)
	}

	var parking models.Parking

	result := config.DB.
		Preload("Images").
		Preload("ParkingZones", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, parking_id")
		}).
		First(&parking, "id = ?", parkingID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Parking not found", nil)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve parking", result.Error.Error())
	}

	// map response
	response := models.ParkingDetailResponse{
		ID:          parking.ID,
		Type:        parking.Type,
		Contact:     parking.Contact,
		Address:     parking.Address,
		Description: parking.Description,
		CoordinateX: parking.CoordinateX,
		CoordinateY: parking.CoordinateY,
		Images:      parking.Images,
	}

	for _, zone := range parking.ParkingZones {
		response.Zones = append(response.Zones, models.ZoneResponse{
			ID:   zone.ID,
			Name: zone.Name,
		})
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Parking retrieved successfully", response)
}

// ListParking godoc
// @Summary List parking
// @Description Retrieve all parking
// @Tags Parking
// @Accept json
// @Produce json
// @Success 200 {array} models.ParkingResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/parking [get]
func ListParking(c *fiber.Ctx) error {
	var parking []models.ParkingResponse

	err := config.DB.
		Table("parking p").
		Select(`
			p.id,
			p.type,
			p.contact,
			p.address,
			p.description,
			p.coordinate_x,
			p.coordinate_y,
			COUNT(ps.id) FILTER (WHERE ps.status = 'available') AS available_slots
		`).
		Joins("LEFT JOIN parking_zones z ON z.parking_id = p.id").
		Joins("LEFT JOIN parking_slots ps ON ps.zone_id = z.id").
		Group("p.id").
		Order("p.created_at DESC").
		Scan(&parking).Error

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve parking", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Parking retrieved successfully", parking)
}

// UpdateParking godoc
// @Summary Update a parking
// @Description Update a parking by UUID
// @Tags Parking
// @Accept json
// @Produce json
// @Param id path string true "Parking UUID"
// @Param parking body models.Parking true "Parking payload"
// @Success 200 {object} models.Parking
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/parking/{id} [put]
// func UpdateParking(c *fiber.Ctx) error {
// 	id := c.Params("id")

// 	parkingID, err := uuid.Parse(id)
// 	if err != nil {
// 		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID format", nil)
// 	}

// 	var updateData models.Parking

// 	if err := c.BodyParser(&updateData); err != nil {
// 		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
// 	}

// 	var parking models.Parking

// 	result := config.DB.First(&parking, "id = ?", parkingID)
// 	if result.Error != nil {
// 		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 			return utils.ErrorResponse(c, fiber.StatusNotFound, "Parking not found", nil)
// 		}
// 		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve parking", result.Error.Error())
// 	}

// 	if err := config.DB.Model(&parking).Updates(updateData).Error; err != nil {
// 		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update parking", err.Error())
// 	}

// 	return utils.SuccessResponse(c, fiber.StatusOK, "Parking updated successfully", parking)
// }

// DeleteParking godoc
// @Summary Delete a parking
// @Description Delete a parking by UUID
// @Tags Parking
// @Accept json
// @Produce json
// @Param id path string true "Parking UUID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/parking/{id} [delete]
// func DeleteParking(c *fiber.Ctx) error {
// 	id := c.Params("id")

// 	parkingID, err := uuid.Parse(id)
// 	if err != nil {
// 		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID format", nil)
// 	}

// 	result := config.DB.Delete(&models.Parking{}, "id = ?", parkingID)
// 	if result.Error != nil {
// 		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete parking", result.Error.Error())
// 	}

// 	if result.RowsAffected == 0 {
// 		return utils.ErrorResponse(c, fiber.StatusNotFound, "Parking not found", nil)
// 	}

// 	return utils.SuccessResponse(c, fiber.StatusOK, "Parking deleted successfully", nil)
// }

// GetParkingSlots godoc
// @Summary Get parking slots by zone
// @Description Retrieve all parking slots in a zone
// @Tags ParkingSlot
// @Accept json
// @Produce json
// @Param parking_id path string true "Parking UUID"
// @Param zone_id path string true "Zone UUID"
// @Success 200 {array} models.ParkingSlotResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/parking/{parking_id}/zones/{zone_id}/slots [get]
func GetParkingSlots(c *fiber.Ctx) error {
	parkingID, err := uuid.Parse(c.Params("parking_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid parking UUID format", nil)
	}

	zoneID, err := uuid.Parse(c.Params("zone_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid zone UUID format", nil)
	}

	var zone models.ParkingZone
	err = config.DB.
		Where("id = ? AND parking_id = ?", zoneID, parkingID).
		First(&zone).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Zone not found in this parking", nil)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to verify zone", err.Error())
	}

	// Get slots
	var slots []models.ParkingSlot
	err = config.DB.
		Where("zone_id = ?", zoneID).
		Order("name ASC").
		Find(&slots).Error

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve slots", err.Error())
	}

	var response []models.ParkingSlotResponse
	for _, s := range slots {
		response = append(response, models.ParkingSlotResponse{
			ID:     s.ID,
			ZoneID: s.ZoneID,
			Name:   s.Name,
			Status: s.Status,
		})
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Slots retrieved successfully", response)
}
