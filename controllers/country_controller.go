package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/services/implementations"
)

type CountryController struct {
	countryService *implementations.CountryService
}

func NewCountryController(countryService *implementations.CountryService) *CountryController {
	return &CountryController{countryService: countryService}
}

func (ctrl *CountryController) List(c fiber.Ctx) error {
	countries, err := ctrl.countryService.FindAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": countries})
}

func (ctrl *CountryController) ListStates(c fiber.Ctx) error {
	countryID := c.Params("id")
	if countryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Country ID is required"})
	}

	states, err := ctrl.countryService.FindStatesByCountryID(c.Context(), countryID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": states})
}

func (ctrl *CountryController) ListDistricts(c fiber.Ctx) error {
	stateID := c.Params("id")
	if stateID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid state ID"})
	}

	districts, err := ctrl.countryService.FindDistrictsByStateID(c.Context(), stateID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": districts})
}
