# Security

> **This is the source of truth for security controls, encryption, and compliance.** Other documents (INFRASTRUCTURE.md, ARCHITECTURE.md) reference this document rather than duplicating its content.

This document defines the security architecture, controls, and threat model for MedVault.

> **Disclaimer:** MedVault is a Proof of Concept. This document defines controls aligned with the HIPAA Privacy Rule (45 CFR §164.500–534), Security Rule (45 CFR §164.302–318), and Breach Notification Rule (45 CFR §164.400–414). It is not a substitute for a formal HIPAA compliance program or legal review.

---

## Protected Health Information (PHI)

### Definition

PHI is any individually identifiable health information transmitted or maintained in any form or medium. Under HIPAA, 18 identifiers must be removed for data to be considered de-identified (Safe Harbor method, 45 CFR §164.514(b)).

### PHI in MedVault

| Data Element | Contains PHI | Identifiers Present |
|--------------|--------------|---------------------|
| Case (symptoms, status) | Yes | PatientID, TenantID |
| Diagnosis (notes) | Yes | DoctorID, PatientID |
| Medical images | Yes | Linked to PatientID via Case |
| Audit logs | Yes | UserID, TenantID, IP address |
| User accounts | Indirectly | Email (identifier) |
| Tenant configuration | No | Organization-level only |

### Minimum Necessary Standard

Each role accesses only the minimum PHI required for its function:

| Role | PHI Accessible |
|------|----------------|
| Patient | Own cases, own symptoms, own images, own diagnoses |
| Doctor | Assigned cases, images for assigned cases, diagnoses written |
| Administrator | All cases within tenant (for management), audit logs (no diagnosis content) |

### De-identification

When PHI must be de-identified (e.g., for analytics, reporting, or research):

- **Safe Harbor Method (45 CFR §164.514(b)):** Remove all 18 identifiers
- **Expert Determination (45 CFR §164.514(c)):** Statistical expert certifies re-identification risk is very small
- De-identified data is no longer considered PHI under HIPAA

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
| SQL injection | High | Parameterized queries via pgx + sqlc code generation |
| Cross-site scripting (XSS) | Medium | Input validation, output encoding |
| Cross-site request forgery (CSRF) | Medium | SameSite cookies, CSRF tokens |
| Brute force | Medium | Rate limiting on auth endpoints, account lockout |
| Session hijacking | Medium | Short-lived JWT, client session state |
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
- **Token storage:** client session state (in-memory for the current browser session)
- **Token transmission:** Authorization header for API calls
- **Login flow:** Two-step (authenticate → select tenant)

### JWT Claims

```json
{
  "sub": "user_id",
  "tenant_id": "selected_tenant_id",
  "role": "patient | doctor | administrator",
  "iat": 1704067200,
  "exp": 1704068100
}
```

The `role` is per-tenant, from the `user_tenants` table. A user may have different roles in different tenants.

### Login Flow

```
Step 1 — Authenticate:
POST /auth/login { email, password }
→ Temporary JWT (no tenant) + list of available tenants with roles

Step 2 — Select tenant:
POST /auth/select-tenant { tenant_id }  (with temporary JWT)
→ Final JWT { user_id, tenant_id, role }
```

To switch tenants, the user calls select-tenant again without re-authenticating.

### Authentication Telemetry

- Successful login continues into the audit trail once a user is identified.
- Failed login is recorded as a structured security log with request ID, email, IP, and user agent.
- Failed login is not stored in `audit_logs` because it has no trusted tenant context.

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
  4. Backend issues new access + refresh tokens (same tenant_id + role)
  5. Client retries original request
