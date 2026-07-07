-- Remove user_agent column from audit_logs

ALTER TABLE audit_logs DROP COLUMN IF EXISTS user_agent;
