package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/myplate/backend/internal/models"
	"github.com/myplate/backend/internal/repositories"
)

type AuthService struct {
	userRepo *repositories.UserRepository
	jwtSecret string
	telegramToken string
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		jwtSecret:     os.Getenv("JWT_SECRET"),
		telegramToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
	}
}

func (s *AuthService) ValidateTelegramInitData(initData string) (map[string]string, error) {
	// Parse initData
	params, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("invalid initData format: %w", err)
	}
	
	// Extract hash
	hash := params.Get("hash")
	if hash == "" {
		return nil, fmt.Errorf("hash not found in initData")
	}
	
	// Remove hash from params for validation
	params.Del("hash")
	
	// Create data-check-string
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params.Get(k)))
	}
	dataCheckString := strings.Join(parts, "\n")
	
	// Calculate secret key
	secretKey := hmac.New(sha256.New, []byte("WebAppData"))
	secretKey.Write([]byte(s.telegramToken))
	
	// Calculate hash
	calculatedHash := hmac.New(sha256.New, secretKey.Sum(nil))
	calculatedHash.Write([]byte(dataCheckString))
	calculatedHashHex := hex.EncodeToString(calculatedHash.Sum(nil))
	
	// Validate
	if calculatedHashHex != hash {
		return nil, fmt.Errorf("invalid hash")
	}
	
	// Parse user data
	result := make(map[string]string)
	for k, v := range params {
		result[k] = v[0]
	}
	
	return result, nil
}

func (s *AuthService) AuthenticateTelegram(initData string) (*models.User, string, error) {
	// Validate initData
	data, err := s.ValidateTelegramInitData(initData)
	if err != nil {
		return nil, "", err
	}
	
	// Extract user info (simplified - in production, parse JSON from 'user' field)
	telegramID := data["id"]
	if telegramID == "" {
		return nil, "", fmt.Errorf("user ID not found")
	}
	
	// Create or update user
	var telegramIDInt int64
	fmt.Sscanf(telegramID, "%d", &telegramIDInt)
	
	user, err := s.userRepo.CreateOrUpdate(
		telegramIDInt,
		data["username"],
		data["first_name"],
		data["last_name"],
	)
	if err != nil {
		return nil, "", err
	}
	
	// Generate JWT token
	token, err := s.GenerateJWT(user.ID)
	if err != nil {
		return nil, "", err
	}
	
	return user, token, nil
}

func (s *AuthService) GenerateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateJWT(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})
	
	if err != nil {
		return 0, err
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, fmt.Errorf("invalid user_id in token")
		}
		return int(userID), nil
	}
	
	return 0, fmt.Errorf("invalid token")
}

// CreateTestUser - создает тестового пользователя и возвращает токен (для разработки)
func (s *AuthService) CreateTestUser() (*models.User, string, error) {
	// Создаем или получаем тестового пользователя
	user, err := s.userRepo.CreateOrUpdate(
		123456789, // Тестовый telegram_id
		"testuser",
		"Test",
		"User",
	)
	if err != nil {
		return nil, "", err
	}
	
	// Генерируем токен
	token, err := s.GenerateJWT(user.ID)
	if err != nil {
		return nil, "", err
	}
	
	return user, token, nil
}

// Register создает нового пользователя с email и паролем
func (s *AuthService) Register(email, password, firstName, lastName string) (*models.User, string, error) {
	// Проверяем, существует ли пользователь с таким email
	existingUser, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка при проверке email: %w", err)
	}
	if existingUser != nil {
		return nil, "", fmt.Errorf("пользователь с таким email уже существует")
	}

	// Хешируем пароль
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка при хешировании пароля: %w", err)
	}

	// Создаем пользователя
	user, err := s.userRepo.Create(email, string(passwordHash), firstName, lastName)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка при создании пользователя: %w", err)
	}

	// Генерируем JWT токен
	token, err := s.GenerateJWT(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка при генерации токена: %w", err)
	}

	return user, token, nil
}

// Login авторизует пользователя по email и паролю
func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	// Получаем пользователя по email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка при поиске пользователя: %w", err)
	}
	if user == nil {
		return nil, "", fmt.Errorf("неверный email или пароль")
	}

	// Проверяем пароль
	if user.PasswordHash == "" {
		return nil, "", fmt.Errorf("у пользователя не установлен пароль")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", fmt.Errorf("неверный email или пароль")
	}

	// Генерируем JWT токен
	token, err := s.GenerateJWT(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка при генерации токена: %w", err)
	}

	// Очищаем password_hash перед возвратом
	user.PasswordHash = ""

	return user, token, nil
}

// GetUserProfile получает профиль пользователя по ID
func (s *AuthService) GetUserProfile(userID int) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении профиля: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("пользователь не найден")
	}
	return user, nil
}

// UpdatePassword обновляет пароль пользователя
func (s *AuthService) UpdatePassword(userID int, oldPassword, newPassword string) error {
	// Получаем пользователя
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("ошибка при получении пользователя: %w", err)
	}
	if user == nil {
		return fmt.Errorf("пользователь не найден")
	}

	// Если у пользователя есть пароль, проверяем старый
	if user.PasswordHash != "" {
		userWithPassword, err := s.userRepo.GetByEmail(user.Email)
		if err != nil {
			return fmt.Errorf("ошибка при проверке пароля: %w", err)
		}
		if userWithPassword == nil {
			return fmt.Errorf("пользователь не найден")
		}

		err = bcrypt.CompareHashAndPassword([]byte(userWithPassword.PasswordHash), []byte(oldPassword))
		if err != nil {
			return fmt.Errorf("неверный текущий пароль")
		}
	}

	// Хешируем новый пароль
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("ошибка при хешировании пароля: %w", err)
	}

	// Обновляем пароль
	err = s.userRepo.UpdatePassword(userID, string(newPasswordHash))
	if err != nil {
		return fmt.Errorf("ошибка при обновлении пароля: %w", err)
	}

	return nil
}

// UpdateProfile обновляет профиль пользователя
func (s *AuthService) UpdateProfile(userID int, firstName, lastName string) (*models.User, error) {
	err := s.userRepo.UpdateProfile(userID, firstName, lastName)
	if err != nil {
		return nil, fmt.Errorf("ошибка при обновлении профиля: %w", err)
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении обновленного профиля: %w", err)
	}

	return user, nil
}


