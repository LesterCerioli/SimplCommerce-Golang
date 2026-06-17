package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/initializers"
)

type MenuController struct {
	svc *initializers.Services
}

func NewMenuController(svc *initializers.Services) *MenuController {
	return &MenuController{svc: svc}
}

func (h *MenuController) ListMenus(c fiber.Ctx) error {
	menus, err := h.svc.MenuService.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": menus})
}

func (h *MenuController) GetMenu(c fiber.Ctx) error {
	menuID := c.Params("id")
	if menuID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid menu id"})
	}

	menu, err := h.svc.MenuService.GetByID(c.Context(), menuID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": menu})
}

func (h *MenuController) CreateMenu(c fiber.Ctx) error {
	var req struct {
		Name        string `json:"name"`
		IsPublished bool   `json:"isPublished"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "name is required"})
	}

	menu, err := h.svc.MenuService.Create(c.Context(), req.Name, req.IsPublished)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": menu})
}

func (h *MenuController) UpdateMenu(c fiber.Ctx) error {
	menuID := c.Params("id")
	if menuID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid menu id"})
	}

	var req struct {
		Name        string `json:"name"`
		IsPublished bool   `json:"isPublished"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	menu, err := h.svc.MenuService.Update(c.Context(), menuID, req.Name, req.IsPublished)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": menu})
}

func (h *MenuController) DeleteMenu(c fiber.Ctx) error {
	menuID := c.Params("id")
	if menuID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid menu id"})
	}

	if err := h.svc.MenuService.Delete(c.Context(), menuID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": nil})
}

func (h *MenuController) AddMenuItem(c fiber.Ctx) error {
	menuID := c.Params("id")
	if menuID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid menu id"})
	}

	var req struct {
		ParentID     *string `json:"parentId"`
		EntityID     *string `json:"entityId"`
		CustomLink   string  `json:"customLink"`
		Name         string  `json:"name"`
		DisplayOrder int     `json:"displayOrder"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	item, err := h.svc.MenuService.AddItem(c.Context(), menuID, req.ParentID, req.EntityID, req.CustomLink, req.Name, req.DisplayOrder)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": item})
}

func (h *MenuController) UpdateMenuItem(c fiber.Ctx) error {
	itemID := c.Params("id")
	if itemID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid menu item id"})
	}

	var req struct {
		ParentID     *string `json:"parentId"`
		EntityID     *string `json:"entityId"`
		CustomLink   string  `json:"customLink"`
		Name         string  `json:"name"`
		DisplayOrder int     `json:"displayOrder"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	item, err := h.svc.MenuService.UpdateItem(c.Context(), itemID, req.ParentID, req.EntityID, req.CustomLink, req.Name, req.DisplayOrder)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": item})
}

func (h *MenuController) DeleteMenuItem(c fiber.Ctx) error {
	itemID := c.Params("id")
	if itemID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid menu item id"})
	}

	if err := h.svc.MenuService.DeleteItem(c.Context(), itemID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": nil})
}
