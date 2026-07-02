-- Symptom entity (child of Case aggregate)
-- Bounded Context: Clinical

CREATE TABLE symptoms (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id     UUID         NOT NULL REFERENCES cases(id) ON DELETE CASCADE,
    tenant_id   UUID         NOT NULL REFERENCES tenants(id),
    description TEXT         NOT NULL,
    severity    VARCHAR(20)  NOT NULL
                CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    reported_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_symptoms_case_id ON symptoms (case_id);
CREATE INDEX idx_symptoms_tenant_id ON symptoms (tenant_id);

-- Enable Row-Level Security
ALTER TABLE symptoms ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_policy ON symptoms
    USING (tenant_id = current_setting('app.current_tenant_id')::UUID);
