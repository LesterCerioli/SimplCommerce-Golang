package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/initializers"
)

type ReviewController struct {
	svc *initializers.Services
}

func NewReviewController(svc *initializers.Services) *ReviewController {
	return &ReviewController{svc: svc}
}

func (h *ReviewController) GetProductReviews(c fiber.Ctx) error {
	productIDStr := c.Query("productId")
	if productIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "productId is required"})
	}
	page := parseInt(c.Query("page", "1"), 1)
	pageSize := parseInt(c.Query("pageSize", "20"), 20)

	reviews, total, err := h.svc.ReviewService.GetProductReviews(c.Context(), productIDStr, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": reviews, "meta": fiber.Map{"page": page, "pageSize": pageSize, "totalItems": total}})
}

func (h *ReviewController) CreateReview(c fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req struct {
		Title        string `json:"title"`
		Comment      string `json:"comment"`
		Rating       int    `json:"rating"`
		ReviewerName string `json:"reviewerName"`
		EntityTypeID string `json:"entityTypeId"`
		EntityID     string `json:"entityId"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	review, err := h.svc.ReviewService.CreateReview(c.Context(), userID, req.Title, req.Comment, req.Rating, req.ReviewerName, req.EntityTypeID, req.EntityID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": review})
}

func (h *ReviewController) UpdateStatus(c fiber.Ctx) error {
	reviewID := c.Params("id")
	if reviewID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid review id"})
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}
	if req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "status is required"})
	}

	if err := h.svc.ReviewService.UpdateStatus(c.Context(), reviewID, req.Status); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": nil})
}

func (h *ReviewController) DeleteReview(c fiber.Ctx) error {
	reviewID := c.Params("id")
	if reviewID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid review id"})
	}

	if err := h.svc.ReviewService.DeleteReview(c.Context(), reviewID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": nil})
}
