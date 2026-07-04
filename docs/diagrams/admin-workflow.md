# Administrator Workflow

## Complete Admin Journey

```mermaid
flowchart TD
    A[Admin visits MedVault] --> B[Login]
    B --> C[Select tenant]
    C --> D[Admin area]

    D --> E[Cases page /cases]
    D --> F[Members page /members]
    D --> G[Audit logs page /audit]

    E --> E1[View all cases]
    E1 --> E2[Filter by status]
    E1 --> E3[Assign doctor to case]
    E3 --> E4[Case status: Open → Assigned]
    E1 --> E5[Close diagnosed case]
    E5 --> E6[Case status: Diagnosed → Closed]

    F --> F1[List tenant members]
    F --> F2[Add user to tenant]
    F --> F3[Remove user from tenant]
    F --> F4[Assign role to user]

    G --> G1[View all audit logs]
    G1 --> G2[Filter by user / resource / action]
```

## Assign Doctor to Case

```mermaid
sequenceDiagram
    actor Admin
    participant Frontend
    participant Backend
    participant DB

    Admin->>Frontend: Open case → Assign Doctor
    Frontend->>Backend: GET /api/v1/tenants/members
    Backend->>Backend: Validate JWT + role = admin
    Backend->>DB: SELECT users WHERE tenant_id = $1 AND role = 'doctor'
    DB-->>Backend: List of doctors
    Backend-->>Frontend: 200 OK {doctors[]}

    Admin->>Frontend: Select doctor
    Frontend->>Backend: POST /api/v1/cases/{id}/assign {doctor_id}
    Backend->>Backend: Validate JWT + role = admin
    Backend->>Backend: Validate case status = Open
    Backend->>Backend: Validate doctor belongs to same tenant
    Backend->>DB: BEGIN TRANSACTION
    Backend->>DB: UPDATE case SET doctor_id = $1, status = 'assigned'
    Backend->>DB: INSERT domain_outbox (DoctorAssigned event)
    Backend->>DB: INSERT audit_log (DoctorAssigned)
    Backend->>DB: COMMIT
    Backend-->>Frontend: 200 OK {case}
    Frontend-->>Admin: Doctor assigned
```

## Close Diagnosed Case

```mermaid
sequenceDiagram
    actor Admin
    participant Frontend
    participant Backend
    participant DB

    Admin->>Frontend: Open diagnosed case → Close
    Frontend->>Backend: POST /api/v1/cases/{id}/close
    Backend->>Backend: Validate JWT + role = admin
    Backend->>Backend: Validate case status = Diagnosed
    Backend->>DB: BEGIN TRANSACTION
    Backend->>DB: UPDATE case SET status = 'closed', closed_at = now()
    Backend->>DB: INSERT domain_outbox (CaseClosed event)
    Backend->>DB: INSERT audit_log (CaseClosed)
    Backend->>DB: COMMIT
    Backend-->>Frontend: 200 OK {case}
    Frontend-->>Admin: Case closed
```

## Add User to Tenant

```mermaid
sequenceDiagram
    actor Admin
    participant Frontend
    participant Backend
    participant DB

    Admin->>Frontend: Click "Add Member"
    Admin->>Frontend: Enter user email + select role
    Frontend->>Backend: POST /api/v1/tenants/members {user_id, role}
    Backend->>Backend: Validate JWT + role = admin
    Backend->>DB: SELECT user WHERE id = $1
    DB-->>Backend: User exists
    Backend->>DB: SELECT user_tenants WHERE user_id = $1 AND tenant_id = $2
    DB-->>Backend: Check not already member
    Backend->>DB: BEGIN TRANSACTION
    Backend->>DB: INSERT user_tenants (user_id, tenant_id, role)
    Backend->>DB: INSERT domain_outbox (UserAddedToTenant event)
    Backend->>DB: INSERT audit_log (UserAddedToTenant)
    Backend->>DB: COMMIT
    Backend-->>Frontend: 201 Created
    Frontend-->>Admin: User added to tenant
```

## Remove User from Tenant

```mermaid
sequenceDiagram
    actor Admin
    participant Frontend
    participant Backend
    participant DB

    Admin->>Frontend: Select member → Remove
    Frontend->>Backend: DELETE /api/v1/tenants/members/{user_id}
    Backend->>Backend: Validate JWT + role = admin
    Backend->>DB: BEGIN TRANSACTION
    Backend->>DB: DELETE user_tenants WHERE user_id = $1 AND tenant_id = $2
    Backend->>DB: INSERT domain_outbox (UserRemovedFromTenant event)
    Backend->>DB: INSERT audit_log (UserRemovedFromTenant)
    Backend->>DB: COMMIT
    Backend-->>Frontend: 204 No Content
    Frontend-->>Admin: User removed
```

## List Tenant Members

```mermaid
sequenceDiagram
    actor Admin
    participant Frontend
    participant Backend
    participant DB

    Admin->>Frontend: Click "Members"
    Frontend->>Backend: GET /api/v1/tenants/members
    Backend->>Backend: Validate JWT + role = admin
    Backend->>DB: SELECT users JOIN user_tenants WHERE tenant_id = $1
    DB-->>Backend: List of members with roles
    Backend-->>Frontend: 200 OK {members[]}
    Frontend-->>Admin: Show member list
```

## View Audit Logs

```mermaid
sequenceDiagram
    actor Admin
    participant Frontend
    participant Backend
    participant DB

    Admin->>Frontend: Click "Audit Logs"
    Frontend->>Backend: GET /api/v1/audit-logs?action=&user_id=&resource_type=&resource_id=
    Backend->>Backend: Validate JWT + role = admin
    Backend->>DB: SELECT audit_logs WHERE tenant_id = $1 AND ... ORDER BY timestamp DESC
    DB-->>Backend: Audit log entries
    Backend-->>Frontend: 200 OK {audit_logs[]}
    Frontend-->>Admin: Show audit log table
```

## Reactivate Suspended Tenant

```mermaid
sequenceDiagram
    actor Admin
    participant Frontend
    participant Backend
    participant DB

    Admin->>Frontend: Click "Reactivate Tenant"
    Frontend->>Backend: POST /api/v1/tenants/{id}/reactivate
    Backend->>Backend: Validate JWT + role = admin
    Backend->>DB: BEGIN TRANSACTION
    Backend->>DB: UPDATE tenants SET status = 'active' WHERE id = $1
    Backend->>DB: INSERT audit_log (TenantReactivated)
    Backend->>DB: COMMIT
    Backend-->>Frontend: 200 OK {tenant}
    Frontend-->>Admin: Tenant reactivated
```
