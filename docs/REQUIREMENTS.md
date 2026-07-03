# Requirements

This document defines the functional and non-functional requirements for MedVault.

---

## Functional Requirements

### Patient

| ID | Requirement | Priority |
|----|-------------|----------|
| P-01 | Register a new account (system-wide, no tenant assignment) | Must |
| P-02 | Authenticate (login/logout) | Must |
| P-03 | Submit a new medical case with symptoms | Must |
| P-04 | Add symptoms to an existing case | Must |
| P-05 | Upload medical images to a case | Must |
| P-06 | View consultation history (own cases) | Must |
| P-07 | View diagnosis for a closed case | Must |

### Doctor

| ID | Requirement | Priority |
|----|-------------|----------|
| D-01 | Authenticate (login/logout) | Must |
| D-02 | View assigned medical cases | Must |
| D-03 | Review uploaded images for assigned cases | Must |
| D-04 | Write diagnosis for an assigned case | Must |
| D-05 | View case history (assigned cases) | Must |

### Administrator

| ID | Requirement | Priority |
|----|-------------|----------|
| A-01 | Authenticate (login/logout) | Must |
| A-02 | Inspect all medical cases within tenant | Must |
| A-03 | Assign doctors to cases | Must |
| A-04 | Close diagnosed cases | Must |
| A-05 | Inspect audit logs | Must |
| A-06 | Add users to tenant (assign role) | Must |
| A-07 | Remove users from tenant | Must |
| A-08 | List tenant members | Must |
| A-09 | Reactivate suspended tenant | Should |
| A-10 | Manage tenant-level resources (future) | Could |

---

## Non-Functional Requirements

### Security

| ID | Requirement | Priority |
|----|-------------|----------|
| NF-S01 | All data encrypted at rest (AES-256) | Must |
| NF-S02 | All data encrypted in transit (TLS 1.2+) | Must |
| NF-S03 | JWT-based authentication with short-lived tokens | Must |
| NF-S04 | Role-based authorization on every endpoint | Must |
| NF-S05 | Tenant isolation on every data access | Must |
| NF-S06 | No hardcoded secrets in codebase | Must |
| NF-S07 | Secrets stored in AWS Secrets Manager | Must |
| NF-S08 | Audit logging for all state-changing operations | Must |
| NF-S09 | No PHI in logs | Must |
| NF-S10 | bcrypt password hashing (cost factor ≥ 12) | Must |
| NF-S11 | Rate limiting on authentication endpoints | Must |

### Performance

| ID | Requirement | Priority |
|----|-------------|----------|
| NF-P01 | API response time < 200ms (p95) | Should |
| NF-P02 | Image upload supports files up to 50MB | Must |
| NF-P03 | Concurrent user support: 100 per tenant | Should |

### Availability

| ID | Requirement | Priority |
|----|-------------|----------|
| NF-A01 | 99.9% uptime for API endpoints | Should |
| NF-A02 | Database automated backups (daily) | Must |
| NF-A03 | Database point-in-time recovery | Should |

### Compliance

| ID | Requirement | Priority |
|----|-------------|----------|
| NF-C01 | HIPAA-inspired architectural controls | Must |
| NF-C02 | No real patient data (PoC only) | Must |
| NF-C03 | Audit trail retention: 90 days minimum | Must |

### Scalability

| ID | Requirement | Priority |
|----|-------------|----------|
| NF-SC01 | Horizontal scaling via ECS Fargate | Must |
| NF-SC02 | Database connection pooling | Should |

### Infrastructure

| ID | Requirement | Priority |
|----|-------------|----------|
| NF-I01 | All resources managed via Terraform | Must |
| NF-I02 | Modular structure representing platform capabilities | Must |
| NF-I03 | Remote state in S3 with versioning and encryption | Must |
| NF-I04 | Security by default: private subnets, encryption, least privilege | Must |
| NF-I05 | No hardcoded values — variables for all configurable parameters | Must |

### Observability

| ID | Requirement | Priority |
|----|-------------|----------|
| NF-O01 | Structured JSON logging | Must |
| NF-O02 | CloudWatch Logs integration | Must |
| NF-O03 | Request tracing with correlation IDs | Should |

---

