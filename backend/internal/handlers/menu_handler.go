package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/myplate/backend/internal/models"
	"github.com/myplate/backend/internal/services"
)

type MenuHandler struct {
	menuService *services.MenuService
}

func NewMenuHandler(menuService *services.MenuService) *MenuHandler {
	return &MenuHandler{
		menuService: menuService,
	}
}

func (h *MenuHandler) Generate(c *fiber.Ctx) error {
	var req models.MenuGenerateRequest
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверное тело запроса"})
	}
	
	menu, err := h.menuService.GenerateMenu(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	return c.JSON(menu)
}

func (h *MenuHandler) GetDaily(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	dateStr := c.Query("date")
	
	var date time.Time
	var err error
	if dateStr != "" {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Неверный формат даты"})
		}
	} else {
		date = time.Now()
	}
	
	menu, err := h.menuService.GetDaily(userID, date)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	return c.JSON(menu)
}

func (h *MenuHandler) GetAll(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	
	menus, err := h.menuService.GetAllByUserID(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	return c.JSON(menus)
}

func (h *MenuHandler) GetByID(c *fiber.Ctx) error {
	menuID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID меню"})
	}
	
	menu, err := h.menuService.GetByID(menuID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	if menu == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Меню не найдено"})
	}
	
	return c.JSON(menu)
}

