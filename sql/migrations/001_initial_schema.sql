-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    username TEXT,
    first_name TEXT,
    last_name TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User goals table
CREATE TABLE user_goals (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    daily_calories INT NOT NULL,
    target_proteins NUMERIC(10,2),
    target_fats NUMERIC(10,2),
    target_carbs NUMERIC(10,2),
    protein_ratio NUMERIC(5,2),
    fat_ratio NUMERIC(5,2),
    carb_ratio NUMERIC(5,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id)
);

-- Recipes table
CREATE TABLE recipes (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    calories INT NOT NULL,
    proteins NUMERIC(10,2),
    fats NUMERIC(10,2),
    carbs NUMERIC(10,2),
    price NUMERIC(10,2) DEFAULT 0,
    cooking_time INT NOT NULL, -- in minutes
    servings INT DEFAULT 1,
    meal_type TEXT CHECK (meal_type IN ('breakfast', 'lunch', 'dinner', 'snack')),
    diet_type TEXT[], -- e.g., ['vegetarian', 'vegan', 'gluten-free']
    allergens TEXT[], -- e.g., ['nuts', 'dairy', 'eggs']
    ingredients JSONB NOT NULL, -- Array of {name, quantity, unit}
    instructions TEXT[],
    image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Pantry items table
CREATE TABLE pantry_items (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    quantity NUMERIC(10,2) NOT NULL,
    unit TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Menus table
CREATE TABLE menus (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    total_calories INT,
    total_price NUMERIC(10,2),
    total_time INT, -- in minutes
    meals JSONB NOT NULL, -- Array of {recipe_id, meal_type, calories, price, time}
    ingredients_used JSONB, -- Array of ingredients used from pantry
    missing_ingredients JSONB, -- Array of missing ingredients
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, date)
);

-- Shopping lists table
CREATE TABLE shopping_lists (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    menu_id INT REFERENCES menus(id) ON DELETE CASCADE,
    items JSONB NOT NULL, -- Array of {name, quantity, unit, reason}
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_users_telegram_id ON users(telegram_id);
CREATE INDEX idx_user_goals_user_id ON user_goals(user_id);
CREATE INDEX idx_pantry_items_user_id ON pantry_items(user_id);
CREATE INDEX idx_menus_user_id ON menus(user_id);
CREATE INDEX idx_menus_date ON menus(date);
CREATE INDEX idx_shopping_lists_user_id ON shopping_lists(user_id);
CREATE INDEX idx_shopping_lists_menu_id ON shopping_lists(menu_id);
CREATE INDEX idx_recipes_meal_type ON recipes(meal_type);
CREATE INDEX idx_recipes_calories ON recipes(calories);
CREATE INDEX idx_recipes_price ON recipes(price);
CREATE INDEX idx_recipes_cooking_time ON recipes(cooking_time);

-- GIN indexes for JSONB and array columns
CREATE INDEX idx_recipes_ingredients ON recipes USING GIN (ingredients);
CREATE INDEX idx_recipes_diet_type ON recipes USING GIN (diet_type);
CREATE INDEX idx_recipes_allergens ON recipes USING GIN (allergens);
CREATE INDEX idx_menus_meals ON menus USING GIN (meals);


