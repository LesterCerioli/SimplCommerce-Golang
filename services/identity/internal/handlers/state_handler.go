package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/identity/internal/services"
)

type StateHandler struct {
	stateService services.StateOrProvinceService
}

func NewStateHandler(stateService services.StateOrProvinceService) *StateHandler {
	return &StateHandler{stateService: stateService}
}

func (h *StateHandler) Index(c fiber.Ctx) error {
	countryID := c.Params("id")
	if countryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "invalid country id"))
	}

	states, err := h.stateService.FindByCountryID(c.Context(), countryID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(fiber.StatusInternalServerError, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(states))
}

func (h *StateHandler) Show(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "invalid state id"))
	}

	state, err := h.stateService.FindByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.Error(fiber.StatusNotFound, "state not found"))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(state))
}
