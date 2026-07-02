# Requirements

This document defines the functional and non-functional requirements for MedVault.

---

## Functional Requirements

### Patient

| ID | Requirement | Priority |
|----|-------------|----------|
| P-01 | Register a new account within a tenant | Must |
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
| A-06 | Manage tenant-level resources (future) | Could |

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

### Observability

| ID | Requirement | Priority |
|----|-------------|----------|
| NF-O01 | Structured JSON logging | Must |
| NF-O02 | CloudWatch Logs integration | Must |
| NF-O03 | Request tracing with correlation IDs | Should |

---

## API Requirements

### Endpoints

| Method | Endpoint | Auth | Role | Description |
|--------|----------|------|------|-------------|
| POST | /api/v1/auth/register | No | — | Register new user |
| POST | /api/v1/auth/login | No | — | Authenticate user |
| POST | /api/v1/auth/refresh | Yes | — | Refresh access token |
| GET | /api/v1/users/me | Yes | Any | Get current user profile |
| GET | /api/v1/cases | Yes | Patient, Doctor, Admin | List cases (filtered by role) |
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
