# Новые функции MyPlateService

## 1. Генерация меню на неделю

### Endpoint
```
GET /menu/weekly
```

### Параметры (query)
- `adults` (int, обязательный) - количество взрослых
- `children` (int, опционально, по умолчанию 0) - количество детей
- `diet_type` (string, опционально) - тип диеты
- `allergies` (string, опционально) - список аллергенов через запятую
- `max_total_time` (int, опционально) - максимальное время приготовления
- `max_time_per_meal` (int, опционально) - максимальное время на одно блюдо
- `consider_pantry` (bool, опционально) - учитывать кладовую
- `pantry_importance` (string, опционально) - важность кладовой: "ignore", "prefer", "strict"

### Пример запроса
```bash
curl -X GET "http://localhost:8080/menu/weekly?adults=2&children=1" \
  -H "Authorization: Bearer $TOKEN"
```

### Ответ
```json
{
  "days": [
    {
      "date": "2025-12-02T00:00:00Z",
      "total_calories": 1950,
      "total_time": 45,
      "meals": [...],
      "ingredients_used": [...],
      "missing_ingredients": [...]
    },
    ...
  ]
}
```

### Особенности
- Генерирует меню на 7 дней
- Использует целевые калории из `user_goals`
- Учитывает количество людей (adults + children * 0.7)
- Пересчитывает ингредиенты с учетом количества порций

## 2. Параметр количества человек

### Обновленные endpoints
Все endpoints генерации меню теперь поддерживают параметры:
- `adults` (int, по умолчанию 1) - количество взрослых
- `children` (int, по умолчанию 0) - количество детей

### Формула расчета
```
totalServings = adults + children * 0.7
ingredient.amount = baseAmount * (totalServings / recipeServings)
```

### Пример
```json
{
  "user_id": 1,
  "target_calories": 2000,
  "adults": 2,
  "children": 1,
  "consider_pantry": true
}
```

## 3. Админ-доступ

### Роли пользователей
- `user` - обычный пользователь (по умолчанию)
- `admin` - администратор

### Middleware
- `AuthMiddleware` - проверяет JWT токен и устанавливает `user_id` и `user_role` в контекст
- `AdminMiddleware` - проверяет, что `user_role == "admin"`

### Endpoints для администраторов

#### POST /admin/recipes
Создание нового рецепта.

**Request:**
```json
{
  "title": "Название рецепта",
  "description": "Описание",
  "tags": ["breakfast", "vegetarian"],
  "ingredients": [
    {
      "name": "Яйца",
      "amount": 2,
      "unit": "шт"
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

#### POST /admin/recipes/import
Импорт нескольких рецептов из JSON.

**Request:**
```json
{
  "recipes": [
    {
      "title": "Рецепт 1",
      ...
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

#### GET /admin/recipes/export
Экспорт всех рецептов в JSON.

**Response:**
```json
{
  "recipes": [
    {
      "title": "Название",
      "description": "Описание",
      "tags": ["breakfast", "vegetarian"],
      "ingredients": [...],
      "calories": 350,
      "proteins": 20.0,
      "fats": 18.0,
      "carbs": 28.0,
      "cooking_time": 10,
      "servings": 1,
      "instructions": [...]
    }
  ]
}
```

## 4. Миграции

### 003_add_user_role.sql
Добавляет поле `role` в таблицу `users`:
- По умолчанию: `'user'`
- Возможные значения: `'user'`, `'admin'`
- Индекс для быстрого поиска

### Применение миграции
```bash
docker exec -i myplateservice-postgres psql -U myplateservice -d myplateservice < sql/migrations/003_add_user_role.sql
```

## 5. Обновления JWT

JWT токены теперь включают поле `role`:
```json
{
  "user_id": 1,
  "role": "admin",
  "exp": 1234567890
}
```

## 6. Назначение роли администратора

Для назначения роли администратора пользователю:
```sql
UPDATE users SET role = 'admin' WHERE id = <user_id>;
```

Или при регистрации (если нужно):
```sql
INSERT INTO users (email, password_hash, first_name, last_name, role)
VALUES ('admin@example.com', '$2a$10$...', 'Admin', 'User', 'admin');
```

## Примеры использования

### Генерация недельного меню для семьи
```bash
curl -X GET "http://localhost:8080/menu/weekly?adults=2&children=2&consider_pantry=true" \
  -H "Authorization: Bearer $TOKEN"
```

### Создание рецепта (только для админов)
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

### Экспорт всех рецептов
```bash
curl -X GET http://localhost:8080/admin/recipes/export \
  -H "Authorization: Bearer $ADMIN_TOKEN" > recipes.json
```

### Импорт рецептов
```bash
curl -X POST http://localhost:8080/admin/recipes/import \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d @recipes.json
```

