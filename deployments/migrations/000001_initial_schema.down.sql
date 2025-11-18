-- Rollback initial schema migration

-- Drop tables in reverse order
DROP TABLE IF EXISTS audit_log CASCADE;
DROP TABLE IF EXISTS api_keys CASCADE;
DROP TABLE IF EXISTS feature_flags CASCADE;
DROP TABLE IF EXISTS locations CASCADE;
DROP TABLE IF EXISTS query_log CASCADE;
DROP TABLE IF EXISTS panchangam_cache CASCADE;
DROP TABLE IF EXISTS user_preferences CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Drop functions
DROP FUNCTION IF EXISTS clean_expired_cache();
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop types
DROP TYPE IF EXISTS regional_variation;
DROP TYPE IF EXISTS calculation_system;

-- Drop extensions (be careful with this in production)
-- DROP EXTENSION IF EXISTS "pg_trgm";
-- DROP EXTENSION IF EXISTS "uuid-ossp";
