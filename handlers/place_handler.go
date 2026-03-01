package handlers

import (
	"errors"
	"strconv"

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
func CreatePlace(c *fiber.Ctx) error {
	var place models.Place

	if err := c.BodyParser(&place); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if validationErrs := utils.ValidateStruct(place); len(validationErrs) > 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", validationErrs)
	}

	if place.ID == uuid.Nil {
		place.ID = uuid.New()
	}

	result := config.DB.Create(&place)
	if result.Error != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create place", result.Error.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Place created successfully", place)
}

// GetPlace godoc
// @Summary Get place by ID
// @Description Retrieve a place by UUID
// @Tags Place
// @Accept json
// @Produce json
// @Param id path string true "Place UUID"
// @Success 200 {object} models.Place
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/places/{id} [get]
func GetPlace(c *fiber.Ctx) error {
	id := c.Params("id")

	placeID, err := uuid.Parse(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID format", nil)
	}

	var place models.Place

	result := config.DB.
		Preload("ParkingZones").
		Preload("Images").
		First(&place, "id = ?", placeID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Place not found", nil)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve place", result.Error.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Place retrieved successfully", place)
}

// ListPlaces godoc
// @Summary List places
// @Description List places with optional filtering and pagination
// @Tags Place
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param type query string false "Place type"
// @Success 200 {array} models.Place
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/places [get]
func ListPlaces(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	placeType := c.Query("type", "")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var places []models.Place
	var total int64

	query := config.DB.Model(&models.Place{})

	if placeType != "" {
		query = query.Where("type = ?", placeType)
	}

	query.Count(&total)

	result := query.
		Preload("ParkingZones").
		Offset(offset).
		Limit(pageSize).
		Find(&places)

	if result.Error != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve places", result.Error.Error())
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	paginatedData := utils.PaginatedData{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		Items:      places,
	}

	return utils.PaginatedResponse(c, fiber.StatusOK, "Places retrieved successfully", paginatedData)
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
func UpdatePlace(c *fiber.Ctx) error {
	id := c.Params("id")

	placeID, err := uuid.Parse(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID format", nil)
	}

	var updateData models.Place

	if err := c.BodyParser(&updateData); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	var place models.Place

	result := config.DB.First(&place, "id = ?", placeID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Place not found", nil)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve place", result.Error.Error())
	}

	if err := config.DB.Model(&place).Updates(updateData).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update place", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Place updated successfully", place)
}

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
func DeletePlace(c *fiber.Ctx) error {
	id := c.Params("id")

	placeID, err := uuid.Parse(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID format", nil)
	}

	result := config.DB.Delete(&models.Place{}, "id = ?", placeID)
	if result.Error != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete place", result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Place not found", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Place deleted successfully", nil)
}