```

### Automatic Logoff

HIPAA requires automatic logoff after a period of inactivity (45 CFR §164.312(a)(2)(iii)).

| Setting | Value | Rationale |
|---------|-------|-----------|
| Access token lifetime | 15 minutes | Short-lived, automatic expiration |
| Refresh token lifetime | 7 days | Balances security and usability |
| Inactivity timeout (frontend) | 15 minutes | Automatic session termination |
| Re-authentication required | After inactivity timeout | User must authenticate again |

**Implementation:**
- Frontend currently clears the in-memory session on explicit sign-out
- Inactivity timeout is not wired in yet
- Refresh tokens are invalidated on the server

### Emergency Access Procedure

HIPAA requires emergency access procedures (45 CFR §164.312(a)(2)(ii)).

**Break-Glass Procedure:**

1. **Trigger:** Normal authentication is unavailable or insufficient for emergency patient care
2. **Authorization:** Only the Security Officer or designated administrator can authorize emergency access
3. **Documentation:** Emergency access must be documented immediately including:
   - Reason for emergency access
   - User who accessed the system
   - Date and time of access
   - PHI accessed
4. **Review:** All emergency access events are reviewed within 24 hours
5. **Revocation:** Emergency access is revoked once normal operations resume

**Emergency Access Roles:**

| Role | Emergency Capability |
|------|---------------------|
| Security Officer | Authorize emergency access; review and revoke |
| System Administrator | Execute emergency access; maintain audit trail |
| On-Call Physician | Request emergency access for patient care |

### Tenant Switch Flow

```
Client → Backend
  1. User selects a different tenant from the UI
  2. Client calls POST /auth/select-tenant { tenant_id }
  3. Backend validates user belongs to that tenant
  4. Backend issues new JWT with new tenant_id + role
  5. Client replaces stored token
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

## Privacy Rule (45 CFR §164.500–534)

### Notice of Privacy Practices (NPP)

Patients must receive a written notice describing:

- How their PHI may be used and disclosed
- Their rights regarding their PHI
- The organization's legal duties
- How to file a complaint

**Implementation:** NPP is presented during registration and available in the patient portal at all times.

### Patient Rights

| Right | HIPAA Reference | Implementation |
|-------|-----------------|----------------|
| Access to PHI | §164.524 | Patient can view own cases, symptoms, images, diagnoses via API |
| Amendment of PHI | §164.526 | Patient can request amendment; admin reviews and responds within 60 days |
| Accounting of Disclosures | §164.528 | Audit log tracks all PHI access; patient can request report of disclosures |
| Request Restrictions | §164.522(a) | Patient can request restrictions on uses/disclosures (admin approves) |
| Confidential Communications | §164.522(b) | Patient can request alternative communication methods |
| Paper Copy of NPP | §164.520 | Available on request |

### Uses and Disclosures

| Category | Requires Authorization | Examples |
|----------|----------------------|----------|
| Treatment, Payment, Operations (TPO) | No | Doctor views case for diagnosis; admin processes billing |
| Required by Law | No | Court orders, public health reporting |
| Patient Authorization | Yes | Research, marketing, sale of PHI, psychotherapy notes |
| Minimum Necessary | Always | Each role accesses only what is needed (see PHI section) |

### Business Associate Agreements (BAAs)

Any third party that creates, receives, maintains, or transmits PHI must sign a BAA.

| Vendor | Service | BAA Required |
|--------|---------|--------------|
| AWS | RDS, S3, ECS, CloudWatch, CloudFront | Yes |
| AWS | Route 53 (DNS only, no PHI) | No |

**BAA Contents (45 CFR §164.504(e)):**
- Establish permitted and required uses of PHI
- Require safeguards to prevent unauthorized use/disclosure
- Report breaches and security incidents
- Ensure subcontractors agree to same restrictions
- Return or destroy PHI upon contract termination

---

## Breach Notification Rule (45 CFR §164.400–414)

### Breach Definition

A breach is the acquisition, access, use, or disclosure of PHI in a manner not permitted by the Privacy Rule that compromises the security or privacy of the PHI.

**Exceptions (not considered breaches):**
- Unintentional acquisition by workforce member acting in good faith
- Inadvertent disclosure between authorized persons
- Disclosure where recipient could not reasonably retain PHI

### Breach Assessment

When a suspected breach is discovered:

1. **Contain** — Immediately mitigate the breach
2. **Assess** — Determine if PHI was involved and the risk level
3. **Document** — Record all findings in the breach log
4. **Notify** — Follow notification requirements below

### Risk Assessment Factors

