-- Add unique constraint to session names
ALTER TABLE sessions ADD CONSTRAINT unique_session_name UNIQUE (name);
