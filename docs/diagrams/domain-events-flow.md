# Domain Events Flow (Transactional Outbox)

## Outbox Pattern Overview

```mermaid
flowchart LR
    subgraph Command Side
        A[Command Handler] --> B[Aggregate]
        B --> C[Mutate State]
    end

    subgraph Transaction
        D[(PostgreSQL)]
        C --> D
        E[Domain Event] --> D
    end

    subgraph Outbox Poller
        F[SELECT unpublished events]
        G[Dispatch to handlers]
        H[Mark as published]
    end

    subgraph Read Side
        I[Projection Handler]
        J[(Read Model)]
    end

    D -->|~1s polling| F
    F --> G
    G --> I
    I --> J
    G --> H
```

## Event Flow Sequence

```mermaid
sequenceDiagram
    actor Caller
    participant Handler as Command Handler
    participant Aggregate
    participant DB as PostgreSQL
    participant Outbox as Outbox Poller
    participant Projection as Projection Handler
    participant ReadModel as Read Model

    Caller->>Handler: Execute command
    Handler->>Aggregate: Validate + mutate
    Aggregate-->>Handler: Domain events

    Handler->>DB: BEGIN TRANSACTION
    Handler->>DB: INSERT/UPDATE aggregate state
    Handler->>DB: INSERT domain_outbox (event, published = false)
    Handler->>DB: COMMIT

    Note over Outbox: Polls every ~1 second

    Outbox->>DB: SELECT FROM domain_outbox WHERE published = false
    DB-->>Outbox: Unpublished events

    loop For each event
        Outbox->>Projection: Dispatch event
        Projection->>Projection: Process event (idempotent)
        Projection->>ReadModel: Update read model
        Outbox->>DB: UPDATE domain_outbox SET published = true
    end
```

## Domain Events Catalog

```mermaid
flowchart TD
    subgraph Identity & Access
        TE[TenantCreated]
        TS[TenantSuspended]
        TR[TenantReactivated]
        UR[UserRegistered]
        UAT[UserAddedToTenant]
        URF[UserRemovedFromTenant]
        UD[UserDeactivated]
    end

    subgraph Clinical
        CC[CaseCreated]
        SA[SymptomAdded]
        DA[DoctorAssigned]
        DW[DiagnosisWritten]
        CS[CaseClosed]
    end

    subgraph Imaging
        IU[ImageUploaded]
        ID[ImageDeleted]
    end

    subgraph Audit
        ALR[AuditLogRecorded]
    end

    TE --> ALR
    TS --> ALR
    TR --> ALR
    UR --> ALR
    UAT --> ALR
    URF --> ALR
    UD --> ALR
    CC --> ALR
    SA --> ALR
    DA --> ALR
    DW --> ALR
    CS --> ALR
    IU --> ALR
    ID --> ALR
```

## Event → Projection Mapping

| Event | Source Context | Projection Target |
|-------|---------------|-------------------|
| UserRegistered | Identity | UserReadModel |
| UserAddedToTenant | Identity | UserTenantReadModel |
| UserRemovedFromTenant | Identity | UserTenantReadModel |
| CaseCreated | Clinical | CaseReadModel |
| SymptomAdded | Clinical | CaseReadModel |
| DoctorAssigned | Clinical | CaseReadModel |
| DiagnosisWritten | Clinical | CaseReadModel |
| CaseClosed | Clinical | CaseReadModel |
| ImageUploaded | Imaging | ImageReadModel |
| AuditLogRecorded | Audit | AuditLogReadModel |

## Delivery Guarantees

| Property | Value |
|----------|-------|
| Delivery guarantee | At-least-once |
| Ordering | Per-aggregate (same aggregate = ordered) |
| Latency | ~1 second (polling interval) |
| Idempotency | Required (projection handlers must be idempotent) |
| Failed events | `attempts` counter incremented; retried with backoff |
| Max retries | Configurable (default: 5) |
| Dead letter | Events exceeding max retries moved to dead-letter table |

## Event Payload Structure

```json
{
  "event_id": "uuid",
  "event_type": "CaseCreated",
  "aggregate_type": "Case",
  "aggregate_id": "uuid",
  "tenant_id": "uuid",
  "occurred_at": "2024-01-01T00:00:00Z",
  "payload": {
    "case_id": "uuid",
    "patient_id": "uuid",
    "status": "open"
  }
}
```
