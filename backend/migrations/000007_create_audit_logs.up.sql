-- AuditLog aggregate root
-- Bounded Context: Audit

CREATE TABLE audit_logs (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID         NOT NULL REFERENCES tenants(id),
    user_id       UUID         NOT NULL REFERENCES users(id),
    action        VARCHAR(100) NOT NULL,
    resource_type VARCHAR(100) NOT NULL,
    resource_id   UUID         NOT NULL,
    ip_address    INET,
    metadata      JSONB,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_logs_tenant_id ON audit_logs (tenant_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs (tenant_id, user_id);
CREATE INDEX idx_audit_logs_resource ON audit_logs (tenant_id, resource_type, resource_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs (tenant_id, created_at DESC);

-- Enable Row-Level Security
ALTER TABLE audit_logs ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_policy ON audit_logs
    USING (tenant_id = current_setting('app.current_tenant_id')::UUID);

-- Audit logs are immutable: prevent UPDATE and DELETE
CREATE OR REPLACE FUNCTION prevent_audit_log_mutation()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'audit logs are immutable';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_audit_logs_immutable
    BEFORE UPDATE OR DELETE ON audit_logs
    FOR EACH ROW
    EXECUTE FUNCTION prevent_audit_log_mutation();
