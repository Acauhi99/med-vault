-- Image aggregate root
-- Bounded Context: Imaging

CREATE TABLE images (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID         NOT NULL REFERENCES tenants(id),
    case_id      UUID         NOT NULL REFERENCES cases(id),
    patient_id   UUID         NOT NULL REFERENCES users(id),
    file_name    VARCHAR(255) NOT NULL,
    content_type VARCHAR(100) NOT NULL
                 CHECK (content_type IN ('image/jpeg', 'image/png', 'image/dicom')),
    s3_key       VARCHAR(512) NOT NULL,
    uploaded_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_images_tenant_id ON images (tenant_id);
CREATE INDEX idx_images_case_id ON images (tenant_id, case_id);
CREATE INDEX idx_images_patient_id ON images (tenant_id, patient_id);

-- Enable Row-Level Security
ALTER TABLE images ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_policy ON images
    USING (tenant_id = current_setting('app.current_tenant_id')::UUID);
