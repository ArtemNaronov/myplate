# Руководство по установке

## Предварительные требования

- **Docker** и **Docker Compose** установлены
- **Make** (опционально, но рекомендуется)
- **Go 1.22+** (для локальной разработки backend)
- **Node.js 20+** (для локальной разработки frontend)

## Быстрый старт с Docker

1. **Запустите все сервисы:**
   ```bash
   make up
   # или
   docker-compose up -d
   ```

2. **Дождитесь готовности сервисов** (особенно инициализации базы данных)

3. **Доступ к приложению:**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - PgAdmin: http://localhost:5050 (admin@myplate.com / admin)

4. **Проверка работоспособности:**
   ```bash
   # Проверка API
   curl http://localhost:8080/health
   
   # Получение тестового токена
   curl http://localhost:8080/auth/test
   ```

## Локальная разработка

### Backend

1. **Перейдите в директорию backend:**
   ```bash
   cd backend
   ```

2. **Установите зависимости:**
   ```bash
   go mod download
   ```

3. **Создайте файл `.env`:**
   ```bash
   cat > .env << EOF
   DATABASE_URL=postgres://myplate:myplate@localhost:5432/myplate?sslmode=disable
   JWT_SECRET=your-secret-key-change-in-production
   TELEGRAM_BOT_TOKEN=your-telegram-bot-token
   PORT=8080
   EOF
   ```

4. **Запустите базу данных:**
   ```bash
   docker-compose up -d postgres
   ```

5. **Примените миграции:**
   Миграции применяются автоматически при первом запуске PostgreSQL контейнера.

   Для ручного применения:
   ```bash
   psql -U myplate -d myplate -f ../sql/migrations/001_initial_schema.sql
   ```

6. **Загрузите тестовые данные:**
   ```bash
   psql -U myplate -d myplate -f ../sql/seed.sql
   ```

7. **Запустите backend:**
   ```bash
   go run cmd/api/main.go
   ```

### Frontend

1. **Перейдите в директорию frontend:**
   ```bash
   cd frontend
   ```

2. **Установите зависимости:**
   ```bash
   npm install
   ```

3. **Создайте файл `.env.local`:**
   ```bash
   cat > .env.local << EOF
   NEXT_PUBLIC_API_URL=http://localhost:8080
   NEXT_PUBLIC_TELEGRAM_BOT_NAME=your-bot-name
   EOF
   ```

4. **Запустите frontend:**
   ```bash
   npm run dev
   ```

## Настройка базы данных

База данных автоматически инициализируется с:
- Схемой из `sql/migrations/001_initial_schema.sql`
- Тестовыми данными из `sql/seed.sql` (10 рецептов, 2 пользователя, 41 продукт в кладовой)

### Ручная загрузка данных

```bash
# Через Docker
make seed
# или
docker-compose exec postgres psql -U myplate -d myplate -f /docker-entrypoint-initdb.d/seed.sql

# Напрямую
psql -U myplate -d myplate -f sql/seed.sql
```

### Сброс базы данных

```bash
# Остановить контейнеры
docker-compose down

# Удалить volumes (удалит все данные!)
docker-compose down -v

# Запустить заново
docker-compose up -d
```

## Тестирование API

### Проверка работоспособности
```bash
curl http://localhost:8080/health
```

### Получение списка рецептов
```bash
curl http://localhost:8080/recipes
```

### Тестовая авторизация
```bash
curl http://localhost:8080/auth/test
```

Сохраните токен:
```bash
export TOKEN="<ваш-токен>"
```

### Генерация меню
```bash
curl -X POST http://localhost:8080/menus/generate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "target_calories": 2000,
    "consider_pantry": true,
    "pantry_importance": "prefer"
  }'
```

Подробные примеры использования API см. в [API_EXAMPLES.md](./API_EXAMPLES.md).

## Интеграция с Telegram WebApp

