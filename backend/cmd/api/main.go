package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/myplate/backend/internal/handlers"
	"github.com/myplate/backend/internal/middleware"
	"github.com/myplate/backend/internal/repositories"
	"github.com/myplate/backend/internal/services"
	"github.com/myplate/backend/pkg/database"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	
	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()
	
	// Initialize repositories
	userRepo := repositories.NewUserRepository()
	recipeRepo := repositories.NewRecipeRepository()
	goalsRepo := repositories.NewGoalsRepository()
	pantryRepo := repositories.NewPantryRepository()
	menuRepo := repositories.NewMenuRepository()
	shoppingRepo := repositories.NewShoppingListRepository()
	
	// Initialize services
	authService := services.NewAuthService(userRepo)
	recipeService := services.NewRecipeService(recipeRepo)
	userService := services.NewUserService(goalsRepo)
	pantryService := services.NewPantryService(pantryRepo)
	menuService := services.NewMenuService(recipeRepo, menuRepo, pantryRepo, shoppingRepo, goalsRepo)
	shoppingService := services.NewShoppingListService(shoppingRepo, menuRepo, recipeRepo, pantryRepo)
	adminRecipeService := services.NewAdminRecipeService(recipeRepo)
	
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	recipeHandler := handlers.NewRecipeHandler(recipeService)
	userHandler := handlers.NewUserHandler(userService)
	pantryHandler := handlers.NewPantryHandler(pantryService)
	menuHandler := handlers.NewMenuHandler(menuService)
	shoppingHandler := handlers.NewShoppingListHandler(shoppingService)
	adminRecipeHandler := handlers.NewAdminRecipeHandler(adminRecipeService)
	
	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})
	
	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
	
	// Public routes (должны быть ДО создания защищенной группы)
	app.Get("/auth/test", authHandler.TestAuth) // Тестовая авторизация для разработки
	app.Post("/auth/telegram", authHandler.AuthenticateTelegram)
	app.Post("/auth/register", authHandler.Register) // Регистрация
	app.Post("/auth/login", authHandler.Login)         // Вход
	app.Get("/recipes", recipeHandler.GetAll)
	app.Get("/recipes/:id", recipeHandler.GetByID)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	
	// Protected routes
	api := app.Group("/", middleware.AuthMiddleware(authService))
	
	// Auth routes (требуют авторизации)
	api.Get("/auth/profile", authHandler.GetProfile)              // Получить профиль
	api.Put("/auth/profile", authHandler.UpdateProfile)           // Обновить профиль
	api.Put("/auth/password", authHandler.UpdatePassword)         // Обновить пароль
	
	// Menu routes
	api.Post("/menus/generate", menuHandler.Generate)
	api.Get("/menu/weekly", menuHandler.GenerateWeekly) // Генерация недельного меню
	api.Post("/menu/weekly/save", menuHandler.SaveWeeklyMenu) // Сохранение недельного меню
	api.Get("/menus/weekly", menuHandler.GetWeeklyMenus) // Получение всех недельных меню
	api.Get("/menus/daily", menuHandler.GetDaily)
	api.Get("/menus", menuHandler.GetAll)
	api.Get("/menus/:id", menuHandler.GetByID)
	api.Delete("/menus/:id", menuHandler.Delete) // Удаление меню
	
	// User goals routes
	api.Post("/users/goals", userHandler.SetGoals)
	api.Get("/users/goals", userHandler.GetGoals)
	
	// Pantry routes
	api.Get("/pantry", pantryHandler.GetAll)
	api.Post("/pantry", pantryHandler.Create)
	api.Delete("/pantry/:id", pantryHandler.Delete)
	
	// Shopping list routes
	api.Get("/shopping-list/:menu_id", shoppingHandler.GetByMenuID)
	
	// Admin routes (требуют роль admin)
	admin := api.Group("/admin", middleware.AdminMiddleware())
	admin.Post("/recipes", adminRecipeHandler.Create)
	admin.Post("/recipes/import", adminRecipeHandler.Import)
	admin.Get("/recipes/export", adminRecipeHandler.Export)
	
	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}


