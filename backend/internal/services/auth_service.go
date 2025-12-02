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