1. **Создайте Telegram бота:**
   - Откройте [@BotFather](https://t.me/botfather) в Telegram
   - Отправьте команду `/newbot`
   - Следуйте инструкциям для создания бота
   - Сохраните полученный токен

2. **Настройте backend:**
   - Добавьте токен в `backend/.env`:
     ```
     TELEGRAM_BOT_TOKEN=your-bot-token-here
     ```

3. **Настройте frontend:**
   - Добавьте имя бота в `frontend/.env.local`:
     ```
     NEXT_PUBLIC_TELEGRAM_BOT_NAME=your-bot-name
     ```

4. **Настройте WebApp URL:**
   - В настройках бота через BotFather установите WebApp URL
   - URL должен указывать на ваш frontend (например, `https://yourdomain.com`)

5. **Проверка:**
   - Frontend автоматически определяет, запущен ли он в Telegram WebApp
   - Если нет, используется тестовая авторизация

## Решение проблем

### Проблемы с подключением к базе данных

**Симптомы:**
- Ошибка "connection refused"
- Ошибка "database does not exist"

**Решения:**
1. Проверьте, запущен ли контейнер PostgreSQL:
   ```bash
   docker-compose ps
   ```

2. Проверьте логи PostgreSQL:
   ```bash
   docker-compose logs postgres
   ```

3. Проверьте строку подключения в `backend/.env`:
   ```
   DATABASE_URL=postgres://myplate:myplate@localhost:5432/myplate?sslmode=disable
   ```
   Для Docker используйте `postgres` вместо `localhost`:
   ```
   DATABASE_URL=postgres://myplate:myplate@postgres:5432/myplate?sslmode=disable
   ```

4. Убедитесь, что база данных инициализирована:
   ```bash
   docker-compose exec postgres psql -U myplate -d myplate -c "\dt"
   ```

### Конфликты портов

**Симптомы:**
- Ошибка "port already in use"
- Сервис не запускается

**Решения:**
1. Измените порты в `docker-compose.yml`:
   ```yaml
   ports:
     - "3001:3000"  # Frontend
     - "8081:8080"  # Backend
     - "5433:5432"  # PostgreSQL
   ```

2. Обновите `frontend/.env.local`:
   ```
   NEXT_PUBLIC_API_URL=http://localhost:8081
   ```

3. Остановите конфликтующие сервисы:
   ```bash
   # Найти процесс на порту
   lsof -i :3000
   # Остановить процесс
   kill <PID>
   ```

### Проблемы с миграциями

**Симптомы:**
- Ошибка "relation does not exist"
- Таблицы не создаются

**Решения:**
1. Проверьте наличие файла миграции:
   ```bash
   ls -la sql/migrations/
   ```

2. Проверьте права доступа:
   ```bash
   chmod 644 sql/migrations/001_initial_schema.sql
   ```

3. Примените миграции вручную:
   ```bash
   docker-compose exec postgres psql -U myplate -d myplate -f /docker-entrypoint-initdb.d/001_initial_schema.sql
   ```

4. Проверьте логи PostgreSQL при запуске:
   ```bash
   docker-compose logs postgres | grep -i error
   ```

### Проблемы с авторизацией

**Симптомы:**
- Ошибка "Требуется заголовок авторизации"
- 401 Unauthorized

**Решения:**
1. Получите новый токен:
   ```bash
   curl http://localhost:8080/auth/test
   ```

2. Проверьте заголовок Authorization:
   ```bash
   curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8080/menus
   ```

3. Убедитесь, что `JWT_SECRET` установлен в `backend/.env`

### Проблемы с генерацией меню

**Симптомы:**
- Ошибка "не найдена подходящая комбинация меню"

**Решения:**
1. Добавьте больше продуктов в кладовую:
   ```bash
   curl -X POST http://localhost:8080/pantry \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"name": "Яйца", "quantity": 10, "unit": "шт"}'
   ```

2. Уменьшите строгость ограничений:
   - Увеличьте `max_total_time`
   - Уберите `diet_type` или `allergies`
   - Измените `pantry_importance` на "prefer" или "ignore"

3. Проверьте, что в базе есть рецепты:
   ```bash
   curl http://localhost:8080/recipes | jq length
   ```

### Проблемы с frontend

**Симптомы:**
- Страница не загружается
- Ошибки в консоли браузера

**Решения:**
1. Проверьте, что frontend запущен:
   ```bash
   docker-compose ps frontend
   # или
   curl http://localhost:3000
   ```

2. Проверьте логи frontend:
   ```bash
   docker-compose logs frontend
   ```

3. Очистите кэш браузера (Ctrl+Shift+R или Cmd+Shift+R)

4. Проверьте переменные окружения:
   ```bash
   docker-compose exec frontend env | grep NEXT_PUBLIC
   ```

5. Пересоберите frontend:
   ```bash
   docker-compose up -d --build frontend
   ```

## Структура проекта

```
MyPlate/
├── backend/              # Go API сервер
│   ├── cmd/api/         # Точка входа приложения
│   ├── internal/        # Внутренние пакеты
│   │   ├── handlers/    # HTTP обработчики
│   │   ├── services/    # Бизнес-логика
│   │   ├── repositories/# Доступ к данным
│   │   ├── models/      # Доменные модели
│   │   └── middleware/  # Middleware
│   └── pkg/             # Публичные пакеты
│       └── database/     # Подключение к БД
├── frontend/            # Next.js приложение
│   ├── app/             # App Router страницы
│   ├── components/      # React компоненты
│   └── lib/             # Утилиты
├── sql/                 # SQL файлы
│   ├── migrations/      # Миграции схемы
│   └── seed.sql        # Тестовые данные
├── docker-compose.yml   # Конфигурация Docker
├── Dockerfile.backend   # Dockerfile для backend
├── Dockerfile.frontend  # Dockerfile для frontend
├── Makefile             # Команды для сборки
├── README.md            # Основная документация
├── API_EXAMPLES.md      # Примеры использования API
└── SETUP.md             # Это руководство
```

## Следующие шаги

1. ✅ Установите и запустите приложение
2. ✅ Протестируйте API через curl или Postman
3. ✅ Настройте Telegram бота (опционально)
4. ✅ Добавьте свои рецепты в базу данных
5. ✅ Настройте переменные окружения для продакшена
6. ✅ Настройте CI/CD pipeline
7. ✅ Добавьте мониторинг и логирование

## Дополнительные ресурсы

- [Основная документация](./README.md) - Полное описание проекта
- [Примеры использования API](./API_EXAMPLES.md) - Практические примеры
- [Makefile](./Makefile) - Доступные команды

## Получение помощи

Если у вас возникли проблемы:

1. Проверьте логи: `docker-compose logs`
2. Убедитесь, что все переменные окружения установлены
3. Проверьте версии зависимостей
4. Создайте issue в репозитории с описанием проблемы
