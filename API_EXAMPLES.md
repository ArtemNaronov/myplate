# Примеры использования API

Этот документ содержит практические примеры использования API MyPlate.

## Базовые примеры

### 1. Получение токена (тестовая авторизация)

```bash
curl -X GET http://localhost:8080/auth/test
```

**Ответ:**
```json
{
  "user": {
    "id": 1,
    "telegram_id": 123456789,
    "username": "testuser",
    "first_name": "Test",
    "last_name": "User"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

Сохраните токен для последующих запросов:
```bash
export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 2. Получение списка рецептов

```bash
curl -X GET http://localhost:8080/recipes
```

### 3. Получение деталей рецепта

```bash
curl -X GET http://localhost:8080/recipes/1
```

## Работа с меню

### Генерация меню

```bash
curl -X POST http://localhost:8080/menus/generate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "target_calories": 2000,
    "diet_type": "vegetarian",
    "allergies": ["nuts"],
    "max_total_time": 60,
    "consider_pantry": true,
    "pantry_importance": "prefer"
  }'
```

**Минимальный запрос:**
```bash
curl -X POST http://localhost:8080/menus/generate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "target_calories": 2000,
    "consider_pantry": false
  }'
```

### Получение списка всех меню

```bash
curl -X GET http://localhost:8080/menus \
  -H "Authorization: Bearer $TOKEN"
```

### Получение меню по ID

```bash
curl -X GET http://localhost:8080/menus/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Получение меню на конкретную дату

```bash
curl -X GET "http://localhost:8080/menus/daily?date=2025-12-02" \
  -H "Authorization: Bearer $TOKEN"
```

## Работа с кладовой

### Получение списка продуктов в кладовой

```bash
curl -X GET http://localhost:8080/pantry \
  -H "Authorization: Bearer $TOKEN"
```

### Добавление продукта в кладовую

```bash
curl -X POST http://localhost:8080/pantry \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Молоко",
    "quantity": 500,
    "unit": "мл"
  }'
```

### Удаление продукта из кладовой

```bash
curl -X DELETE http://localhost:8080/pantry/1 \
  -H "Authorization: Bearer $TOKEN"
```

## Работа с целями

### Установка целей пользователя

```bash
curl -X POST http://localhost:8080/users/goals \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "daily_calories": 2000,
    "protein_ratio": 30.0,
    "fat_ratio": 30.0,
    "carb_ratio": 40.0
  }'
```

### Получение целей пользователя

```bash
curl -X GET http://localhost:8080/users/goals \
  -H "Authorization: Bearer $TOKEN"
```

## Работа со списками покупок

### Получение списка покупок для меню

```bash
curl -X GET http://localhost:8080/shopping-list/1 \
  -H "Authorization: Bearer $TOKEN"
```

## Примеры на JavaScript/TypeScript

### Использование с Axios

```typescript
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Добавление токена к запросам
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Получение токена
async function getTestToken() {
  const response = await api.get('/auth/test');
  localStorage.setItem('token', response.data.token);
  return response.data.token;
}

// Генерация меню
async function generateMenu() {
  const response = await api.post('/menus/generate', {
    user_id: 1,
    target_calories: 2000,
    diet_type: 'vegetarian',
    consider_pantry: true,
    pantry_importance: 'prefer',
  });
  return response.data;
}

// Получение списка рецептов
async function getRecipes() {
  const response = await api.get('/recipes');
  return response.data;
}

// Добавление продукта в кладовую
async function addPantryItem(name: string, quantity: number, unit: string) {
  const response = await api.post('/pantry', {
    name,
    quantity,
    unit,
  });
  return response.data;
}
```

### Использование с Fetch API

```javascript
const API_URL = 'http://localhost:8080';

// Получение токена
async function getTestToken() {
  const response = await fetch(`${API_URL}/auth/test`);
  const data = await response.json();
  localStorage.setItem('token', data.token);
  return data.token;
}

