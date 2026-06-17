package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/notifications/internal/services"
)

type NotificationHandler struct {
	notificationService services.NotificationService
}

func NewNotificationHandler(notificationService services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

func (h *NotificationHandler) GetUserNotifications(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return response.Error(fiber.StatusUnauthorized, "unauthorized")
	}

	notifications, err := h.notificationService.GetUserNotifications(c.Context(), userID)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(notifications))
}

func (h *NotificationHandler) MarkAsRead(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return response.Error(fiber.StatusUnauthorized, "unauthorized")
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid id")
	}

	notification, err := h.notificationService.FindByID(c.Context(), uint(id))
	if err != nil {
		return response.Error(fiber.StatusNotFound, "notification not found")
	}

	if notification.UserID != userID {
		return response.Error(fiber.StatusForbidden, "forbidden")
	}

	if err := h.notificationService.MarkAsRead(c.Context(), uint(id)); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(nil))
}

func (h *NotificationHandler) UnreadCount(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return response.Error(fiber.StatusUnauthorized, "unauthorized")
	}

	count, err := h.notificationService.UnreadCount(c.Context(), userID)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(map[string]int64{"unreadCount": count}))
}
