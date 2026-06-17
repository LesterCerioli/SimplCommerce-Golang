package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/identity/internal/services"
)

type CountryHandler struct {
	countryService services.CountryService
}

func NewCountryHandler(countryService services.CountryService) *CountryHandler {
	return &CountryHandler{countryService: countryService}
}

func (h *CountryHandler) Index(c fiber.Ctx) error {
	countries, err := h.countryService.FindAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(fiber.StatusInternalServerError, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(countries))
}

func (h *CountryHandler) Show(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "invalid country id"))
	}

	country, err := h.countryService.FindByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.Error(fiber.StatusNotFound, "country not found"))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(country))
}