// Генерация меню
async function generateMenu() {
  const token = localStorage.getItem('token');
  const response = await fetch(`${API_URL}/menus/generate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify({
      user_id: 1,
      target_calories: 2000,
      consider_pantry: true,
    }),
  });
  return await response.json();
}

// Получение списка рецептов
async function getRecipes() {
  const response = await fetch(`${API_URL}/recipes`);
  return await response.json();
}
```

## Примеры на Python

```python
import requests

API_URL = "http://localhost:8080"

# Получение токена
def get_test_token():
    response = requests.get(f"{API_URL}/auth/test")
    data = response.json()
    return data["token"]

# Генерация меню
def generate_menu(token):
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json",
    }
    data = {
        "user_id": 1,
        "target_calories": 2000,
        "diet_type": "vegetarian",
        "consider_pantry": True,
        "pantry_importance": "prefer",
    }
    response = requests.post(
        f"{API_URL}/menus/generate",
        headers=headers,
        json=data
    )
    return response.json()

# Получение списка рецептов
def get_recipes():
    response = requests.get(f"{API_URL}/recipes")
    return response.json()

# Добавление продукта в кладовую
def add_pantry_item(token, name, quantity, unit):
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json",
    }
    data = {
        "name": name,
        "quantity": quantity,
        "unit": unit,
    }
    response = requests.post(
        f"{API_URL}/pantry",
        headers=headers,
        json=data
    )
    return response.json()

# Использование
if __name__ == "__main__":
    token = get_test_token()
    print(f"Token: {token}")
    
    recipes = get_recipes()
    print(f"Recipes: {len(recipes)}")
    
    menu = generate_menu(token)
    print(f"Menu ID: {menu['id']}")
    
    pantry_item = add_pantry_item(token, "Яйца", 10, "шт")
    print(f"Added: {pantry_item['name']}")
```

## Примеры на Go

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

const API_URL = "http://localhost:8080"

type AuthResponse struct {
    User  User   `json:"user"`
    Token string `json:"token"`
}

type User struct {
    ID        int    `json:"id"`
    TelegramID int64 `json:"telegram_id"`
    Username  string `json:"username"`
}

type MenuGenerateRequest struct {
    UserID           int     `json:"user_id"`
    TargetCalories   int     `json:"target_calories"`
    DietType         string  `json:"diet_type,omitempty"`
    ConsiderPantry   bool    `json:"consider_pantry"`
    PantryImportance string  `json:"pantry_importance"`
}

func getTestToken() (string, error) {
    resp, err := http.Get(API_URL + "/auth/test")
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var authResp AuthResponse
    if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
        return "", err
    }

    return authResp.Token, nil
}

func generateMenu(token string, req MenuGenerateRequest) (map[string]interface{}, error) {
    jsonData, err := json.Marshal(req)
    if err != nil {
        return nil, err
    }

    httpReq, err := http.NewRequest("POST", API_URL+"/menus/generate", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }

    httpReq.Header.Set("Authorization", "Bearer "+token)
    httpReq.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(httpReq)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }

    return result, nil
}

func main() {
    token, err := getTestToken()
    if err != nil {
        fmt.Printf("Error getting token: %v\n", err)
        return
    }
    fmt.Printf("Token: %s\n", token)

    req := MenuGenerateRequest{
        UserID:           1,
        TargetCalories:   2000,
        DietType:         "vegetarian",
        ConsiderPantry:   true,
        PantryImportance: "prefer",
    }

    menu, err := generateMenu(token, req)
    if err != nil {
        fmt.Printf("Error generating menu: %v\n", err)
        return
    }

    fmt.Printf("Menu ID: %.0f\n", menu["id"])
}
```

## Типичные сценарии использования

### Сценарий 1: Генерация меню с учетом кладовой

1. Получить список продуктов в кладовой
2. Сгенерировать меню с `consider_pantry: true` и `pantry_importance: "prefer"`
3. Получить список покупок для недостающих ингредиентов

```bash
# 1. Получить кладовую
curl -X GET http://localhost:8080/pantry -H "Authorization: Bearer $TOKEN"

# 2. Сгенерировать меню
curl -X POST http://localhost:8080/menus/generate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "target_calories": 2000,
    "consider_pantry": true,
    "pantry_importance": "prefer"
  }'

# 3. Получить список покупок (используйте menu_id из ответа)
curl -X GET http://localhost:8080/shopping-list/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Сценарий 2: Вегетарианское меню без аллергенов

```bash
curl -X POST http://localhost:8080/menus/generate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "target_calories": 1800,
    "diet_type": "vegetarian",
    "allergies": ["nuts", "dairy"],
    "max_total_time": 45,
    "consider_pantry": false
  }'
```

### Сценарий 3: Быстрое меню (минимум времени)

```bash
curl -X POST http://localhost:8080/menus/generate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "target_calories": 2000,
    "max_total_time": 30,
    "consider_pantry": true,
    "pantry_importance": "strict"
  }'
```

## Обработка ошибок

### Пример обработки ошибок в JavaScript

```javascript
async function generateMenuWithErrorHandling() {
  try {
    const response = await fetch(`${API_URL}/menus/generate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify({
        user_id: 1,
        target_calories: 2000,
      }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Ошибка при генерации меню');
    }

    const menu = await response.json();
    return menu;
  } catch (error) {
    console.error('Ошибка:', error.message);
    throw error;
  }
}
```

### Типичные ошибки

**401 Unauthorized:**
```json
{
  "error": "Требуется заголовок авторизации"
}
```
Решение: Проверьте наличие и корректность токена.

**400 Bad Request:**
```json
{
  "error": "Неверное тело запроса"
}
```
Решение: Проверьте формат JSON и обязательные поля.

**500 Internal Server Error:**
```json
{
  "error": "не найдена подходящая комбинация меню"
}
```
Решение: Попробуйте изменить параметры (увеличить допуск по калориям, добавить больше продуктов в кладовую, убрать строгие ограничения).

