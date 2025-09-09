-- Revert: rename whatsapp_jid back to device_jid
DO $$ 
BEGIN 
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='sessions' AND column_name='whatsapp_jid') THEN
        ALTER TABLE sessions ADD COLUMN device_jid TEXT;
        UPDATE sessions SET device_jid = whatsapp_jid WHERE whatsapp_jid IS NOT NULL;
        ALTER TABLE sessions DROP COLUMN whatsapp_jid;
    END IF;
END $$;
