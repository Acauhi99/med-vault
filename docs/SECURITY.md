# Security

This document defines the security architecture, controls, and threat model for MedVault.

> **Disclaimer:** MedVault is a Proof of Concept. It demonstrates HIPAA-inspired architectural decisions but is not HIPAA certified or production-ready.

---

## Threat Model

### Assets to Protect

| Asset | Sensitivity | Classification |
|-------|-------------|----------------|
| Patient health information (PHI) | Critical | Protected |
| User credentials | Critical | Secret |
| JWT signing keys | Critical | Secret |
| Database credentials | Critical | Secret |
| Medical images | High | Protected |
| Audit logs | High | Protected |
| Tenant configuration | Medium | Internal |

### Threat Actors

| Actor | Motivation | Capability |
|-------|------------|------------|
| Unauthorized user | Access PHI | Low (web-based attacks) |
| Malicious insider | Exfiltrate data | Medium (legitimate access) |
| Compromised account | Lateral movement | Medium (valid credentials) |
| External attacker | Disrupt service | Medium (automated attacks) |

### Attack Vectors

| Vector | Risk | Mitigation |
|--------|------|------------|
| SQL injection | High | Parameterized queries, ORM |
| Cross-site scripting (XSS) | Medium | Input validation, output encoding |
| Cross-site request forgery (CSRF) | Medium | SameSite cookies, CSRF tokens |
| Brute force | Medium | Rate limiting, account lockout |
| Session hijacking | Medium | Short-lived JWT, httpOnly cookies |
| Man-in-the-middle | High | TLS 1.2+ everywhere |
| Data exfiltration | High | Encryption at rest, access controls |
| Privilege escalation | High | RBAC enforcement, tenant isolation |
| Insecure direct object reference | High | Tenant-scoped queries, authorization checks |

---

## Authentication

### JWT Strategy

- **Algorithm:** HMAC-SHA256 (symmetric) for PoC
- **Access token lifetime:** 15 minutes
- **Refresh token lifetime:** 7 days
- **Token storage:** httpOnly, Secure, SameSite=Strict cookies
- **Token transmission:** Authorization header for API calls

### JWT Claims

```json
{
  "sub": "user_id",
  "tenant_id": "tenant_id",
  "role": "patient | doctor | administrator",
  "iat": 1704067200,
  "exp": 1704068100
}
```

### Password Policy

- Minimum 12 characters
- At least one uppercase, one lowercase, one number, one special character
- bcrypt hashing with cost factor ≥ 12
- No password reuse (future)

### Token Refresh Flow

```
Client → Backend
  1. Access token expires (401 response)
  2. Client sends refresh token
  3. Backend validates refresh token
  4. Backend issues new access token
  5. Client retries original request
```

---

## Authorization

### Role-Based Access Control (RBAC)

| Role | Permissions |
|------|-------------|
| Patient | Create cases, add symptoms, upload images, view own cases |
| Doctor | View assigned cases, view images, write diagnoses |
| Administrator | Assign doctors, close cases, view all cases, view audit logs |

### Enforcement Points

1. **HTTP middleware** validates JWT and extracts role
2. **Route handlers** check role before executing business logic
3. **Domain commands** validate permissions within aggregate boundaries
4. **Repository queries** enforce tenant isolation

---

## Multi-Tenant Isolation

### Data Isolation

- Every table includes `tenant_id` column
- Every query includes `WHERE tenant_id = $1`
- Row-Level Security (RLS) as defense-in-depth
- No cross-tenant joins allowed

### Storage Isolation

- S3 paths prefixed with `/{tenant_id}/`
- Pre-signed URLs scoped to tenant path
- Bucket policy enforces tenant prefix

### Network Isolation

- Each tenant shares the same network (PoC)
- Future: per-tenant VPC peering for strict isolation

---

## Encryption

### At Rest

| Resource | Method | Key Management |
|----------|--------|----------------|
| RDS PostgreSQL | AES-256 | AWS KMS |
| S3 medical images | AES-256 (SSE-S3) | AWS managed |
| S3 audit logs | AES-256 (SSE-S3) | AWS managed |
| Secrets | AES-256 | AWS Secrets Manager |

### In Transit

