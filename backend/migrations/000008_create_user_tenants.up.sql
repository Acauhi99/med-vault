-- User-Tenant relationship
-- Bounded Context: Identity & Access
-- A user can belong to multiple tenants with a role per tenant.

CREATE TABLE user_tenants (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id   UUID         NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    role        VARCHAR(20)  NOT NULL
                CHECK (role IN ('patient', 'doctor', 'administrator')),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),

    UNIQUE (user_id, tenant_id)
);

CREATE INDEX idx_user_tenants_user_id ON user_tenants (user_id);
CREATE INDEX idx_user_tenants_tenant_id ON user_tenants (tenant_id);
CREATE INDEX idx_user_tenants_role ON user_tenants (tenant_id, role);

-- Enable Row-Level Security
ALTER TABLE user_tenants ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_policy ON user_tenants
    USING (tenant_id = current_setting('app.current_tenant_id')::UUID);
