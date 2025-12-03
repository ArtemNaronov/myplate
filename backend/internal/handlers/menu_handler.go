package handlers

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/myplate/backend/internal/models"
	"github.com/myplate/backend/internal/services"
	"github.com/myplate/backend/pkg/database"
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

// GenerateWeekly генерирует меню на неделю
// GET /menu/weekly?adults=2&children=1&diet_type=vegetarian&allergies=nuts,dairy
func (h *MenuHandler) GenerateWeekly(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	
	var req models.WeeklyMenuRequest
	req.UserID = userID
	
	// Парсим обязательные параметры
	adultsStr := c.Query("adults")
	if adultsStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Параметр 'adults' обязателен"})
	}
	adults, err := strconv.Atoi(adultsStr)
	if err != nil || adults < 1 {
		return c.Status(400).JSON(fiber.Map{"error": "Параметр 'adults' должен быть положительным числом"})
	}
	req.Adults = adults
	
	childrenStr := c.Query("children")
	if childrenStr == "" {
		req.Children = 0
	} else {
		children, err := strconv.Atoi(childrenStr)
		if err != nil || children < 0 {
			return c.Status(400).JSON(fiber.Map{"error": "Параметр 'children' должен быть неотрицательным числом"})
		}
		req.Children = children
	}
	
	// Опциональные параметры
	req.DietType = c.Query("diet_type")
	maxTotalTimeStr := c.Query("max_total_time")
	if maxTotalTimeStr != "" {
		req.MaxTotalTime, _ = strconv.Atoi(maxTotalTimeStr)
		if req.MaxTotalTime == 0 {
			req.MaxTotalTime = 0 // Не учитываем, если 0
		}
	}
	maxTimePerMealStr := c.Query("max_time_per_meal")
	if maxTimePerMealStr != "" {
		req.MaxTimePerMeal, _ = strconv.Atoi(maxTimePerMealStr)
		if req.MaxTimePerMeal == 0 {
			req.MaxTimePerMeal = 0 // Не учитываем, если 0
		}
	}
	req.ConsiderPantry = c.Query("consider_pantry") == "true"
	req.PantryImportance = c.Query("pantry_importance")
	if req.PantryImportance == "" {
		req.PantryImportance = "prefer"
	}
	
	// Парсим allergies из query (через запятую)
	if allergiesStr := c.Query("allergies"); allergiesStr != "" {
		req.Allergies = strings.Split(allergiesStr, ",")
		// Убираем пробелы
		for i := range req.Allergies {
			req.Allergies[i] = strings.TrimSpace(req.Allergies[i])
		}
	}
	
	weeklyMenu, err := h.menuService.GenerateWeeklyMenu(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	return c.JSON(weeklyMenu)
}

// SaveWeeklyMenu сохраняет недельное меню в базу данных
// POST /menu/weekly/save
func (h *MenuHandler) SaveWeeklyMenu(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	
	var weeklyMenu models.WeeklyMenu
	if err := c.BodyParser(&weeklyMenu); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверное тело запроса: " + err.Error()})
	}
	
	menu, err := h.menuService.SaveWeeklyMenu(userID, &weeklyMenu)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	return c.Status(201).JSON(menu)
}

// GetWeeklyMenus получает все сохраненные недельные меню пользователя
// GET /menus/weekly
func (h *MenuHandler) GetWeeklyMenus(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	
	menus, err := h.menuService.GetWeeklyMenus(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	// Преобразуем Menu в формат с meals как JSON
	type WeeklyMenuResponse struct {
		ID                 int                    `json:"id"`
		UserID             int                    `json:"user_id"`
		Date               string                 `json:"date"`
		TotalCalories      int                    `json:"total_calories"`
		TotalTime          int                    `json:"total_time"`
		MenuType           string                 `json:"menu_type"`
		Meals              interface{}            `json:"meals"` // JSON данные недели
		IngredientsUsed    models.Ingredients     `json:"ingredients_used,omitempty"`
		MissingIngredients models.Ingredients     `json:"missing_ingredients,omitempty"`
		CreatedAt          string                 `json:"created_at"`
		UpdatedAt          string                 `json:"updated_at"`
	}
	
	// Получаем сырые данные meals из базы
	responses := make([]WeeklyMenuResponse, 0, len(menus))
	for _, menu := range menus {
		// Получаем meals напрямую из базы как JSON
		query := `SELECT meals FROM menus WHERE id = $1`
		var mealsJSON []byte
		err := database.DB.QueryRow(query, menu.ID).Scan(&mealsJSON)
		if err != nil {
			continue
		}
		
		var mealsData interface{}
		json.Unmarshal(mealsJSON, &mealsData)
		
		responses = append(responses, WeeklyMenuResponse{
			ID:                 menu.ID,
			UserID:             menu.UserID,
			Date:               menu.Date.Format(time.RFC3339),
			TotalCalories:      menu.TotalCalories,
			TotalTime:          menu.TotalTime,
			MenuType:           menu.MenuType,
			Meals:              mealsData,
			IngredientsUsed:    menu.IngredientsUsed,
			MissingIngredients: menu.MissingIngredients,
			CreatedAt:          menu.CreatedAt.Format(time.RFC3339),
			UpdatedAt:          menu.UpdatedAt.Format(time.RFC3339),
		})
	}
	
	return c.JSON(responses)
}

// Delete удаляет меню по ID
// DELETE /menus/:id
func (h *MenuHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	menuID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID меню"})
	}
	
	err = h.menuService.DeleteMenu(menuID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Меню не найдено"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	return c.Status(200).JSON(fiber.Map{"message": "Меню успешно удалено"})
}

