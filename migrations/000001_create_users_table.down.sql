-- Drop trigger
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_is_active;

-- Drop table
-- Drop constraints first
ALTER TABLE users DROP CONSTRAINT IF EXISTS uni_users_email;
ALTER TABLE users DROP CONSTRAINT IF EXISTS uni_users_username;

-- Drop table
DROP TABLE IF EXISTS users;
