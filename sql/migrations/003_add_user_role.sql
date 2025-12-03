-- Миграция: добавление роли пользователя
-- Добавляем поле role в таблицу users

ALTER TABLE users
ADD COLUMN role TEXT NOT NULL DEFAULT 'user';

-- Создаем индекс для быстрого поиска по роли
CREATE INDEX idx_users_role ON users(role);

-- Комментарии
COMMENT ON COLUMN users.role IS 'Роль пользователя: user или admin';

