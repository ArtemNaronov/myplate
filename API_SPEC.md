# API Specification - MyPlateService

## Базовый URL
```
http://localhost:8080
```

## Аутентификация

Большинство endpoints требуют JWT токен в заголовке:
```
Authorization: Bearer <token>
```

JWT токен содержит:
```json
{
  "user_id": 1,
  "role": "user" | "admin",
  "exp": 1234567890
}
```

---

## 1. Генерация меню на неделю

### `GET /menu/weekly`

Генерирует меню на 7 дней с учетом количества людей, анти-повторов и баланса БЖУ.

**Требует авторизации:** Да

**Query параметры:**

| Параметр | Тип | Обязательный | Описание |
|----------|-----|--------------|----------|
| `adults` | int | Да | Количество взрослых |
| `children` | int | Нет (по умолчанию 0) | Количество детей |
| `diet_type` | string | Нет | Тип диеты: "vegetarian", "vegan", "gluten-free" |
| `allergies` | string | Нет | Список аллергенов через запятую: "nuts,dairy,eggs" |
| `max_total_time` | int | Нет | Максимальное время приготовления в минутах |
| `max_time_per_meal` | int | Нет | Максимальное время на одно блюдо в минутах |
| `consider_pantry` | bool | Нет | Учитывать кладовую (по умолчанию false) |
| `pantry_importance` | string | Нет | Важность кладовой: "ignore", "prefer", "strict" (по умолчанию "prefer") |

**Пример запроса:**
```bash
curl -X GET "http://localhost:8080/menu/weekly?adults=2&children=1&diet_type=vegetarian&allergies=nuts" \
  -H "Authorization: Bearer $TOKEN"
```

**Ответ:**
```json
{
  "week": [
    {
      "day": 1,
      "breakfast": {
        "id": 1,
        "name": "Омлет",
        "description": "Классический омлет",
        "calories": 350,
        "proteins": 20.0,
        "fats": 18.0,
        "carbs": 28.0,
        "cooking_time": 10,
        "servings": 1,
        "meal_type": "breakfast",
        "ingredients": [
          {"name": "Яйца", "quantity": 2, "unit": "шт"}
        ],
        "instructions": ["Шаг 1", "Шаг 2"]
      },
      "lunch": {
        "id": 5,
        "name": "Салат",
        ...
      },
      "dinner": {
        "id": 8,
        "name": "Рыба",
        ...
      },
      "totalCalories": 1950,
      "totalProteins": 95.0,
      "totalFats": 70.0,
      "totalCarbs": 180.0,
      "totalTime": 45,
      "ingredients_used": [...],
      "missing_ingredients": [...]
    },
    ...
  ]
}
```

**Особенности:**
- Калорийная цель дня: `adults * 2000 + children * 1400`
- Распределение калорий: завтрак 25%, обед 40%, ужин 35%
- Анти-повторы: один рецепт не используется 3 дня подряд
- Автоматическая оптимизация баланса БЖУ (25% белки, 30% жиры, 45% углеводы)
- Пересчет ингредиентов: `totalServings = adults + children * 0.7`

---

## 2. Генерация дневного меню (обновлено)

### `POST /menus/generate`

Генерирует меню на один день с учетом количества людей.

**Требует авторизации:** Да

**Request Body:**
```json
{
  "user_id": 1,
  "target_calories": 2000,
  "adults": 2,
  "children": 1,
  "diet_type": "vegetarian",
  "allergies": ["nuts"],
  "max_total_time": 60,
  "max_time_per_meal": 30,
  "consider_pantry": true,
  "pantry_importance": "prefer"
}
```

**Параметры:**
- `adults` (int, опционально, по умолчанию 1) - количество взрослых
- `children` (int, опционально, по умолчанию 0) - количество детей

**Формула расчета ингредиентов:**
```
totalServings = adults + children * 0.7
ingredient.amount = baseAmount * (totalServings / recipeServings)
```

---

## 3. Админ endpoints

### `POST /admin/recipes`

Создает новый рецепт.

**Требует авторизации:** Да (роль `admin`)

**Request Body:**
```json
{
  "title": "Название рецепта",
  "description": "Описание рецепта",
  "tags": ["breakfast", "vegetarian", "eggs"],
  "ingredients": [
    {
      "name": "Яйца",
      "amount": 2,
      "unit": "шт"
    },
    {
      "name": "Молоко",
      "amount": 100,
      "unit": "мл"
    }
  ],
  "calories": 350,
  "proteins": 20.0,
  "fats": 18.0,
  "carbs": 28.0,
  "cooking_time": 10,
  "servings": 1,
  "instructions": ["Шаг 1", "Шаг 2"]
}
```

