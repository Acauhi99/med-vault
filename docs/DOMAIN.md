# Domain Model

This document defines the domain model for MedVault using Domain-Driven Design (DDD) tactical and strategic patterns with CQRS.

---

## Ubiquitous Language

| Term | Definition |
|------|------------|
| **Tenant** | A healthcare organization. Top-level boundary for all data isolation. |
| **User** | An authenticated person within a tenant. Has exactly one role. |
| **Patient** | A user who submits symptoms and medical images. |
| **Doctor** | A user who reviews cases and writes diagnoses. |
| **Administrator** | A user who manages tenant resources and inspects audit information. |
| **Case** | A medical consultation request initiated by a patient. Contains symptoms and images. |
| **Symptom** | A patient-reported health concern attached to a case. |
| **Image** | A medical image (X-ray, scan, etc.) uploaded by a patient for a case. |
| **Diagnosis** | A doctor's written assessment of a case. |
| **Audit Log** | An immutable record of a state-changing operation. |

---

## Strategic Design — Bounded Contexts

```
┌─────────────────────────────────────────────────────────────────────┐
│                          MedVault Platform                          │
├─────────────────┬─────────────────┬─────────────────┬───────────────┤
│    Identity     │    Clinical     │     Imaging     │    Audit      │
│    & Access     │                 │                 │               │
├─────────────────┼─────────────────┼─────────────────┼───────────────┤
│  Tenants        │  Cases          │  Images         │  Audit Logs   │
│  Users          │  Symptoms       │                 │               │
│  Authentication │  Diagnoses      │                 │               │
│  Authorization  │                 │                 │               │
└─────────────────┴─────────────────┴─────────────────┴───────────────┘
```

### Context Map

```
Identity & Access ──provides──▶ Clinical (user identity, tenant context)
Identity & Access ──provides──▶ Imaging (user identity, tenant context)
Identity & Access ──provides──▶ Audit (user identity, tenant context)
Clinical ──publishes──▶ Audit (CaseCreated, CaseUpdated, DiagnosisWritten)
Imaging ──publishes──▶ Audit (ImageUploaded, ImageAccessed)
```

**Relationship types:**
- Identity & Access is a **upstream** provider to all other contexts
- Clinical and Imaging are **downstream** consumers of Identity & Access
- Clinical and Imaging are **upstream** publishers to Audit

---

## Bounded Context: Identity & Access

### Aggregates

#### Tenant (Aggregate Root)

```
Tenant
├── TenantID (Value Object)
├── Name (Value Object)
├── Status (Value Object: Active | Suspended)
└── CreatedAt (Value Object)
```

**Invariants:**
- Tenant name must be non-empty and unique within the system
- Status transitions: Active → Suspended (no reactivation in PoC)

#### User (Aggregate Root)

```
User
├── UserID (Value Object)
├── TenantID (Value Object)
├── Email (Value Object)
├── PasswordHash (Value Object)
├── Role (Value Object: Patient | Doctor | Administrator)
├── Status (Value Object: Active | Inactive)
├── CreatedAt (Value Object)
└── UpdatedAt (Value Object)
```

**Invariants:**
- Email must be unique within a tenant
- Password must meet minimum complexity requirements
- Role cannot be changed after creation
- Status transitions: Active ↔ Inactive

### Commands

| Command | Description |
|---------|-------------|
| `CreateTenant` | Register a new healthcare organization |
| `SuspendTenant` | Disable a tenant (all users lose access) |
| `RegisterUser` | Create a new user within a tenant |
| `AuthenticateUser` | Validate credentials, issue JWT |
| `DeactivateUser` | Disable a user account |

### Queries

| Query | Description |
|-------|-------------|
| `GetUserByID` | Retrieve user by ID (tenant-scoped) |
| `GetUserByEmail` | Retrieve user by email (tenant-scoped) |
| `ListUsersByTenant` | List all users for a tenant |

### Domain Events

| Event | Trigger |
|-------|---------|
| `TenantCreated` | New tenant registered |
| `TenantSuspended` | Tenant status changed to Suspended |
| `UserRegistered` | New user created |
| `UserDeactivated` | User status changed to Inactive |

---

## Bounded Context: Clinical

### Aggregates

#### Case (Aggregate Root)

