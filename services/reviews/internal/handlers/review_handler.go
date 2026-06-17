package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/reviews/internal/models"
	"github.com/simplcommerce-go/services/reviews/internal/services"
)

type ReviewHandler struct {
	reviewService services.ReviewService
}

func NewReviewHandler(reviewService services.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

type createReviewRequest struct {
	Title        string `json:"title"`
	Comment      string `json:"comment"`
	Rating       int    `json:"rating"`
	ReviewerName string `json:"reviewerName"`
	EntityTypeID string `json:"entityTypeId"`
	EntityID     uint   `json:"entityId"`
}

type statusUpdateRequest struct {
	Status string `json:"status"`
}

func (h *ReviewHandler) GetProductReviews(c fiber.Ctx) error {
	productID, err := strconv.ParseUint(c.Query("productId"), 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid productId")
	}

	reviews, err := h.reviewService.GetEntityReviews(c.Context(), "Product", uint(productID))
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(reviews))
}

func (h *ReviewHandler) Create(c fiber.Ctx) error {
	var req createReviewRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return response.Error(fiber.StatusUnauthorized, "unauthorized")
	}

	review := &models.Review{
		UserID:       userID,
		Title:        req.Title,
		Comment:      req.Comment,
		Rating:       req.Rating,
		ReviewerName: req.ReviewerName,
		Status:       "Pending",
		EntityTypeID: req.EntityTypeID,
		EntityID:     req.EntityID,
	}

	if err := h.reviewService.Create(c.Context(), review); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(review))
}

func (h *ReviewHandler) UpdateStatus(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid id")
	}

	var req statusUpdateRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	switch req.Status {
	case "Approved":
		err = h.reviewService.Approve(c.Context(), uint(id))
	case "Rejected":
		err = h.reviewService.Reject(c.Context(), uint(id))
	default:
		return response.Error(fiber.StatusBadRequest, "invalid status, must be Approved or Rejected")
	}

	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(nil))
}

func (h *ReviewHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid id")
	}

	if err := h.reviewService.Delete(c.Context(), uint(id)); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(nil))
}
