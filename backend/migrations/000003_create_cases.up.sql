-- Case aggregate root
-- Bounded Context: Clinical

CREATE TABLE cases (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID         NOT NULL REFERENCES tenants(id),
    patient_id  UUID         NOT NULL REFERENCES users(id),
    doctor_id   UUID         REFERENCES users(id),
    status      VARCHAR(20)  NOT NULL DEFAULT 'open'
                CHECK (status IN ('open', 'assigned', 'diagnosed', 'closed')),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    closed_at   TIMESTAMPTZ
);

CREATE INDEX idx_cases_tenant_id ON cases (tenant_id);
CREATE INDEX idx_cases_patient_id ON cases (tenant_id, patient_id);
CREATE INDEX idx_cases_doctor_id ON cases (tenant_id, doctor_id);
CREATE INDEX idx_cases_status ON cases (tenant_id, status);

-- Enable Row-Level Security
ALTER TABLE cases ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_policy ON cases
    USING (tenant_id = current_setting('app.current_tenant_id')::UUID);
