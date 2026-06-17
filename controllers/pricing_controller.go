package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/initializers"
	"simplcommerce/services/implementations"
)

type PricingController struct {
	svc *initializers.Services
}

func NewPricingController(svc *initializers.Services) *PricingController {
	return &PricingController{svc: svc}
}

func (h *PricingController) ListCartRules(c fiber.Ctx) error {
	rules, err := h.svc.PricingService.GetCartRules(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": rules})
}

func (h *PricingController) GetCartRule(c fiber.Ctx) error {
	ruleID := parseUint(c.Params("id"))

	rule, err := h.svc.PricingService.GetCartRule(c.Context(), ruleID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": rule})
}

func (h *PricingController) CreateCartRule(c fiber.Ctx) error {
	var req implementations.CreateCartRuleRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	rule, err := h.svc.PricingService.CreateCartRule(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": rule})
}

func (h *PricingController) UpdateCartRule(c fiber.Ctx) error {
	ruleID := parseUint(c.Params("id"))

	var req implementations.CreateCartRuleRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	rule, err := h.svc.PricingService.UpdateCartRule(c.Context(), ruleID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": rule})
}

func (h *PricingController) DeleteCartRule(c fiber.Ctx) error {
	ruleID := parseUint(c.Params("id"))

	if err := h.svc.PricingService.DeleteCartRule(c.Context(), ruleID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": nil})
}

func (h *PricingController) ValidateCoupon(c fiber.Ctx) error {
	var req struct {
		Code string `json:"code"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}
	if req.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "code is required"})
	}

	result, err := h.svc.PricingService.ValidateCoupon(c.Context(), req.Code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": result})
}

func (h *PricingController) ListCatalogRules(c fiber.Ctx) error {
	rules, err := h.svc.PricingService.GetCatalogRules(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": rules})
}

func (h *PricingController) CreateCatalogRule(c fiber.Ctx) error {
	var req implementations.CreateCatalogRuleRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	rule, err := h.svc.PricingService.CreateCatalogRule(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": rule})
}

func (h *PricingController) UpdateCatalogRule(c fiber.Ctx) error {
	ruleID := parseUint(c.Params("id"))

	var req implementations.CreateCatalogRuleRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	rule, err := h.svc.PricingService.UpdateCatalogRule(c.Context(), ruleID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": rule})
}

func (h *PricingController) DeleteCatalogRule(c fiber.Ctx) error {
	ruleID := parseUint(c.Params("id"))

	if err := h.svc.PricingService.DeleteCatalogRule(c.Context(), ruleID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": nil})
}
