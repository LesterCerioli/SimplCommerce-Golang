package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/identity/internal/models"
	"github.com/simplcommerce-go/services/identity/internal/services"
)

type VendorHandler struct {
	vendorService services.VendorService
}

func NewVendorHandler(vendorService services.VendorService) *VendorHandler {
	return &VendorHandler{vendorService: vendorService}
}

func (h *VendorHandler) Index(c fiber.Ctx) error {
	vendors, err := h.vendorService.FindAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(fiber.StatusInternalServerError, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(vendors))
}

func (h *VendorHandler) Create(c fiber.Ctx) error {
	var vendor models.Vendor
	if err := c.Bind().JSON(&vendor); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "invalid request body"))
	}

	result, err := h.vendorService.Create(c.Context(), &vendor)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(fiber.StatusInternalServerError, err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(result))
}

func (h *VendorHandler) Update(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "invalid vendor id"))
	}

	var vendor models.Vendor
	if err := c.Bind().JSON(&vendor); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "invalid request body"))
	}

	result, err := h.vendorService.Update(c.Context(), uint(id), &vendor)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(fiber.StatusInternalServerError, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(result))
}
