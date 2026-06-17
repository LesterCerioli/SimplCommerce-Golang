package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/cms/internal/models"
	"github.com/simplcommerce-go/services/cms/internal/services"
)

type MenuHandler struct {
	menuService services.MenuService
}

func NewMenuHandler(menuService services.MenuService) *MenuHandler {
	return &MenuHandler{menuService: menuService}
}

type createMenuRequest struct {
	Name        string `json:"name"`
	IsPublished bool   `json:"isPublished"`
}

type createMenuItemRequest struct {
	Name         string `json:"name"`
	CustomLink   string `json:"customLink"`
	EntityID     *uint  `json:"entityId"`
	DisplayOrder int    `json:"displayOrder"`
	ParentID     *uint  `json:"parentId"`
}

func (h *MenuHandler) List(c fiber.Ctx) error {
	menus, err := h.menuService.FindAll(c.Context())
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(response.Success(menus))
}

func (h *MenuHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid id")
	}

	menu, err := h.menuService.FindByID(c.Context(), uint(id))
	if err != nil {
		return response.Error(fiber.StatusNotFound, "menu not found")
	}

	return c.JSON(response.Success(menu))
}

func (h *MenuHandler) Create(c fiber.Ctx) error {
	var req createMenuRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	menu := &models.Menu{
		Name:        req.Name,
		IsPublished: req.IsPublished,
	}

	if err := h.menuService.Create(c.Context(), menu); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(menu))
}

func (h *MenuHandler) Update(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid id")
	}

	var req createMenuRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	menu, err := h.menuService.FindByID(c.Context(), uint(id))
	if err != nil {
		return response.Error(fiber.StatusNotFound, "menu not found")
	}

	menu.Name = req.Name
	menu.IsPublished = req.IsPublished

	if err := h.menuService.Update(c.Context(), menu); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(menu))
}

func (h *MenuHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid id")
	}

	if err := h.menuService.Delete(c.Context(), uint(id)); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(nil))
}

func (h *MenuHandler) AddItem(c fiber.Ctx) error {
	menuID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid menu id")
	}

	var req createMenuItemRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	item := &models.MenuItem{
		MenuID:       uint(menuID),
		Name:         req.Name,
		CustomLink:   req.CustomLink,
		EntityID:     req.EntityID,
		DisplayOrder: req.DisplayOrder,
		ParentID:     req.ParentID,
	}

	if err := h.menuService.AddItem(c.Context(), item); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(item))
}

func (h *MenuHandler) UpdateItem(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid id")
	}

	var req createMenuItemRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	item, err := h.menuService.FindItemByID(c.Context(), uint(id))
	if err != nil {
		return response.Error(fiber.StatusNotFound, "menu item not found")
	}

	item.Name = req.Name
	item.CustomLink = req.CustomLink
	item.EntityID = req.EntityID
	item.DisplayOrder = req.DisplayOrder
	item.ParentID = req.ParentID

	if err := h.menuService.UpdateItem(c.Context(), item); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(item))
}

func (h *MenuHandler) DeleteItem(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid id")
	}

	if err := h.menuService.DeleteItem(c.Context(), uint(id)); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(nil))
}
