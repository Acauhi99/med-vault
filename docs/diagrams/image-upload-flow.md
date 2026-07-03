# Image Upload Flow

## Pre-Signed URL Upload (Client → S3 Direct)

```mermaid
sequenceDiagram
    actor Patient
    participant Frontend
    participant Backend
    participant S3
    participant DB

    rect rgb(240, 248, 255)
        Note over Patient, DB: Step 1 — Request Upload URL
        Patient->>Frontend: Select image file to upload
        Frontend->>Frontend: Validate file type (JPEG, PNG, DICOM)
        Frontend->>Frontend: Validate file size (≤ 50MB)
        Frontend->>Backend: POST /api/v1/cases/{id}/images/upload-url {filename, content_type}
        Backend->>Backend: Validate JWT + role = patient
        Backend->>Backend: Validate case belongs to patient
        Backend->>Backend: Validate case status = Open
        Backend->>S3: Generate pre-signed PUT URL
        Note right of S3: Path: /{tenant_id}/{case_id}/{image_id}/{filename}
        Note right of S3: Expiry: 15 minutes
        Note right of S3: Content-Type restricted
        S3-->>Backend: Pre-signed URL
        Backend-->>Frontend: 200 OK {upload_url, image_id}
    end

    rect rgb(240, 255, 240)
        Note over Patient, S3: Step 2 — Upload Directly to S3
        Patient->>S3: PUT image via pre-signed URL
        Note right of S3: TLS encrypted
        Note right of S3: No backend bandwidth used
        S3-->>Patient: 200 OK (upload complete)
    end

    rect rgb(255, 248, 240)
        Note over Patient, DB: Step 3 — Confirm Upload
        Patient->>Frontend: Upload complete
        Frontend->>Backend: POST /api/v1/cases/{id}/images {image_id, filename, content_type, size}
        Backend->>Backend: Validate metadata
        Backend->>DB: BEGIN TRANSACTION
        Backend->>DB: INSERT image (tenant_id, case_id, patient_id, s3_key, metadata)
        Backend->>DB: INSERT domain_outbox (ImageUploaded event)
        Backend->>DB: INSERT audit_log (ImageUploaded)
        Backend->>DB: COMMIT
        Backend-->>Frontend: 201 Created {image}
        Frontend-->>Patient: Image uploaded successfully
    end
```

## Image Download (Doctor Viewing)

```mermaid
sequenceDiagram
    actor Doctor
    participant Frontend
    participant Backend
    participant S3
    participant DB

    Doctor->>Frontend: Open case → View images
    Frontend->>Backend: GET /api/v1/cases/{id}/images
    Backend->>Backend: Validate JWT + role = doctor
    Backend->>Backend: Validate case is assigned to this doctor
    Backend->>DB: SELECT images WHERE case_id = $1 AND tenant_id = $2
    DB-->>Backend: Image metadata
    Backend-->>Frontend: 200 OK {images[]}

    loop For each image
        Frontend->>Backend: GET /api/v1/images/{id}/download-url
        Backend->>Backend: Validate access (patient, assigned doctor, or admin)
        Backend->>S3: Generate pre-signed GET URL
        Note right of S3: Expiry: 15 minutes
        Note right of S3: Scoped to image path
        S3-->>Backend: Pre-signed URL
        Backend->>DB: INSERT audit_log (ImageAccessed)
        Backend-->>Frontend: 200 OK {download_url}
        Frontend->>S3: GET image via pre-signed URL
        S3-->>Frontend: Image data
    end

    Frontend-->>Doctor: Display images
```

## Security Constraints

| Constraint | Implementation |
|------------|----------------|
| File types | Allowlist: image/jpeg, image/png, image/dicom |
| Max file size | 50MB |
| Upload URL expiry | 15 minutes |
| Download URL expiry | 15 minutes |
| Storage path | `/{tenant_id}/{case_id}/{image_id}/{filename}` |
| Encryption at rest | AES-256 (SSE-S3) |
| Encryption in transit | TLS 1.2+ |
| Access control | Patient (own), Doctor (assigned), Admin (tenant) |
| Audit | Every upload and download logged |
