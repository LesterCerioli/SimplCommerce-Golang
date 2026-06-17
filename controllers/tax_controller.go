package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/initializers"
)

type TaxController struct {
	svc *initializers.Services
}

func NewTaxController(svc *initializers.Services) *TaxController {
	return &TaxController{svc: svc}
}

func (h *TaxController) ListTaxClasses(c fiber.Ctx) error {
	classes, err := h.svc.TaxService.GetTaxClasses(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": classes})
}

func (h *TaxController) CreateTaxClass(c fiber.Ctx) error {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "name is required"})
	}

	tc, err := h.svc.TaxService.CreateTaxClass(c.Context(), req.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": tc})
}

func (h *TaxController) UpdateTaxClass(c fiber.Ctx) error {
	tcID := parseUint(c.Params("id"))

	var req struct {
		Name string `json:"name"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	tc, err := h.svc.TaxService.UpdateTaxClass(c.Context(), tcID, req.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": tc})
}

func (h *TaxController) DeleteTaxClass(c fiber.Ctx) error {
	tcID := parseUint(c.Params("id"))
	if err := h.svc.TaxService.DeleteTaxClass(c.Context(), tcID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": nil})
}

func (h *TaxController) ListTaxRates(c fiber.Ctx) error {
	page := parseInt(c.Query("page", "1"), 1)
	pageSize := parseInt(c.Query("pageSize", "20"), 20)

	rates, total, err := h.svc.TaxService.GetTaxRates(c.Context(), page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": rates, "meta": fiber.Map{"page": page, "pageSize": pageSize, "totalItems": total}})
}

func (h *TaxController) CreateTaxRate(c fiber.Ctx) error {
	var req struct {
		TaxClassID        uint   `json:"taxClassId"`
		CountryID         string `json:"countryId"`
		StateOrProvinceID *uint  `json:"stateOrProvinceId"`
		Rate              float64 `json:"rate"`
		ZipCode           string `json:"zipCode"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	tr, err := h.svc.TaxService.CreateTaxRate(c.Context(), req.TaxClassID, req.CountryID, req.StateOrProvinceID, req.Rate, req.ZipCode)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": tr})
}

func (h *TaxController) UpdateTaxRate(c fiber.Ctx) error {
	trID := parseUint(c.Params("id"))

	var req struct {
		TaxClassID        *uint   `json:"taxClassId"`
		CountryID         *string `json:"countryId"`
		StateOrProvinceID *uint   `json:"stateOrProvinceId"`
		Rate              *float64 `json:"rate"`
		ZipCode           *string `json:"zipCode"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	tr, err := h.svc.TaxService.UpdateTaxRate(c.Context(), trID, req.TaxClassID, req.CountryID, req.StateOrProvinceID, req.Rate, req.ZipCode)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": tr})
}

func (h *TaxController) DeleteTaxRate(c fiber.Ctx) error {
	trID := parseUint(c.Params("id"))

	if err := h.svc.TaxService.DeleteTaxRate(c.Context(), trID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": nil})
}
