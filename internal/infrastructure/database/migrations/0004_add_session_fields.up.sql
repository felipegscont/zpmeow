-- Add additional fields needed for session management based on backup analysis
-- These fields are based on the original implementation from the backup files

-- Add connected field to track connection state
ALTER TABLE sessions ADD COLUMN connected BOOLEAN DEFAULT FALSE;

-- Add apikey field for session authentication (replacing user/admin tokens)
ALTER TABLE sessions ADD COLUMN apikey TEXT DEFAULT '';

-- Update existing sessions to have a generated apikey
UPDATE sessions SET apikey = gen_random_uuid()::text WHERE apikey = '';

-- Make apikey NOT NULL after setting values
ALTER TABLE sessions ALTER COLUMN apikey SET NOT NULL;

-- Add unique constraint on apikey to ensure uniqueness
CREATE UNIQUE INDEX idx_sessions_apikey ON sessions (apikey) WHERE apikey != '';