| Connection | Protocol | Certificate |
|------------|----------|-------------|
| Client → CloudFront | TLS 1.2+ | ACM managed |
| CloudFront → ALB | TLS 1.2+ | ACM managed |
| ALB → ECS | HTTP (internal VPC) | — |
| ECS → RDS | TLS 1.2+ | RDS CA bundle |
| ECS → S3 | HTTPS | AWS managed |

### Application-Level

- JWT signing: HMAC-SHA256 with secret from Secrets Manager
- Password hashing: bcrypt with cost factor ≥ 12

---

## Secrets Management

### Storage

| Secret | Location | Rotation |
|--------|----------|----------|
| Database credentials | AWS Secrets Manager | Automatic (30 days) |
| JWT signing key | AWS Secrets Manager | Manual (90 days) |
| S3 bucket keys | AWS managed | Automatic |

### Code Rules

- Never hardcode secrets in source code
- Never log secret values
- Never expose secrets in error messages
- Use environment variables or Secrets Manager SDK
- Pre-commit hooks scan for secrets (future)

---

## Audit Logging

### What to Log

| Event | Details |
|-------|---------|
| Authentication | Login success/failure, token refresh, logout |
| Authorization | Access denied, role check failures |
| Data mutations | Create, update, delete operations |
| File operations | Upload, download, delete |
| Admin operations | Tenant management, user management |

### Log Entry Format

```json
{
  "timestamp": "2024-01-01T00:00:00Z",
  "user_id": "uuid",
  "tenant_id": "uuid",
  "action": "case.created",
  "resource_type": "case",
  "resource_id": "uuid",
  "ip_address": "192.168.1.1",
  "metadata": {}
}
```

### Retention

- Application logs: 90 days (CloudWatch)
- Audit logs: 1 year (S3 with lifecycle)
- CloudTrail: 90 days (management events)

---

## Network Security

### VPC Design

```
┌─────────────────────────────────────────────────────┐
│                      VPC                            │
├─────────────────┬─────────────────┬─────────────────┤
│   Public Subnet │  Private Subnet │  Private Subnet │
│   (ALB only)    │  (ECS Fargate)  │  (RDS, S3)     │
└─────────────────┴─────────────────┴─────────────────┘
```

### Security Groups

| Group | Inbound | Outbound |
|-------|---------|----------|
| ALB SG | 443 from 0.0.0.0/0 | 8080 to ECS SG |
| ECS SG | 8080 from ALB SG | 5432 to RDS SG, 443 to S3 |
| RDS SG | 5432 from ECS SG | None |

### Network ACLs

- Default deny-all inbound
- Allow HTTP/HTTPS from ALB subnet
- Allow PostgreSQL from ECS subnet

---

## Data Protection

### PHI Handling

- Never store PHI in logs
- Never include PHI in error messages
- Never expose PHI in API responses without authorization
- Medical images encrypted at rest and in transit

### Data Retention

- Patient data retained per tenant configuration
- Audit logs retained for compliance
- Soft deletes for data recovery

### Data Disposal

- Soft delete with retention period
- Hard delete after retention period (future)
- S3 lifecycle policies for old images

---

## Compliance (HIPAA-Inspired)

### Shared Responsibility Model

| AWS Responsibility | Application Responsibility |
|-------------------|---------------------------|
| Physical security | Application-level encryption |
| Network infrastructure | Access controls |
| Database encryption | Tenant isolation |
| Secrets management | Audit logging |
| Compliance certifications | PHI handling procedures |

### Key Controls

| Control | Implementation |
|---------|----------------|
| Access control | RBAC + tenant isolation |
| Audit controls | Comprehensive audit logging |
| Integrity controls | Referential integrity, validation |
| Transmission security | TLS 1.2+ everywhere |
| Encryption | AES-256 at rest, TLS in transit |
| Emergency access | Administrator override (documented) |

---

## Security Checklist

- [ ] No hardcoded secrets
- [ ] All endpoints require authentication
- [ ] All endpoints enforce RBAC
- [ ] All queries include tenant_id
- [ ] Passwords hashed with bcrypt
- [ ] JWT tokens have short expiration
- [ ] Audit logging on all mutations
- [ ] PHI excluded from logs
- [ ] TLS enforced on all connections
- [ ] S3 bucket has no public access
- [ ] RDS is in private subnet
- [ ] Secrets stored in Secrets Manager