| Factor | Consideration |
|--------|---------------|
| Nature and extent of PHI | Types of identifiers, clinical data involved |
| Unauthorized person | Who accessed the data (internal vs external) |
| Whether PHI actually acquired/viewed | Evidence of actual access vs potential access |
| Extent of risk mitigation | How quickly the breach was contained |

### Notification Requirements

| Notification | Deadline | Content |
|--------------|----------|---------|
| Individual Notice | Within 60 days of discovery | Description, types of PHI, steps to protect, investigation summary, contact info |
| HHS (≥500 individuals) | Within 60 days of discovery | Breach details via HHS portal |
| HHS (<500 individuals) | Annual log (by March 1) | Annual summary of small breaches |
| Media (≥500 in a state) | Within 60 days of discovery | Same as individual notice |

### Breach Response Team

| Role | Responsibility |
|------|----------------|
| Security Officer | Leads investigation, coordinates response |
| Privacy Officer | Determines notification requirements |
| Legal Counsel | Advises on legal obligations |
| IT/Infrastructure | Contains breach, preserves evidence |
| Communications | Drafts and sends notifications |

### Breach Documentation

All breaches must be documented regardless of size:

```json
{
  "breach_id": "uuid",
  "discovered_at": "2024-01-01T00:00:00Z",
  "reported_at": "2024-01-01T00:00:00Z",
  "description": "Description of the breach",
  "phi_involved": ["case", "diagnosis"],
  "individuals_affected": 100,
  "root_cause": "Compromised credentials",
  "mitigation_steps": ["Password reset", "Session invalidation"],
  "notifications_sent": {
    "individuals": true,
    "hhs": true,
    "media": false
  },
  "status": "resolved"
}
```

**Retention:** Breach documentation must be retained for **6 years** from the date of creation.

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
| S3 medical images | AWS KMS CMK | Customer managed |
| S3 audit logs | AWS KMS CMK | Customer managed |
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

| Log Type | Retention Period | HIPAA Requirement |
|----------|-----------------|-------------------|
| Application logs | 90 days (CloudWatch) | Not specifically required |
| Audit logs | 6 years (S3 with lifecycle) | 45 CFR §164.530(j) — documentation retention |
| CloudTrail | 90 days (management events) | Not specifically required |
| Breach documentation | 6 years | 45 CFR §164.530(j) |
| Security incident records | 6 years | 45 CFR §164.530(j) |

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
- Full request logging at edge controls is minimized unless needed for active investigation, reducing unnecessary PHI/header exposure

### Data Retention

- Patient data retained per tenant configuration
- Audit logs retained for compliance
- Soft deletes for data recovery

### Data Disposal

- Soft delete with retention period
- Hard delete after retention period (future)
- S3 lifecycle policies for old images

---

## Administrative Safeguards (45 CFR §164.308)

### Security Management Process

| Standard | Implementation |
|----------|----------------|
| Risk Analysis (§164.308(a)(1)(ii)(A)) | Threat model documented; risk assessment performed annually |
| Risk Management (§164.308(a)(1)(ii)(B)) | Security controls implemented based on risk analysis; cost-benefit analysis for each control |
| Sanction Policy (§164.308(a)(1)(ii)(C)) | Workforce members who violate security policies face disciplinary action up to termination |
| Information System Activity Review (§164.308(a)(1)(ii)(D)) | Audit logs reviewed weekly; anomalies investigated within 48 hours |

### Assigned Security Responsibility

| Role | Responsibility | HIPAA Reference |
|------|----------------|-----------------|
| Security Officer | Overall HIPAA Security Rule compliance; leads risk analysis; approves security policies | §164.308(a)(2) |
| Privacy Officer | Privacy Rule compliance; manages NPP; handles patient rights requests | §164.308(a)(2) |
| IT Administrator | Implements technical controls; manages access provisioning | §164.308(a)(2) |

### Workforce Security

