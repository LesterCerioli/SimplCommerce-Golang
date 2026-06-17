package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/initializers"
)

type NotificationController struct {
	svc *initializers.Services
}

func NewNotificationController(svc *initializers.Services) *NotificationController {
	return &NotificationController{svc: svc}
}

func (h *NotificationController) ListNotifications(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	notifications, err := h.svc.NotificationService.GetUserNotifications(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": notifications})
}

func (h *NotificationController) MarkAsRead(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	notifID := parseUint(c.Params("id"))

	if err := h.svc.NotificationService.MarkAsRead(c.Context(), notifID, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": nil})
}

func (h *NotificationController) UnreadCount(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	count, err := h.svc.NotificationService.GetUnreadCount(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": fiber.Map{"count": count}})
}
