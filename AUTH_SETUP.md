# Настройка полноценной авторизации

## Что было добавлено

### Backend

1. **Миграция базы данных** (`sql/migrations/002_add_email_password.sql`)
   - Добавлены поля `email` и `password_hash` в таблицу `users`
   - `telegram_id` теперь может быть NULL
   - Добавлен индекс для `email`
   - Добавлена проверка: пользователь должен иметь либо `telegram_id`, либо `email`

2. **Обновлена модель User**
   - Добавлены поля `Email` и `PasswordHash`
   - `TelegramID` теперь указатель (`*int64`), может быть `nil`

3. **Расширен UserRepository**
   - `GetByEmail()` - получение пользователя по email
   - `Create()` - создание пользователя с email и паролем
   - `UpdatePassword()` - обновление пароля
   - `UpdateProfile()` - обновление профиля

4. **Расширен AuthService**
   - `Register()` - регистрация нового пользователя
   - `Login()` - вход по email и паролю
   - `GetUserProfile()` - получение профиля пользователя
   - `UpdatePassword()` - обновление пароля
   - `UpdateProfile()` - обновление профиля
   - Используется `bcrypt` для хеширования паролей

5. **Расширен AuthHandler**
   - `POST /auth/register` - регистрация
   - `POST /auth/login` - вход
   - `GET /auth/profile` - получение профиля (требует авторизации)
   - `PUT /auth/profile` - обновление профиля (требует авторизации)
   - `PUT /auth/password` - обновление пароля (требует авторизации)

### Frontend

1. **Страница регистрации** (`/auth/register`)
   - Форма с полями: email, пароль, имя, фамилия
   - Валидация пароля (минимум 6 символов)
   - Автоматический вход после регистрации

2. **Страница входа** (`/auth/login`)
   - Форма с полями: email, пароль
   - Обработка ошибок авторизации
   - Редирект на главную после успешного входа

3. **Страница профиля** (`/profile`)
   - Просмотр информации о пользователе
   - Обновление имени и фамилии
   - Изменение пароля
   - Выход из аккаунта

4. **Обновлена навигация**
   - Кнопка "Войти" для неавторизованных пользователей
   - Ссылка "Профиль" и кнопка "Выйти" для авторизованных

5. **Обновлен TelegramProvider**
   - Проверка валидности токена при загрузке
   - Удаление невалидного токена
   - Поддержка обычной авторизации (не только Telegram)

## Применение миграции

### Вариант 1: Через psql

```bash
psql $DATABASE_URL -f sql/migrations/002_add_email_password.sql
```

### Вариант 2: Через Docker

```bash
docker exec -i <postgres_container> psql -U <user> -d <database> < sql/migrations/002_add_email_password.sql
```

### Вариант 3: Вручную

Выполните SQL команды из файла `sql/migrations/002_add_email_password.sql` в вашей базе данных.

## Установка зависимостей

Backend теперь требует пакет `golang.org/x/crypto` для хеширования паролей:

```bash
cd backend
go get golang.org/x/crypto/bcrypt
go mod tidy
```

## Использование API

### Регистрация

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "Иван",
    "last_name": "Иванов"
  }'
```

Ответ:
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "user@example.com",
    "first_name": "Иван",
    "last_name": "Иванов"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Вход

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Получение профиля

```bash
curl -X GET http://localhost:8080/auth/profile \
  -H "Authorization: Bearer <token>"
```

### Обновление профиля

```bash
curl -X PUT http://localhost:8080/auth/profile \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Петр",
    "last_name": "Петров"
  }'
```

### Обновление пароля

```bash
curl -X PUT http://localhost:8080/auth/password \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "password123",
    "new_password": "newpassword456"
  }'
```

## Безопасность

1. ✅ Пароли хешируются с помощью `bcrypt` (стоимость по умолчанию)
2. ✅ Минимальная длина пароля - 6 символов
3. ✅ JWT токены имеют срок действия 7 дней
4. ✅ Email должен быть уникальным
5. ✅ Все защищенные endpoints требуют валидный JWT токен

## Совместимость

Система поддерживает оба способа авторизации:
- **Telegram WebApp** - через `/auth/telegram`
- **Email/Password** - через `/auth/register` и `/auth/login`

Пользователь может иметь либо `telegram_id`, либо `email`, либо оба.

## Тестирование

1. Примените миграцию
2. Установите зависимости: `cd backend && go get golang.org/x/crypto/bcrypt && go mod tidy`
3. Запустите backend: `cd backend && go run cmd/api/main.go`
4. Откройте frontend: `cd frontend && npm run dev`
5. Перейдите на `/auth/register` и создайте аккаунт
6. Войдите через `/auth/login`
7. Проверьте профиль на `/profile`