| Standard | Implementation |
|----------|----------------|
| Authorization/Supervision (§164.308(a)(3)(ii)(A)) | Access granted based on role; supervisor approval required |
| Workforce Clearance (§164.308(a)(3)(ii)(B)) | Background checks for workforce members with PHI access |
| Termination Procedures (§164.308(a)(3)(ii)(C)) | Access revoked within 24 hours of termination; credentials disabled; devices returned |

### Information Access Management

| Standard | Implementation |
|----------|----------------|
| Access Authorization (§164.308(a)(4)(ii)(B)) | RBAC enforced; access granted per minimum necessary standard |
| Access Modification (§164.308(a)(4)(ii)(C)) | Access changes logged and reviewed; modification requires approval |

### Security Awareness and Training

| Standard | Implementation |
|----------|----------------|
| Security Reminders (§164.308(a)(5)(ii)(A)) | Monthly security awareness communications |
| Protection from Malicious Software (§164.308(a)(5)(ii)(B)) | Workforce trained on malware prevention; endpoint protection deployed |
| Log-in Monitoring (§164.308(a)(5)(ii)(C)) | Failed login attempts tracked; account lockout after 5 failures |
| Password Management (§164.308(a)(5)(ii)(D)) | Password policy enforced; complexity requirements documented |

### Security Incident Procedures

| Standard | Implementation |
|----------|----------------|
| Response and Reporting (§164.308(a)(6)(ii)) | Incident response plan documented; breach response team defined |

**Incident Response Process:**

1. **Detection** — Monitor audit logs, alerts, and reports
2. **Containment** — Isolate affected systems immediately
3. **Eradication** — Remove threat and vulnerabilities
4. **Recovery** — Restore systems and verify integrity
5. **Lessons Learned** — Document findings and update controls
6. **Reporting** — Notify affected parties per Breach Notification Rule

### Contingency Plan (§164.308(a)(7))

| Standard | Implementation |
|----------|----------------|
| Data Backup Plan (§164.308(a)(7)(ii)(A)) | RDS automated backups daily; S3 versioning enabled; backups encrypted |
| Disaster Recovery Plan (§164.308(a)(7)(ii)(B)) | Recovery procedures documented; RTO: 4 hours; RPO: 1 hour |
| Emergency Mode Operation (§164.308(a)(7)(ii)(C)) | Emergency access procedures documented; read-only access available during outages |
| Testing and Revision (§164.308(a)(7)(ii)(D)) | Contingency plan tested annually; results documented |
| Applications and Data Criticality (§164.308(a)(7)(ii)(E)) | Critical systems identified; priority restoration order defined |

**Critical Systems Priority:**

| Priority | System | RTO | RPO |
|----------|--------|-----|-----|
| 1 | Database (RDS PostgreSQL) | 1 hour | 15 minutes |
| 2 | Backend API (ECS Fargate) | 2 hours | N/A |
| 3 | Medical Image Storage (S3) | 4 hours | 1 hour |
| 4 | Frontend (CloudFront + S3) | 4 hours | N/A |

### Evaluation (§164.308(a)(8))

- Security controls evaluated annually
- Penetration testing performed by qualified third party
- Vulnerability scanning performed quarterly
- Results documented and remediation tracked

---

## Physical Safeguards (45 CFR §164.310)

### Facility Access Controls (§164.310(a))

| Standard | Implementation |
|----------|----------------|
| Contingency Operations (§164.310(a)(2)(i)) | AWS handles physical facility contingency |
| Facility Security Plan (§164.310(a)(2)(ii)) | AWS data center security; see AWS compliance programs |
| Access Control and Validation (§164.310(a)(2)(iii)) | AWS SOC 2 Type II certified; physical access restricted |
| Maintenance Records (§164.310(a)(2)(iv)) | AWS maintains physical maintenance logs |

### Workstation Use (§164.310(b))

- Workstations accessing PHI must be in secure areas
- Screen lock required after 15 minutes of inactivity
- PHI must not be displayed on unattended screens
- Removable media prohibited for PHI storage

### Workstation Security (§164.310(c))

- Endpoint protection (antivirus, firewall) required
- Operating system and software kept current with patches
- Full disk encryption required on devices with PHI access
- Remote wipe capability for mobile devices

