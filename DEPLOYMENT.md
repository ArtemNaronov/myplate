# Руководство по деплою

## Обзор архитектуры

Приложение состоит из трех компонентов:
1. **Frontend** (Next.js) - можно разместить на GitHub Pages
2. **Backend** (Go API) - нужен отдельный хостинг
3. **База данных** (PostgreSQL) - нужен отдельный хостинг

## Деплой Frontend на GitHub Pages

### Предварительные требования

1. Репозиторий на GitHub
2. Включенные GitHub Pages в настройках репозитория:
   - Settings → Pages → Source: GitHub Actions

### Настройка

1. **Настройте basePath в `next.config.js`** (если репозиторий не в корне):
   ```javascript
   basePath: process.env.NODE_ENV === 'production' ? '/MyPlate' : '',
   assetPrefix: process.env.NODE_ENV === 'production' ? '/MyPlate' : '',
   ```

2. **Добавьте секреты в GitHub** (Settings → Secrets and variables → Actions):
   - `NEXT_PUBLIC_API_URL` - URL вашего backend API (например, `https://api.myplate.com`)

3. **Запустите workflow**:
   - При push в `main` ветку автоматически запустится деплой
   - Или вручную через Actions → Deploy to GitHub Pages → Run workflow

### Результат

Frontend будет доступен по адресу:
- `https://<username>.github.io/<repository-name>/`

## Деплой Backend

GitHub Pages не поддерживает backend, поэтому нужен отдельный хостинг.

### Варианты хостинга Backend

#### 1. Railway (рекомендуется)

1. Создайте аккаунт на [Railway](https://railway.app)
2. Создайте новый проект
3. Подключите GitHub репозиторий
4. Добавьте PostgreSQL сервис
5. Настройте переменные окружения:
   ```
   DATABASE_URL=<connection_string>
   JWT_SECRET=<your-secret>
   TELEGRAM_BOT_TOKEN=<your-token>
   PORT=8080
   ```
6. Railway автоматически определит Go проект и задеплоит

#### 2. Render

1. Создайте аккаунт на [Render](https://render.com)
2. Создайте новый Web Service
3. Подключите GitHub репозиторий
4. Настройки:
   - Build Command: `cd backend && go build -o bin/api ./cmd/api`
   - Start Command: `./backend/bin/api`
5. Добавьте PostgreSQL базу данных
6. Настройте переменные окружения

#### 3. Fly.io

1. Установите [flyctl](https://fly.io/docs/getting-started/installing-flyctl/)
2. Создайте `fly.toml` в корне проекта:
   ```toml
   app = "myplate-api"
   primary_region = "iad"

   [build]
     builder = "paketobuildpacks/builder:base"

   [[services]]
     internal_port = 8080
     protocol = "tcp"
   ```
3. Запустите: `fly launch`

#### 4. Heroku

1. Создайте `Procfile` в корне:
   ```
   web: cd backend && ./bin/api
   ```
2. Установите Heroku CLI
3. `heroku create myplate-api`
4. `git push heroku main`

### Настройка базы данных

#### Варианты:

1. **Railway PostgreSQL** - автоматически при создании проекта
2. **Supabase** - бесплатный PostgreSQL хостинг
3. **Neon** - serverless PostgreSQL
4. **Render PostgreSQL** - встроенная база данных
5. **AWS RDS** - для production

### Применение миграций

После создания базы данных:

```bash
# Через psql
psql $DATABASE_URL -f sql/migrations/001_initial_schema.sql

# Или через Docker
docker exec -i <container> psql -U user -d dbname < sql/migrations/001_initial_schema.sql
```

## Полный процесс деплоя

### Шаг 1: Деплой базы данных

1. Создайте PostgreSQL базу данных на выбранном хостинге
2. Примените миграции:
   ```bash
   psql $DATABASE_URL -f sql/migrations/001_initial_schema.sql
   psql $DATABASE_URL -f sql/seed.sql
   ```

### Шаг 2: Деплой Backend

1. Выберите хостинг (Railway/Render/Fly.io)
2. Подключите репозиторий
3. Настройте переменные окружения:
   - `DATABASE_URL`
   - `JWT_SECRET`
   - `TELEGRAM_BOT_TOKEN`
   - `PORT=8080`
4. Деплой запустится автоматически

### Шаг 3: Деплой Frontend

1. В GitHub репозитории:
   - Settings → Pages → Source: GitHub Actions
2. Добавьте секрет `NEXT_PUBLIC_API_URL` с URL вашего backend
3. Push в `main` ветку
4. Frontend автоматически задеплоится

### Шаг 4: Настройка CORS

Убедитесь, что backend разрешает запросы с вашего frontend домена:

В `backend/cmd/api/main.go`:
```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "https://<username>.github.io,https://<your-domain>.com",
    AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
    AllowHeaders: "Origin,Content-Type,Accept,Authorization",
}))
```

## Переменные окружения

### Frontend (GitHub Secrets)

- `NEXT_PUBLIC_API_URL` - URL backend API

### Backend (на хостинге)

- `DATABASE_URL` - строка подключения к PostgreSQL
- `JWT_SECRET` - секретный ключ для JWT (сгенерируйте случайный)
- `TELEGRAM_BOT_TOKEN` - токен Telegram бота
- `PORT` - порт (обычно 8080)

## Проверка деплоя

1. **Frontend**: Откройте `https://<username>.github.io/<repo>/`
2. **Backend**: Проверьте `https://<your-api-url>/health`
3. **API**: Проверьте `https://<your-api-url>/recipes`

## Troubleshooting

### Frontend не загружается

- Проверьте, что `basePath` настроен правильно
- Убедитесь, что workflow завершился успешно
- Проверьте консоль браузера на ошибки

### Backend не отвечает

- Проверьте логи на хостинге
- Убедитесь, что переменные окружения установлены
- Проверьте подключение к базе данных

### CORS ошибки

- Добавьте домен frontend в `AllowOrigins` в backend
- Проверьте, что заголовки настроены правильно

## Альтернативные варианты деплоя

### Vercel (для Frontend)

Vercel лучше подходит для Next.js, чем GitHub Pages:

1. Подключите репозиторий к Vercel
2. Настройте переменные окружения
3. Автоматический деплой при каждом push

### Docker Compose на VPS

Для полного контроля можно использовать VPS:

1. Арендуйте VPS (DigitalOcean, Linode, etc.)
2. Установите Docker и Docker Compose
3. Склонируйте репозиторий
4. Запустите `docker-compose up -d`

## Мониторинг

После деплоя настройте:

1. **Логирование**: Sentry, LogRocket
2. **Мониторинг**: UptimeRobot, Pingdom
3. **Аналитика**: Google Analytics, Plausible

## Безопасность

1. ✅ Используйте сильные секретные ключи
2. ✅ Включите HTTPS везде
3. ✅ Настройте rate limiting на backend
4. ✅ Регулярно обновляйте зависимости
5. ✅ Используйте переменные окружения для секретов


