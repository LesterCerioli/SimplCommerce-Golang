package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/reviews/internal/models"
	"github.com/simplcommerce-go/services/reviews/internal/services"
)

type ReplyHandler struct {
	replyService services.ReplyService
}

func NewReplyHandler(replyService services.ReplyService) *ReplyHandler {
	return &ReplyHandler{replyService: replyService}
}

type createReplyRequest struct {
	Comment     string `json:"comment"`
	ReplierName string `json:"replierName"`
}

func (h *ReplyHandler) Create(c fiber.Ctx) error {
	reviewID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid review id")
	}

	var req createReplyRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return response.Error(fiber.StatusUnauthorized, "unauthorized")
	}

	reply := &models.Reply{
		ReviewID:    uint(reviewID),
		UserID:      userID,
		Comment:     req.Comment,
		ReplierName: req.ReplierName,
		Status:      "Pending",
	}

	if err := h.replyService.Create(c.Context(), reply); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(reply))
}