## API Requirements

> **Source of truth:** The full API contract (schemas, examples, validation rules) is defined in [`spec/openapi.yaml`](../spec/openapi.yaml). The endpoints below are a summary for quick reference.

### Endpoints

| Method | Endpoint | Auth | Role | Description |
|--------|----------|------|------|-------------|
| POST | /api/v1/auth/register | No | — | Register new user (no tenant) |
| POST | /api/v1/auth/login | No | — | Authenticate user, return available tenants |
| POST | /api/v1/auth/select-tenant | Yes | — | Select tenant, receive JWT with tenant + role |
| POST | /api/v1/auth/refresh | Yes | — | Refresh access token |
| GET | /api/v1/users/me | Yes | Any | Get current user profile |
| GET | /api/v1/tenants/members | Yes | Admin | List tenant members |
| POST | /api/v1/tenants/members | Yes | Admin | Add user to tenant with role |
| DELETE | /api/v1/tenants/members/{user_id} | Yes | Admin | Remove user from tenant |
| POST | /api/v1/tenants/{id}/reactivate | Yes | Admin | Reactivate suspended tenant |
| GET | /api/v1/cases | Yes | Patient, Doctor, Admin | List cases (filtered by role, filterable by status) |
| POST | /api/v1/cases | Yes | Patient | Create new case |
| GET | /api/v1/cases/{id} | Yes | Patient, Doctor, Admin | Get case details |
| POST | /api/v1/cases/{id}/symptoms | Yes | Patient | Add symptom to case |
| POST | /api/v1/cases/{id}/assign | Yes | Admin | Assign doctor to case |
| POST | /api/v1/cases/{id}/diagnosis | Yes | Doctor | Write diagnosis |
| POST | /api/v1/cases/{id}/close | Yes | Admin | Close case |
| POST | /api/v1/cases/{id}/images/upload-url | Yes | Patient | Request pre-signed upload URL |
| POST | /api/v1/cases/{id}/images | Yes | Patient | Confirm image upload |
| GET | /api/v1/cases/{id}/images | Yes | Patient, Doctor, Admin | List images for case |
| GET | /api/v1/images/{id}/download-url | Yes | Patient, Doctor, Admin | Get pre-signed download URL |
| GET | /api/v1/audit-logs | Yes | Admin | List audit logs |

### Response Format

```json
{
  "data": {},
  "error": null,
  "meta": {
    "request_id": "uuid",
    "timestamp": "2024-01-01T00:00:00Z"
  }
}
```

### Error Format

```json
{
  "data": null,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Description of the error",
    "details": []
  },
  "meta": {
    "request_id": "uuid",
    "timestamp": "2024-01-01T00:00:00Z"
  }
}
```

### Error Codes

| HTTP Status | Error Code | When |
|-------------|------------|------|
| 400 | `VALIDATION_ERROR` | Request body or parameters fail validation |
| 401 | `UNAUTHORIZED` | Missing, invalid, or expired token |
| 403 | `FORBIDDEN` | Valid token but insufficient role/permissions |
| 404 | `NOT_FOUND` | Resource does not exist or belongs to another tenant |
| 409 | `CONFLICT` | Duplicate resource (e.g., email already registered) |
| 422 | `BUSINESS_RULE_VIOLATION` | Domain invariant violated (e.g., close case without diagnosis) |
| 429 | `RATE_LIMITED` | Too many requests (retry-after header included) |
| 500 | `INTERNAL_ERROR` | Unexpected server error |

**Rules:**
- `404` is returned instead of `403` when a resource belongs to another tenant (prevents tenant enumeration)
- `error.details[]` contains field-level validation errors (field name + message)
- `error.code` is a machine-readable string, never changes between API versions
- `error.message` is human-readable, safe to display to end users
- `meta.request_id` is included in every response (success and error) for log correlation

---

## Out of Scope (PoC)

The following are explicitly out of scope for the initial PoC:

- Real-time notifications
- Video consultations
- Payment processing
- Multi-language support
- Mobile applications
- CI/CD pipeline
- Production deployment
- Load testing
- Penetration testing
- API versioning strategy (all endpoints use `/api/v1/` prefix; versioning policy deferred to post-PoC)