```
Case
├── CaseID (Value Object)
├── TenantID (Value Object)
├── PatientID (Value Object) → references User
├── DoctorID (Value Object | nil) → references User
├── Status (Value Object: Open | Assigned | Diagnosed | Closed)
├── Symptoms (Entity Collection)
├── Diagnosis (Value Object | nil)
├── CreatedAt (Value Object)
├── UpdatedAt (Value Object)
└── ClosedAt (Value Object | nil)
```

**Invariants:**
- Case must belong to exactly one patient
- Status transitions: Open → Assigned → Diagnosed → Closed
- Doctor assignment only allowed when status is Open
- Diagnosis can only be written when status is Assigned
- Case can only be closed after diagnosis is written

#### Symptom (Entity)

```
Symptom
├── SymptomID (Value Object)
├── Description (Value Object)
├── Severity (Value Object: Low | Medium | High | Critical)
└── ReportedAt (Value Object)
```

**Invariants:**
- Description must be non-empty
- Severity must be one of the defined values

#### Diagnosis (Value Object)

```
Diagnosis
├── Notes (Value Object)
├── DoctorID (Value Object)
├── WrittenAt (Value Object)
```

**Invariants:**
- Notes must be non-empty
- Must be written by the assigned doctor

### Commands

| Command | Description |
|---------|-------------|
| `CreateCase` | Patient creates a new medical case |
| `AddSymptom` | Patient adds a symptom to an existing case |
| `AssignDoctor` | Administrator assigns a doctor to a case |
| `WriteDiagnosis` | Doctor writes a diagnosis for an assigned case |
| `CloseCase` | Administrator closes a diagnosed case |

### Queries

| Query | Description |
|-------|-------------|
| `GetCaseByID` | Retrieve a case with all symptoms and diagnosis |
| `ListCasesByPatient` | List all cases for a patient (tenant-scoped) |
| `ListCasesByDoctor` | List all cases assigned to a doctor (tenant-scoped) |
| `ListOpenCases` | List all cases with status Open (admin view) |
| `ListAllCases` | List all cases for a tenant (admin view) |

### Domain Events

| Event | Trigger |
|-------|---------|
| `CaseCreated` | New case created by patient |
| `SymptomAdded` | Symptom added to a case |
| `DoctorAssigned` | Doctor assigned to a case |
| `DiagnosisWritten` | Diagnosis written for a case |
| `CaseClosed` | Case closed by administrator |

---

## Bounded Context: Imaging

### Aggregates

#### Image (Aggregate Root)

```
Image
├── ImageID (Value Object)
├── TenantID (Value Object)
├── CaseID (Value Object) → references Case
├── PatientID (Value Object) → references User
├── FileName (Value Object)
├── ContentType (Value Object)
├── S3Key (Value Object)
├── UploadedAt (Value Object)
```

**Invariants:**
- Image must belong to exactly one case
- Image must be uploaded by the case's patient
- ContentType must be an allowed medical image format

### Commands

| Command | Description |
|---------|-------------|
| `RequestUploadURL` | Generate pre-signed S3 URL for image upload |
| `ConfirmUpload` | Record image metadata after successful upload |
| `DeleteImage` | Remove image from case (soft delete) |

### Queries

| Query | Description |
|-------|-------------|
| `GetImageByID` | Retrieve image metadata |
| `ListImagesByCase` | List all images for a case |
| `GetImageDownloadURL` | Generate pre-signed URL for image download |

### Domain Events

| Event | Trigger |
|-------|---------|
| `ImageUploaded` | Image metadata recorded after upload |
| `ImageDeleted` | Image soft-deleted |

---

## Bounded Context: Audit

### Aggregates

#### AuditLog (Aggregate Root)

```
AuditLog
├── AuditLogID (Value Object)
├── TenantID (Value Object)
├── UserID (Value Object)
├── Action (Value Object)
├── ResourceType (Value Object)
├── ResourceID (Value Object)
├── Timestamp (Value Object)
├── IPAddress (Value Object)
└── Metadata (Value Object: optional JSON)
```

**Invariants:**
- Audit logs are immutable once created
- Every state-changing operation must produce an audit log entry
- Audit logs are tenant-scoped

### Commands

| Command | Description |
|---------|-------------|
| `RecordAuditLog` | Create an immutable audit log entry |

