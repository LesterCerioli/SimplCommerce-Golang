package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/pricing/internal/models"
	"github.com/simplcommerce-go/services/pricing/internal/repositories"
)

type CatalogRuleHandler struct {
	catalogRuleRepo repositories.CatalogRuleRepository
}

func NewCatalogRuleHandler(catalogRuleRepo repositories.CatalogRuleRepository) *CatalogRuleHandler {
	return &CatalogRuleHandler{catalogRuleRepo: catalogRuleRepo}
}

type createCatalogRuleRequest struct {
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	IsActive          *bool      `json:"isActive"`
	StartOn           *time.Time `json:"startOn"`
	EndOn             *time.Time `json:"endOn"`
	RuleToApply       string     `json:"ruleToApply"`
	DiscountAmount    float64    `json:"discountAmount"`
	MaxDiscountAmount *float64   `json:"maxDiscountAmount"`
	CustomerGroupIDs  []uint     `json:"customerGroupIds"`
}

type updateCatalogRuleRequest struct {
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	IsActive          *bool      `json:"isActive"`
	StartOn           *time.Time `json:"startOn"`
	EndOn             *time.Time `json:"endOn"`
	RuleToApply       string     `json:"ruleToApply"`
	DiscountAmount    float64    `json:"discountAmount"`
	MaxDiscountAmount *float64   `json:"maxDiscountAmount"`
	CustomerGroupIDs  []uint     `json:"customerGroupIds"`
}

func (h *CatalogRuleHandler) List(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))

	rules, total, err := h.catalogRuleRepo.Paginate(c.Context(), page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(rules, page, pageSize, total))
}

func (h *CatalogRuleHandler) Create(c fiber.Ctx) error {
	var req createCatalogRuleRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	ruleToApply := req.RuleToApply
	if ruleToApply == "" {
		ruleToApply = "by_fixed"
	}

	rule := &models.CatalogRule{
		Name:              req.Name,
		Description:       req.Description,
		IsActive:          isActive,
		StartOn:           req.StartOn,
		EndOn:             req.EndOn,
		RuleToApply:       ruleToApply,
		DiscountAmount:    req.DiscountAmount,
		MaxDiscountAmount: req.MaxDiscountAmount,
	}

	for _, cgID := range req.CustomerGroupIDs {
		rule.CustomerGroups = append(rule.CustomerGroups, models.CatalogRuleCustomerGroup{
			CatalogRuleID:   0,
			CustomerGroupID: cgID,
		})
	}

	if err := h.catalogRuleRepo.Create(c.Context(), rule); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(rule))
}

func (h *CatalogRuleHandler) Update(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid catalog rule id")
	}

	var req updateCatalogRuleRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	rule, err := h.catalogRuleRepo.FindByID(c.Context(), uint(id))
	if err != nil {
		return response.Error(fiber.StatusNotFound, "catalog rule not found")
	}

	if req.Name != "" {
		rule.Name = req.Name
	}
	rule.Description = req.Description
	if req.IsActive != nil {
		rule.IsActive = *req.IsActive
	}
	rule.StartOn = req.StartOn
	rule.EndOn = req.EndOn
	if req.RuleToApply != "" {
		rule.RuleToApply = req.RuleToApply
	}
	rule.DiscountAmount = req.DiscountAmount
	rule.MaxDiscountAmount = req.MaxDiscountAmount

	rule.CustomerGroups = nil
	for _, cgID := range req.CustomerGroupIDs {
		rule.CustomerGroups = append(rule.CustomerGroups, models.CatalogRuleCustomerGroup{
			CatalogRuleID:   rule.ID,
			CustomerGroupID: cgID,
		})
	}

	if err := h.catalogRuleRepo.Update(c.Context(), rule); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(rule))
}

func (h *CatalogRuleHandler) Delete(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid catalog rule id")
	}

	if err := h.catalogRuleRepo.Delete(c.Context(), uint(id)); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(nil))
}
