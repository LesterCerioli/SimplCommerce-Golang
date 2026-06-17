package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/initializers"
)

type ActivityController struct {
	svc *initializers.Services
}

func NewActivityController(svc *initializers.Services) *ActivityController {
	return &ActivityController{svc: svc}
}

func (h *ActivityController) ListActivities(c fiber.Ctx) error {
	page := parseInt(c.Query("page", "1"), 1)
	pageSize := parseInt(c.Query("pageSize", "20"), 20)

	activities, total, err := h.svc.ActivityService.GetActivities(c.Context(), page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": activities, "meta": fiber.Map{"page": page, "pageSize": pageSize, "totalItems": total}})
}
