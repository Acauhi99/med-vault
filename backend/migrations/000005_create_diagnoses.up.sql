-- Diagnosis value object (embedded in Case aggregate)
-- Bounded Context: Clinical

CREATE TABLE diagnoses (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id     UUID         NOT NULL REFERENCES cases(id) ON DELETE CASCADE,
    tenant_id   UUID         NOT NULL REFERENCES tenants(id),
    doctor_id   UUID         NOT NULL REFERENCES users(id),
    notes       TEXT         NOT NULL,
    written_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),

    UNIQUE (case_id) -- one diagnosis per case
);

CREATE INDEX idx_diagnoses_case_id ON diagnoses (case_id);
CREATE INDEX idx_diagnoses_tenant_id ON diagnoses (tenant_id);
CREATE INDEX idx_diagnoses_doctor_id ON diagnoses (tenant_id, doctor_id);

-- Enable Row-Level Security
ALTER TABLE diagnoses ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_policy ON diagnoses
    USING (tenant_id = current_setting('app.current_tenant_id')::UUID);
