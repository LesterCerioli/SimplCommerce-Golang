package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"simplcommerce/services/implementations"
)

type VendorController struct {
	vendorService *implementations.VendorService
}

func NewVendorController(vendorService *implementations.VendorService) *VendorController {
	return &VendorController{vendorService: vendorService}
}

func (ctrl *VendorController) List(c fiber.Ctx) error {
	vendors, err := ctrl.vendorService.FindAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": vendors})
}

func (ctrl *VendorController) Create(c fiber.Ctx) error {
	var req implementations.CreateVendorRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Vendor name is required"})
	}

	vendor, err := ctrl.vendorService.Create(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(vendor)
}

func (ctrl *VendorController) Update(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid vendor ID"})
	}

	var req implementations.UpdateVendorRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	vendor, err := ctrl.vendorService.Update(c.Context(), uint(id), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(vendor)
}
