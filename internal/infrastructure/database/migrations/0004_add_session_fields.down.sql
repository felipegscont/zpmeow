-- Remove additional session fields
DROP INDEX IF EXISTS idx_sessions_apikey;
ALTER TABLE sessions DROP COLUMN IF EXISTS connected;
ALTER TABLE sessions DROP COLUMN IF EXISTS apikey;
