-- Add webhook fields to sessions table
ALTER TABLE sessions ADD COLUMN webhook_url TEXT DEFAULT '';
ALTER TABLE sessions ADD COLUMN webhook_events TEXT DEFAULT '';