### Device and Media Controls (§164.310(d))

| Standard | Implementation |
|----------|----------------|
| Disposal (§164.310(d)(2)(i)) | Electronic media sanitized before disposal; NIST 800-88 guidelines |
| Media Re-use (§164.310(d)(2)(ii)) | All PHI removed before media re-use |
| Accountability (§164.310(d)(2)(iii)) | Device inventory maintained; movement tracked |
| Data Backup and Storage (§164.310(d)(2)(iv)) | Backups stored encrypted; access restricted to authorized personnel |

---

## AWS Shared Responsibility Model

This project operates under the AWS Shared Responsibility Model.

### AWS Responsibilities

| Area | AWS Manages |
|------|-------------|
| Physical security | Data center access, hardware lifecycle |
| Network infrastructure | VPC backbone, availability zones |
| Managed services | RDS, S3, ECS, CloudWatch availability |
| Compliance certifications | SOC, ISO, HIPAA BAA |

### Infrastructure Responsibilities (Terraform)

| Area | We Manage |
|------|-----------|
| Network design | VPC, subnets, route tables, security groups |
| Encryption | RDS encryption, S3 encryption, KMS keys |
| IAM | Roles, policies, least privilege |
| Secrets | Secrets Manager configuration, rotation |
| Logging | CloudTrail, VPC Flow Logs, CloudWatch |
| Backup | RDS backups, S3 versioning |

### Application Responsibilities (Go/Next.js)

| Area | We Manage |
|------|-----------|
| Tenant isolation | `WHERE tenant_id = $1` in every query |
| RBAC | Role checks on every endpoint |
| PHI handling | Never in logs, never in errors |
| Audit logging | Business operation audit trail |
| Input validation | Zod schemas, request validation |
| Authentication | JWT, bcrypt, token lifecycle |

### Developer Responsibilities

| Area | We Manage |
|------|-----------|
| Code quality | Reviews, testing, documentation |
| Secret rotation | Awareness, manual trigger |
| Incident response | Detection, mitigation, reporting |
| Least privilege | IAM policy review, no wildcards |

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

### Technical Safeguards
- [x] No hardcoded secrets
- [x] All endpoints require authentication
- [x] All endpoints enforce RBAC
- [x] All queries include tenant_id
- [x] Passwords hashed with bcrypt
- [x] JWT tokens have short expiration
- [x] Audit logging on all mutations
- [x] PHI excluded from logs
- [x] TLS enforced on all connections
- [x] S3 bucket has no public access
- [x] RDS is in private subnet
- [x] Secrets stored in Secrets Manager
- [x] Automatic logoff after 15 minutes inactivity
- [x] Emergency access procedures documented

### Privacy Rule
- [x] Notice of Privacy Practices defined
- [x] Patient rights documented (access, amendment, accounting, restrictions)
- [x] Minimum Necessary Standard enforced per role
- [x] Business Associate Agreement with AWS
- [x] Uses and Disclosures policy documented
- [x] De-identification methods documented

### Breach Notification Rule
- [x] Breach definition documented
- [x] Breach assessment process defined
- [x] Risk assessment factors documented
- [x] Notification requirements and deadlines defined
- [x] Breach response team identified
- [x] Breach documentation template defined
- [x] 6-year retention for breach records

### Administrative Safeguards
- [x] Security Officer designated
- [x] Privacy Officer designated
- [x] Risk analysis methodology documented
- [x] Risk management process defined
- [x] Sanction policy documented
- [x] Workforce security procedures (background checks, termination)
- [x] Security awareness and training program
- [x] Incident response procedures documented
- [x] Contingency plan (backup, disaster recovery, emergency mode)
- [x] Annual security evaluation scheduled

### Physical Safeguards
- [x] Facility access controls (AWS managed)
- [x] Workstation security policy documented
- [x] Workstation use policy documented
- [x] Device and media controls documented
- [x] Automatic logoff configured

### Documentation Retention
- [x] Audit logs retained for 6 years
- [x] Breach documentation retained for 6 years
- [x] Security incident records retained for 6 years
- [x] Policy documentation retained for 6 years
