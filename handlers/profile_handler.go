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

// CreateProfile godoc
// @Summary Create Profile
// @Description Create a new profile record
// @Tags Profile
// @Accept json
// @Produce json
// @Param profile body models.Profile true "Profile payload"
// @Success 201 {object} utils.Response{data=models.Profile}
// @Failure 400 {object} utils.Response
// @Router /api/v1/profile [post]
func CreateProfile(c *fiber.Ctx) error {
	var profile models.Profile

	if err := c.BodyParser(&profile); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if validationErrs := utils.ValidateStruct(profile); len(validationErrs) > 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", validationErrs)
	}

	if profile.UID == uuid.Nil {
		profile.UID = uuid.New()
	}

	result := config.DB.Create(&profile)
	if result.Error != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create profile", result.Error.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Profile created successfully", profile)
}

// GetProfile godoc
// @Summary Get profile by UID
// @Description Retrieve a profile record by UUID
// @Tags Profile
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=models.Profile}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/profile [get]
func GetProfile(c *fiber.Ctx) error {
	uid, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID format", nil)
	}

	var profile models.Profile

	result := config.DB.First(&profile, "uid = ?", uid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Profile not found", nil)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve profile", result.Error.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Profile retrieved successfully", profile)
}

// ListProfile godoc
// @Summary List profiles
// @Description List profiles with pagination
// @Tags Profile
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} utils.Response{data=[]models.Profile}
// @Failure 500 {object} utils.Response
// @Router /api/v1/profiles [get]
func ListProfile(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var profiles []models.Profile
	var total int64

	config.DB.Model(&models.Profile{}).Count(&total)

	result := config.DB.
		Offset(offset).
		Limit(pageSize).
		Find(&profiles)

	if result.Error != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve profiles", result.Error.Error())
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	paginatedData := utils.PaginatedData{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		Items:      profiles,
	}

	return utils.PaginatedResponse(c, fiber.StatusOK, "Profiles retrieved successfully", paginatedData)
}

// UpdateProfile godoc
// @Summary Update profile
// @Description Update a profile by UUID
// @Tags Profile
// @Accept json
// @Produce json
// @Param profile body models.Profile true "Profile payload"
// @Success 200 {object} utils.Response{data=models.Profile}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/profile [put]
func UpdateProfile(c *fiber.Ctx) error {
	uidStr, ok := c.Locals("user_id").(string)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID format", nil)
	}

	var req models.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	var profile models.Profile
	if err := config.DB.First(&profile, "uid = ?", uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Profile not found", nil)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error", err.Error())
	}

	updates := make(map[string]interface{})
	updates["display_name"] = req.DisplayName
	updates["phone"] = req.Phone

	if len(updates) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No fields to update", nil)
	}

	if err := config.DB.Model(&profile).Updates(updates).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update profile", err.Error())
	}

	if err := config.DB.First(&profile, "uid = ?", uid).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to reload profile", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Profile updated successfully", profile)
}

// DeleteProfile godoc
// @Summary Delete profile
// @Description Delete a profile by UUID
// @Tags Profile
// @Accept json
// @Produce json
// @Param id path string true "Profile UUID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/profile [delete]
func DeleteProfile(c *fiber.Ctx) error {
	uidStr, ok := c.Locals("user_id").(string)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid UUID format", nil)
	}

	result := config.DB.Delete(&models.Profile{}, "uid = ?", uid)
	if result.Error != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete profile", result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Profile not found", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Profile deleted successfully", nil)
}
