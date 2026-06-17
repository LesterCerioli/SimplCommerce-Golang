package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/initializers"
)

type SearchController struct {
	svc *initializers.Services
}

func NewSearchController(svc *initializers.Services) *SearchController {
	return &SearchController{svc: svc}
}

func (h *SearchController) LogSearch(c fiber.Ctx) error {
	var req struct {
		QueryText    string `json:"queryText"`
		ResultsCount int    `json:"resultsCount"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}
	if req.QueryText == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "queryText is required"})
	}

	if err := h.svc.SearchService.LogSearch(c.Context(), req.QueryText, req.ResultsCount); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": nil})
}

func (h *SearchController) PopularSearches(c fiber.Ctx) error {
	limit := parseInt(c.Query("limit", "10"), 10)

	terms, err := h.svc.SearchService.GetPopularSearches(c.Context(), limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": terms})
}
