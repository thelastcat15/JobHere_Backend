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

// CreatePlace godoc
// @Summary Create Place
// @Description Create a new place
// @Tags Place
// @Accept json
// @Produce json
// @Param place body models.Place true "Place payload"
// @Success 201 {object} models.Place
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/places [post]
// func CreatePlace(c *fiber.Ctx) error {
// 	var place models.Place

// 	if err := c.BodyParser(&place); err != nil {
// 		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
// 	}

// 	if validationErrs := utils.ValidateStruct(place); len(validationErrs) > 0 {
// 		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", validationErrs)
// 	}

// 	if place.ID == uuid.Nil {
// 		place.ID = uuid.New()
// 	}

// 	result := config.DB.Create(&place)
// 	if result.Error != nil {
// 		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create place", result.Error.Error())
// 	}

// 	return utils.SuccessResponse(c, fiber.StatusCreated, "Place created successfully", place)
// }

// GetPlace godoc
// @Summary Get place by ID
// @Description Retrieve a place by UUID
// @Tags Place
// @Accept json
// @Produce json
// @Param id path string true "Place UUID"
// @Success 200 {object} models.PlaceDetailResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/places/{place_id} [get]
func GetPlace(c *fiber.Ctx) error {
	place_id := c.Params("place_id")

	placeID, err := uuid.Parse(place_id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID format", nil)
	}

	var place models.Place

	result := config.DB.
		Preload("Images").
		Preload("ParkingZones", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, place_id")
		}).
		First(&place, "id = ?", placeID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Place not found", nil)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve place", result.Error.Error())
	}

	// map response
	response := models.PlaceDetailResponse{
		ID:          place.ID,
		Type:        place.Type,
		Contact:     place.Contact,
		Address:     place.Address,
		Description: place.Description,
		CoordinateX: place.CoordinateX,
		CoordinateY: place.CoordinateY,
		Images:      place.Images,
	}

	for _, zone := range place.ParkingZones {
		response.Zones = append(response.Zones, models.ZoneResponse{
			ID:   zone.ID,
			Name: zone.Name,
		})
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Place retrieved successfully", response)
}

// ListPlaces godoc
// @Summary List places
// @Description Retrieve all places
// @Tags Place
// @Accept json
// @Produce json
// @Success 200 {array} models.PlaceResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/places [get]
func ListPlaces(c *fiber.Ctx) error {
	var places []models.PlaceResponse

	err := config.DB.
		Table("places p").
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
		Joins("LEFT JOIN parking_zones z ON z.place_id = p.id").
		Joins("LEFT JOIN parking_slots ps ON ps.zone_id = z.id").
		Group("p.id").
		Order("p.created_at DESC").
		Scan(&places).Error

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve places", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Places retrieved successfully", places)
}

// UpdatePlace godoc
// @Summary Update a place
// @Description Update a place by UUID
// @Tags Place
// @Accept json
// @Produce json
// @Param id path string true "Place UUID"
// @Param place body models.Place true "Place payload"
// @Success 200 {object} models.Place
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/places/{id} [put]
// func UpdatePlace(c *fiber.Ctx) error {
// 	id := c.Params("id")

// 	placeID, err := uuid.Parse(id)
// 	if err != nil {
// 		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID format", nil)
// 	}

// 	var updateData models.Place

// 	if err := c.BodyParser(&updateData); err != nil {
// 		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
// 	}

// 	var place models.Place

// 	result := config.DB.First(&place, "id = ?", placeID)
// 	if result.Error != nil {
// 		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 			return utils.ErrorResponse(c, fiber.StatusNotFound, "Place not found", nil)
// 		}
// 		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve place", result.Error.Error())
// 	}

// 	if err := config.DB.Model(&place).Updates(updateData).Error; err != nil {
// 		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update place", err.Error())
// 	}

// 	return utils.SuccessResponse(c, fiber.StatusOK, "Place updated successfully", place)
// }

// DeletePlace godoc
// @Summary Delete a place
// @Description Delete a place by UUID
// @Tags Place
// @Accept json
// @Produce json
// @Param id path string true "Place UUID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/places/{id} [delete]
// func DeletePlace(c *fiber.Ctx) error {
// 	id := c.Params("id")

// 	placeID, err := uuid.Parse(id)
// 	if err != nil {
// 		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID format", nil)
// 	}

// 	result := config.DB.Delete(&models.Place{}, "id = ?", placeID)
// 	if result.Error != nil {
// 		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete place", result.Error.Error())
// 	}

// 	if result.RowsAffected == 0 {
// 		return utils.ErrorResponse(c, fiber.StatusNotFound, "Place not found", nil)
// 	}

// 	return utils.SuccessResponse(c, fiber.StatusOK, "Place deleted successfully", nil)
// }

// GetParkingSlots godoc
// @Summary Get parking slots by zone
// @Description Retrieve all parking slots in a zone
// @Tags ParkingSlot
// @Accept json
// @Produce json
// @Param place_id path string true "Place UUID"
// @Param zone_id path string true "Zone UUID"
// @Success 200 {array} models.ParkingSlotResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/places/{place_id}/zones/{zone_id}/slots [get]
func GetParkingSlots(c *fiber.Ctx) error {
	placeID, err := uuid.Parse(c.Params("place_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid place UUID format", nil)
	}

	zoneID, err := uuid.Parse(c.Params("zone_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid zone UUID format", nil)
	}

	var zone models.ParkingZone
	err = config.DB.
		Where("id = ? AND place_id = ?", zoneID, placeID).
		First(&zone).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Zone not found in this place", nil)
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
