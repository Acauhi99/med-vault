# ADR-017: Transactional Outbox for Domain Events

## Status

Accepted

## Context

MedVault uses CQRS with domain events to bridge the write and read sides. Events like `CaseCreated`, `DiagnosisWritten`, and `ImageUploaded` trigger projections that update read models. The system needs a reliable way to deliver these events without losing any, even if the application crashes mid-operation.

**Requirements:**
- Events must not be lost (aggregate state and event delivery must be consistent)
- Read models may be eventually consistent (acceptable latency: ~1s)
- All modules run in the same process (modular monolith)
- No external infrastructure beyond PostgreSQL

## Decision

Use the **Transactional Outbox** pattern with polling, implemented entirely in Go + PostgreSQL.

### How It Works

```
Write Path:
  CommandHandler.Handle()
    → aggregate.Mutate()
    → BEGIN TRANSACTION
        → repo.Save(aggregate)
        → outbox.Save(event)        ← same transaction
    → COMMIT

Read Path (async):
  OutboxPoller (goroutine, every 1s)
    → SELECT FROM domain_outbox WHERE published = false
    → for each event:
        → projectionHandler.Handle(event)
        → mark as published
```

### Outbox Table

```sql
domain_outbox (
    id, event_type, aggregate_type, aggregate_id, tenant_id,
    payload (JSONB), created_at, published, published_at,
    attempts, last_error, failed
)
```

### Guarantees

| Guarantee | Mechanism |
|-----------|-----------|
| No event loss | Event persisted in same transaction as aggregate |
| At-least-once delivery | Poller retries on failure; handler must be idempotent |
| Ordering per aggregate | Events ordered by `created_at` within same aggregate |
| Failure handling | After 3 failed attempts, mark as `failed` for manual inspection |
| Cleanup | Published events deleted after 30 days |

### Idempotency

Handlers must be idempotent — the same event may be delivered twice (poller publishes, crashes before marking as published). Each handler checks if the projection already exists before inserting.

```go
func (h *CaseProjection) Handle(evt CaseCreated) error {
    if h.repo.Exists(evt.CaseID) {
        return nil // already processed
    }
    return h.repo.Insert(...)
}
```

### Retry Strategy

```
Attempt 1: immediate
Attempt 2: after 1s
Attempt 3: after 5s
Attempt 4+: mark as failed, log error
```

## Consequences

### Positive
- No event loss — events survive application crashes
- No external infrastructure (no SQS, RabbitMQ, Kafka)
- Debuggable — events persist in the outbox table for inspection
- Consistent with modular monolith architecture
- Simple to implement (goroutine + SQL)
- Works with PostgreSQL on RDS (no special extensions needed)

### Negative
- At-least-once (not exactly-once) — handlers must be idempotent
- Polling adds latency (~1s, configurable)
- Extra table to manage and clean up
- Outbox table grows until cleanup runs

## Alternatives Considered

| Alternative | Reason for Rejection |
|-------------|---------------------|
| In-process event bus (no persistence) | Events lost on crash; no real CQRS guarantee |
| Synchronous dispatch (same transaction) | Write and read coupled; not true CQRS |
| SQS / SNS | External infrastructure; overkill for modular monolith |
| Kafka / MSK | High complexity; not justified for PoC |
| CDC (Debezium) | Requires additional infrastructure; overkill |

## When to Add a Message Broker

If the architecture evolves to microservices (separate containers/processes), replace the in-process poller with:
- **SQS** — simple queue, at-least-once
- **SNS + SQS** — fan-out to multiple consumers
- **EventBridge** — complex routing rules
- **Kafka (MSK)** — high-throughput, replay, streaming

The outbox table stays — the poller just publishes to the broker instead of calling in-process handlers.

## References

- [Transactional Outbox Pattern](https://microservices.io/patterns/data/transactional-outbox.html)
- [Reliable Microservices Data Exchange with the Outbox Pattern](https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/)
- [DOMAIN.md](../DOMAIN.md) — Domain events and CQRS mapping
