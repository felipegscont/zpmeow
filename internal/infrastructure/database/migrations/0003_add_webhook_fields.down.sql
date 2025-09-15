-- Remove webhook fields from sessions table
ALTER TABLE sessions DROP COLUMN IF EXISTS webhook_url;
ALTER TABLE sessions DROP COLUMN IF EXISTS webhook_events;