### Queries

| Query | Description |
|-------|-------------|
| `GetAuditLogByID` | Retrieve a specific audit log entry |
| `ListAuditLogsByTenant` | List all audit logs for a tenant |
| `ListAuditLogsByUser` | List all audit logs for a specific user |
| `ListAuditLogsByResource` | List all audit logs for a specific resource |

### Domain Events

| Event | Trigger |
|-------|---------|
| `AuditLogRecorded` | New audit log entry created |

---

## Value Objects

All value objects are immutable and implement equality based on their attributes.

| Value Object | Used In | Rules |
|-------------|---------|-------|
| `TenantID` | Everywhere | UUID format |
| `UserID` | Everywhere | UUID format |
| `Email` | User | Valid email format, unique within tenant |
| `PasswordHash` | User | bcrypt hash |
| `CaseID` | Clinical, Imaging | UUID format |
| `ImageID` | Imaging | UUID format |
| `SymptomID` | Clinical | UUID format |
| `AuditLogID` | Audit | UUID format |
| `Role` | Identity & Access | Enum: Patient, Doctor, Administrator |
| `Status` | Various | Context-specific enums |
| `Severity` | Clinical | Enum: Low, Medium, High, Critical |
| `S3Key` | Imaging | String, format: `{tenant_id}/{case_id}/{image_id}/{filename}` |
| `ContentType` | Imaging | Allowlist: image/jpeg, image/png, image/dicom |

---

## Aggregate Relationships

```
Tenant (1) ──── (N) User
Tenant (1) ──── (N) Case
Tenant (1) ──── (N) Image
Tenant (1) ──── (N) AuditLog

User (Patient) (1) ──── (N) Case
User (Doctor) (1) ──── (N) Case

Case (1) ──── (N) Symptom
Case (1) ──── (0..1) Diagnosis
Case (1) ──── (N) Image
```

**Aggregate boundaries:**
- `Case` aggregate contains `Symptom` entities and `Diagnosis` value object
- `Image` is a separate aggregate (references `Case` by ID, does not contain it)
- `Tenant` and `User` are separate aggregates
- `AuditLog` is a separate aggregate (fire-and-forget, no consistency boundary with other aggregates)

---

## CQRS Mapping

### Write Side (Commands)

| Command | Aggregate | Handler |
|---------|-----------|---------|
| `CreateTenant` | Tenant | TenantCommandHandler |
| `RegisterUser` | User | UserCommandHandler |
| `CreateCase` | Case | CaseCommandHandler |
| `AddSymptom` | Case | CaseCommandHandler |
| `AssignDoctor` | Case | CaseCommandHandler |
| `WriteDiagnosis` | Case | CaseCommandHandler |
| `RequestUploadURL` | Image | ImageCommandHandler |
| `ConfirmUpload` | Image | ImageCommandHandler |

### Read Side (Queries)

| Query | Read Model | Handler |
|-------|------------|---------|
| `GetUserByID` | UserReadModel | UserQueryHandler |
| `ListCasesByPatient` | CaseReadModel | CaseQueryHandler |
| `ListCasesByDoctor` | CaseReadModel | CaseQueryHandler |
| `ListImagesByCase` | ImageReadModel | ImageQueryHandler |
| `ListAuditLogsByTenant` | AuditLogReadModel | AuditQueryHandler |

### Event Handlers (Projections)

| Event | Projection |
|-------|------------|
| `UserRegistered` | Update UserReadModel |
| `CaseCreated` | Update CaseReadModel |
| `SymptomAdded` | Update CaseReadModel |
| `DoctorAssigned` | Update CaseReadModel |
| `DiagnosisWritten` | Update CaseReadModel |
| `CaseClosed` | Update CaseReadModel |
| `ImageUploaded` | Update ImageReadModel |
| `AuditLogRecorded` | Update AuditLogReadModel |

---

## Multi-Tenancy in the Domain

Every aggregate root includes `TenantID` as a first-class attribute.

**Enforcement points:**
1. **Command handlers** validate `TenantID` from JWT context
2. **Repository interfaces** include `TenantID` parameter
3. **Database queries** include `WHERE tenant_id = $1`
4. **Read models** are partitioned by `TenantID`

**No cross-tenant access** is permitted at the domain level.
