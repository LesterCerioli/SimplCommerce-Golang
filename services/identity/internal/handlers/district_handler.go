package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/identity/internal/services"
)

type DistrictHandler struct {
	districtService services.DistrictService
}

func NewDistrictHandler(districtService services.DistrictService) *DistrictHandler {
	return &DistrictHandler{districtService: districtService}
}

func (h *DistrictHandler) Index(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "invalid state id"))
	}

	districts, err := h.districtService.FindByStateOrProvinceID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(fiber.StatusInternalServerError, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(districts))
}
