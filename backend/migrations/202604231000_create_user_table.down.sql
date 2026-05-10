DROP TRIGGER IF EXISTS update_app_user_updated_at ON app_user;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS app_user_username_idx;
DROP TABLE IF EXISTS app_user;