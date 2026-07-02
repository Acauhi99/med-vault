-- Tenant aggregate root
-- Bounded Context: Identity & Access

CREATE TABLE tenants (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    status      VARCHAR(20)  NOT NULL DEFAULT 'active'
                CHECK (status IN ('active', 'suspended')),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_tenants_name ON tenants (name);
CREATE INDEX idx_tenants_status ON tenants (status);

-- Enable Row-Level Security
ALTER TABLE tenants ENABLE ROW LEVEL SECURITY;

-- RLS policy: tenant can only see itself
CREATE POLICY tenant_isolation_policy ON tenants
    USING (id = current_setting('app.current_tenant_id')::UUID);
