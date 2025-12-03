-- Add menu_type column to menus table to support both daily and weekly menus
ALTER TABLE menus
ADD COLUMN menu_type TEXT NOT NULL DEFAULT 'daily' CHECK (menu_type IN ('daily', 'weekly'));

COMMENT ON COLUMN menus.menu_type IS 'Type of menu: "daily" for single day menu, "weekly" for 7-day menu';

-- Update unique constraint to allow multiple weekly menus per user
-- (but still only one daily menu per user per date)
ALTER TABLE menus DROP CONSTRAINT IF EXISTS menus_user_id_date_key;

-- Create new unique constraint for daily menus only
CREATE UNIQUE INDEX idx_menus_user_id_date_daily 
ON menus(user_id, date) 
WHERE menu_type = 'daily';

-- Index for menu_type for faster queries
CREATE INDEX idx_menus_menu_type ON menus(menu_type);
CREATE INDEX idx_menus_user_id_menu_type ON menus(user_id, menu_type);

