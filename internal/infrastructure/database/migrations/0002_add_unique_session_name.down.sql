-- Remove unique constraint from session names
ALTER TABLE sessions DROP CONSTRAINT IF EXISTS unique_session_name;
