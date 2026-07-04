# Domain Model

This document defines the domain model for MedVault using Domain-Driven Design (DDD) tactical and strategic patterns with CQRS.

---

## Ubiquitous Language

| Term | Definition |
|------|------------|
| **Tenant** | A healthcare organization. Top-level boundary for all data isolation. |
| **User** | An authenticated person. Exists independently of tenants. Has roles per tenant via `UserTenant`. |
| **Patient** | A user with role `patient` in a specific tenant. |
| **Doctor** | A user with role `doctor` in a specific tenant. |
| **Administrator** | A user with role `administrator` in a specific tenant. |
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
- Status transitions: Active ↔ Suspended

#### User (Aggregate Root)

```
User
├── UserID (Value Object)
├── Email (Value Object)
├── PasswordHash (Value Object)
├── Status (Value Object: Active | Inactive)
├── CreatedAt (Value Object)
└── UpdatedAt (Value Object)
```

**Invariants:**
- Email must be unique (system-wide, not per tenant)
- Password must meet minimum complexity requirements
- Status transitions: Active ↔ Inactive
- A user has no inherent tenant — tenant membership is defined by `UserTenant`
- Registration creates a user without tenant membership; an admin must add the user to a tenant via `AddUserToTenant`

#### UserTenant (Entity)

```
UserTenant
├── UserTenantID (Value Object)
├── UserID (Value Object) → references User
├── TenantID (Value Object) → references Tenant
├── Role (Value Object: Patient | Doctor | Administrator)
└── CreatedAt (Value Object)
```

**Invariants:**
- A user can belong to multiple tenants
- A user has exactly one role per tenant
- The same (user_id, tenant_id) pair cannot be duplicated
- Role cannot be changed after creation (delete + re-create to change)

### Commands

| Command | Description |
|---------|-------------|
| `CreateTenant` | Register a new healthcare organization |
| `SuspendTenant` | Disable a tenant (all users lose access to that tenant) |
| `ReactivateTenant` | Restore a suspended tenant to active status |
| `RegisterUser` | Create a new user (no tenant assignment) |
| `AuthenticateUser` | Validate email + password, return available tenants |
| `SelectTenant` | User selects a tenant, receives JWT with tenant_id + role |
| `AddUserToTenant` | Associate an existing user with a tenant and assign a role |
| `RemoveUserFromTenant` | Remove a user's access to a tenant |
| `DeactivateUser` | Disable a user account (across all tenants) |

### Queries

| Query | Description |
|-------|-------------|
| `GetUserByID` | Retrieve user by ID (tenant-independent) |
| `GetUserByEmail` | Retrieve user by email (tenant-independent) |
| `ListTenantsByUser` | List all tenants a user belongs to (with roles) |
| `ListUsersByTenant` | List all users for a tenant (with roles) |
| `ListTenantMembers` | List tenant members with user details (joined query) |

### Domain Events

| Event | Trigger |
|-------|---------|
| `TenantCreated` | New tenant registered |
| `TenantSuspended` | Tenant status changed to Suspended |
| `TenantReactivated` | Tenant status changed back to Active |
| `UserRegistered` | New user created (no tenant) |
| `UserAddedToTenant` | User associated with a tenant |
| `UserRemovedFromTenant` | User's access to a tenant revoked |
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
| `ListAllCases` | List all cases for a tenant (admin view, filterable by status, sortable by created_at) |

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
| `GetAuditLogByID` | Not implemented in backend yet |
| `ListAuditLogsByTenant` | List all audit logs for a tenant (filterable by resource_type/resource_id) |
| `ListAuditLogsByUser` | Not implemented in backend yet |
| `ListAuditLogsByResource` | Not implemented in backend yet |

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
| `Email` | User | Valid email format, unique system-wide |
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
Tenant (1) ──── (N) UserTenant
User   (1) ──── (N) UserTenant

Tenant (1) ──── (N) Case
Tenant (1) ──── (N) Image
Tenant (1) ──── (N) AuditLog

User (via UserTenant, role=Patient) (1) ──── (N) Case
User (via UserTenant, role=Doctor)  (1) ──── (N) Case

Case (1) ──── (N) Symptom
Case (1) ──── (0..1) Diagnosis
Case (1) ──── (N) Image
```

**Aggregate boundaries:**
- `User` aggregate is independent — no `TenantID` attribute
- `UserTenant` is a separate entity linking User to Tenant with a role
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
| `SelectTenant` | UserTenant | UserCommandHandler |
| `AddUserToTenant` | UserTenant | UserCommandHandler |
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
| `ListTenantsByUser` | UserTenantReadModel | UserQueryHandler |
| `ListUsersByTenant` | UserTenantReadModel | UserQueryHandler |
| `ListCasesByPatient` | CaseReadModel | CaseQueryHandler |
| `ListCasesByDoctor` | CaseReadModel | CaseQueryHandler |
| `ListImagesByCase` | ImageReadModel | ImageQueryHandler |
| `ListAuditLogsByTenant` | AuditLogReadModel | AuditQueryHandler |

### Event Handlers (Projections)

> **Detailed event flow diagram:** See [diagrams/domain-events-flow.md](diagrams/domain-events-flow.md)

Events are delivered via the **Transactional Outbox** pattern (see [ADR-017](adr/017-transactional-outbox.md)). Events are persisted in the same transaction as the aggregate. A poller reads unpublished events and dispatches to projection handlers. Delivery guarantee: at-least-once (handlers must be idempotent).

| Event | Projection |
|-------|------------|
| `UserRegistered` | Update UserReadModel |
| `UserAddedToTenant` | Update UserTenantReadModel |
| `UserRemovedFromTenant` | Update UserTenantReadModel |
| `CaseCreated` | Update CaseReadModel |
| `SymptomAdded` | Update CaseReadModel |
| `DoctorAssigned` | Update CaseReadModel |
| `DiagnosisWritten` | Update CaseReadModel |
| `CaseClosed` | Update CaseReadModel |
| `ImageUploaded` | Update ImageReadModel |
| `AuditLogRecorded` | Update AuditLogReadModel |

**Event delivery flow:**

```
Command Handler
  → Aggregate.Mutate()
  → BEGIN TRANSACTION
      → Repo.Save(aggregate)
      → Outbox.Save(event)        ← same transaction
  → COMMIT

Outbox Poller (goroutine, every ~1s)
  → SELECT FROM domain_outbox WHERE published = false
  → ProjectionHandler.Handle(event)
  → Mark as published
```

---

## Multi-Tenancy in the Domain

Users exist independently of tenants. A user belongs to one or more tenants via the `UserTenant` entity, with a role per tenant. The `tenant_id` for a request comes from the JWT (selected during the `SelectTenant` command), not from the user record.

**Enforcement points:**
1. **Command handlers** validate `TenantID` from JWT context
2. **Repository interfaces** include `TenantID` parameter
3. **Database queries** include `WHERE tenant_id = $1`
4. **Read models** are partitioned by `TenantID`

**No cross-tenant access** is permitted at the domain level. A doctor who belongs to Tenant A and Tenant B can only access the tenant they selected at login.
