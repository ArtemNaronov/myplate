# Настройка переменных окружения

## Файл .env

Файл `.env` содержит все переменные окружения для проекта. Он уже создан с настройками по умолчанию.

### Текущие настройки:

**База данных PostgreSQL:**
- Пользователь: `myplate`
- Пароль: `myplate`
- База данных: `myplate`
- Порт: `5433` (внешний), `5432` (внутри контейнера)

**Backend API:**
- Порт: `8080`
- JWT Secret: `your-secret-key-change-in-production` ⚠️ **Измените для production!**

**Frontend:**
- API URL: `http://localhost:8080`

**PgAdmin:**
- Email: `admin@myplate.com`
- Пароль: `admin`

## Изменение настроек

### Для разработки:

1. Откройте файл `.env`
2. Измените нужные переменные
3. Перезапустите контейнеры:
   ```bash
   docker-compose down
   docker-compose up -d
   ```

### Для production:

**ВАЖНО:** Обязательно измените следующие значения:

1. **POSTGRES_PASSWORD** - используйте сильный пароль
2. **JWT_SECRET** - сгенерируйте случайную строку (минимум 32 символа)
3. **PGADMIN_DEFAULT_PASSWORD** - используйте сильный пароль

Пример генерации JWT_SECRET:
```bash
# Linux/Mac
openssl rand -base64 32

# Или используйте онлайн генератор
```

## Подключение к базе данных

### Из приложения:
```
postgres://myplate:myplate@postgres:5432/myplate?sslmode=disable
```

### Извне (с хоста):
```
postgres://myplate:myplate@localhost:5433/myplate?sslmode=disable
```

### Через PgAdmin:
- URL: http://localhost:5050
- Email: `admin@myplate.com`
- Пароль: `admin`

Для подключения к базе данных в PgAdmin:
- Host: `postgres` (внутри Docker) или `localhost` (снаружи)
- Port: `5432` (внутри Docker) или `5433` (снаружи)
- Database: `myplate`
- Username: `myplate`
- Password: `myplate`

## Безопасность

⚠️ **НИКОГДА не коммитьте файл `.env` в Git!**

Файл `.env` уже добавлен в `.gitignore`. Используйте `.env.example` как шаблон для других разработчиков.

## Проверка переменных

Проверить, что переменные загружены правильно:

```bash
# Внутри контейнера backend
docker exec myplate-api env | grep -E "DATABASE_URL|JWT_SECRET|PORT"

# Внутри контейнера frontend
docker exec myplate-frontend env | grep NEXT_PUBLIC_API_URL
```

