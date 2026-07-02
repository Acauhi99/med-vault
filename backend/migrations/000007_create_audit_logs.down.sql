DROP TRIGGER IF EXISTS trg_audit_logs_immutable ON audit_logs;
DROP FUNCTION IF EXISTS prevent_audit_log_mutation();
DROP POLICY IF EXISTS tenant_isolation_policy ON audit_logs;
DROP TABLE IF EXISTS audit_logs;
