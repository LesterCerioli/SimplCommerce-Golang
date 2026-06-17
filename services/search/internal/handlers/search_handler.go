package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/search/internal/services"
)

type SearchHandler struct {
	searchService services.SearchService
}

func NewSearchHandler(searchService services.SearchService) *SearchHandler {
	return &SearchHandler{searchService: searchService}
}

type searchRequest struct {
	Query string `json:"query"`
}

type searchResult struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Price float64 `json:"price"`
}

func (h *SearchHandler) Search(c fiber.Ctx) error {
	var req searchRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	if err := h.searchService.LogSearch(c.Context(), req.Query, 0); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	mockResults := []searchResult{}
	return c.JSON(response.Success(mockResults))
}

func (h *SearchHandler) Popular(c fiber.Ctx) error {
	queries, err := h.searchService.GetPopularSearches(c.Context(), 10)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(response.Success(queries))
}
