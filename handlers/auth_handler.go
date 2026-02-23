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

// CreateAuth godoc
// @Summary Create Auth
// @Description Create a new auth record
// @Tags Auth
// @Accept json
// @Produce json
// @Param auth body models.Auth true "Auth payload"
// @Success 201 {object} models.Auth
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/auth [post]
func CreateAuth(c *fiber.Ctx) error {
	var auth models.Auth

	if err := c.BodyParser(&auth); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate input
	if validationErrs := utils.ValidateStruct(auth); len(validationErrs) > 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", validationErrs)
	}

	if auth.UID == uuid.Nil {
		auth.UID = uuid.New()
	}

	result := config.DB.Create(&auth)
	if result.Error != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create auth", result.Error.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Auth created successfully", auth)
}

// GetAuth godoc
// @Summary Get auth by ID
// @Description Retrieve an auth record by UUID
// @Tags Auth
// @Accept json
// @Produce json
// @Param id path string true "Auth UUID"
// @Success 200 {object} models.Auth
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/auth/{id} [get]
func GetAuth(c *fiber.Ctx) error {
	id := c.Params("id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID format", nil)
	}

	var auth models.Auth

	result := config.DB.Preload("Profile").First(&auth, "uid = ?", uuid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Auth not found", nil)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve auth", result.Error.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Auth retrieved successfully", auth)
}

// ListAuth godoc
// @Summary List auths
// @Description List auth records with pagination
// @Tags Auth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {array} models.Auth
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth [get]
func ListAuth(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var auths []models.Auth
	var total int64

	config.DB.Model(&models.Auth{}).Count(&total)
	result := config.DB.
		Preload("Profile").
		Offset(offset).
		Limit(pageSize).
		Find(&auths)

	if result.Error != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve auths", result.Error.Error())
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	paginatedData := utils.PaginatedData{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		Items:      auths,
	}

	return utils.PaginatedResponse(c, fiber.StatusOK, "Auths retrieved successfully", paginatedData)
}

// UpdateAuth godoc
// @Summary Update an auth
// @Description Update an auth record by UUID
// @Tags Auth
// @Accept json
// @Produce json
// @Param id path string true "Auth UUID"
// @Param auth body models.Auth true "Auth payload"
// @Success 200 {object} models.Auth
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/auth/{id} [put]
func UpdateAuth(c *fiber.Ctx) error {
	id := c.Params("id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID format", nil)
	}

	var updateData models.Auth

	if err := c.BodyParser(&updateData); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	var auth models.Auth

	result := config.DB.First(&auth, "uid = ?", uuid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Auth not found", nil)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve auth", result.Error.Error())
	}

	if err := config.DB.Model(&auth).Updates(updateData).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update auth", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Auth updated successfully", auth)
}

// DeleteAuth godoc
// @Summary Delete an auth
// @Description Delete an auth record by UUID
// @Tags Auth
// @Accept json
// @Produce json
// @Param id path string true "Auth UUID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/auth/{id} [delete]
func DeleteAuth(c *fiber.Ctx) error {
	id := c.Params("id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID format", nil)
	}

	result := config.DB.Delete(&models.Auth{}, "uid = ?", uuid)
	if result.Error != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete auth", result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Auth not found", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Auth deleted successfully", nil)
}