**Теги:**
- `meal_type`: "breakfast", "lunch", "dinner", "snack"
- `diet_type`: "vegetarian", "vegan", "gluten-free"
- `allergens`: "eggs", "dairy", "nuts", "fish", "gluten"

**Response:**
```json
{
  "id": 1,
  "name": "Название рецепта",
  "description": "Описание рецепта",
  "calories": 350,
  "proteins": 20.0,
  "fats": 18.0,
  "carbs": 28.0,
  "cooking_time": 10,
  "servings": 1,
  "meal_type": "breakfast",
  "diet_type": ["vegetarian"],
  "allergens": ["eggs"],
  "ingredients": [...],
  "instructions": [...],
  "created_at": "2025-12-03T10:00:00Z",
  "updated_at": "2025-12-03T10:00:00Z"
}
```

**Ошибки:**
- `400` - Неверное тело запроса
- `403` - Требуются права администратора
- `409` - Рецепт с таким названием уже существует

---

### `POST /admin/recipes/import`

Импортирует несколько рецептов из JSON.

**Требует авторизации:** Да (роль `admin`)

**Request Body:**
```json
{
  "recipes": [
    {
      "title": "Рецепт 1",
      "description": "Описание",
      "tags": ["breakfast"],
      "ingredients": [
        {"name": "Ингредиент", "amount": 100, "unit": "г"}
      ],
      "calories": 300,
      "proteins": 15.0,
      "fats": 10.0,
      "carbs": 40.0
    },
    {
      "title": "Рецепт 2",
      ...
    }
  ]
}
```

**Response:**
```json
{
  "imported": 2,
  "failed": 0,
  "errors": []
}
```

**Особенности:**
- Использует транзакцию (все или ничего)
- Проверяет дубликаты по названию (case-insensitive)
- Пропускает дубликаты, но продолжает импорт остальных

---

### `GET /admin/recipes/export`

Экспортирует все рецепты в JSON.

**Требует авторизации:** Да (роль `admin`)

**Response:**
```json
{
  "recipes": [
    {
      "title": "Название",
      "description": "Описание",
      "tags": ["breakfast", "vegetarian"],
      "ingredients": [
        {"name": "Ингредиент", "amount": 100, "unit": "г"}
      ],
      "calories": 350,
      "proteins": 20.0,
      "fats": 18.0,
      "carbs": 28.0,
      "cooking_time": 10,
      "servings": 1,
      "instructions": ["Шаг 1", "Шаг 2"]
    }
  ]
}
```

---

## Коды ошибок

| Код | Описание |
|-----|----------|
| `200 OK` | Успешный запрос |
| `201 Created` | Ресурс создан |
| `400 Bad Request` | Неверный формат запроса |
| `401 Unauthorized` | Требуется авторизация или неверный токен |
| `403 Forbidden` | Доступ запрещен (требуется роль admin) |
| `404 Not Found` | Ресурс не найден |
| `409 Conflict` | Конфликт (например, дубликат рецепта) |
| `500 Internal Server Error` | Внутренняя ошибка сервера |

---

## Примеры использования

### Генерация недельного меню для семьи
```bash
curl -X GET "http://localhost:8080/menu/weekly?adults=2&children=2&consider_pantry=true&pantry_importance=prefer" \
  -H "Authorization: Bearer $TOKEN"
```

### Создание рецепта (админ)
```bash
curl -X POST http://localhost:8080/admin/recipes \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Омлет",
    "tags": ["breakfast", "vegetarian"],
    "ingredients": [{"name": "Яйца", "amount": 3, "unit": "шт"}],
    "calories": 250,
    "proteins": 18.0,
    "fats": 15.0,
    "carbs": 2.0,
    "cooking_time": 10,
    "servings": 1
  }'
```

### Импорт рецептов (админ)
```bash
curl -X POST http://localhost:8080/admin/recipes/import \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d @recipes.json
```

### Экспорт рецептов (админ)
```bash
curl -X GET http://localhost:8080/admin/recipes/export \
  -H "Authorization: Bearer $ADMIN_TOKEN" > recipes_export.json
```

