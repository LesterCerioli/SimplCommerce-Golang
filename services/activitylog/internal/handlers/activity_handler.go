package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/activitylog/internal/services"
)

type ActivityHandler struct {
	activityService services.ActivityService
}

func NewActivityHandler(activityService services.ActivityService) *ActivityHandler {
	return &ActivityHandler{activityService: activityService}
}

type logActivityRequest struct {
	ActivityType string `json:"activityType"`
	EntityID     uint   `json:"entityId"`
	EntityTypeID string `json:"entityTypeId"`
}

func (h *ActivityHandler) List(c fiber.Ctx) error {
	activities, err := h.activityService.GetAllActivities(c.Context())
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(response.Success(activities))
}

func (h *ActivityHandler) MostViewed(c fiber.Ctx) error {
	activities, err := h.activityService.GetMostViewed(c.Context(), 10)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(response.Success(activities))
}

func (h *ActivityHandler) LogActivity(c fiber.Ctx) error {
	var req logActivityRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return response.Error(fiber.StatusUnauthorized, "unauthorized")
	}

	if err := h.activityService.LogActivity(c.Context(), req.ActivityType, userID, req.EntityID, req.EntityTypeID); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(nil))
}
