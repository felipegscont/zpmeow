-- Ensure device_jid column exists (no-op if already exists)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='sessions' AND column_name='device_jid') THEN
        ALTER TABLE sessions ADD COLUMN device_jid TEXT;
    END IF;
END $$;
