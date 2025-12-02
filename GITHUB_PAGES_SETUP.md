# Быстрая настройка GitHub Pages

## Шаги для деплоя Frontend на GitHub Pages

### 1. Подготовка репозитория

```bash
# Если репозиторий еще не создан
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/<username>/<repository-name>.git
git push -u origin main
```

### 2. Настройка GitHub Pages

1. Откройте репозиторий на GitHub
2. Перейдите в **Settings** → **Pages**
3. В разделе **Source** выберите **GitHub Actions**
4. Сохраните изменения

### 3. Настройка переменных окружения

1. Перейдите в **Settings** → **Secrets and variables** → **Actions**
2. Нажмите **New repository secret**
3. Добавьте:
   - **Name**: `NEXT_PUBLIC_API_URL`
   - **Value**: URL вашего backend API (например, `https://myplate-api.railway.app`)

### 4. Настройка basePath (если нужно)

Если ваш репозиторий называется не `MyPlate`, обновите `frontend/next.config.js`:

```javascript
basePath: process.env.NODE_ENV === 'production' ? '/<repository-name>' : '',
assetPrefix: process.env.NODE_ENV === 'production' ? '/<repository-name>' : '',
```

### 5. Запуск деплоя

1. Сделайте push в ветку `main`:
   ```bash
   git add .
   git commit -m "Setup GitHub Pages"
   git push
   ```

2. Или запустите вручную:
   - Перейдите в **Actions**
   - Выберите workflow **Deploy to GitHub Pages**
   - Нажмите **Run workflow**

### 6. Проверка

После завершения workflow:
- Frontend будет доступен по адресу: `https://<username>.github.io/<repository-name>/`
- Проверьте статус в **Actions** → **Deploy to GitHub Pages**

## Важно!

⚠️ **GitHub Pages поддерживает только статические сайты**

- ✅ Frontend (Next.js) - можно разместить
- ❌ Backend (Go API) - нужен отдельный хостинг (Railway, Render, Fly.io)
- ❌ База данных (PostgreSQL) - нужен отдельный хостинг

## Деплой Backend

См. [DEPLOYMENT.md](./DEPLOYMENT.md) для инструкций по деплою backend.

Рекомендуемые сервисы:
- **Railway** - самый простой, автоматический деплой из GitHub
- **Render** - бесплатный tier, простой деплой
- **Fly.io** - хороший для Go приложений

## Troubleshooting

### Workflow не запускается

- Убедитесь, что файл `.github/workflows/deploy-pages.yml` существует
- Проверьте, что он в ветке `main`

### Build падает

- Проверьте логи в **Actions**
- Убедитесь, что все зависимости установлены
- Проверьте, что `NEXT_PUBLIC_API_URL` установлен

### Страница не загружается

- Проверьте, что `basePath` настроен правильно
- Убедитесь, что workflow завершился успешно
- Проверьте консоль браузера на ошибки

### API запросы не работают

- Убедитесь, что `NEXT_PUBLIC_API_URL` указывает на правильный URL
- Проверьте CORS настройки на backend
- Убедитесь, что backend доступен и работает


