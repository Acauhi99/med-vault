-- Domain event outbox for Transactional Outbox pattern
-- Events are persisted in the same transaction as the aggregate.
-- A poller reads unpublished events and dispatches to projection handlers.

CREATE TABLE domain_outbox (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type    VARCHAR(100) NOT NULL,
    aggregate_type VARCHAR(100) NOT NULL,
    aggregate_id  UUID         NOT NULL,
    tenant_id     UUID         NOT NULL REFERENCES tenants(id),
    payload       JSONB        NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    published     BOOLEAN      NOT NULL DEFAULT false,
    published_at  TIMESTAMPTZ,
    attempts      INT          NOT NULL DEFAULT 0,
    last_error    TEXT,
    failed        BOOLEAN      NOT NULL DEFAULT false
);

-- Poller query: unpublished events ordered by creation time
CREATE INDEX idx_outbox_unpublished ON domain_outbox (created_at)
    WHERE published = false AND failed = false;

-- Tenant filtering for audit/debug
CREATE INDEX idx_outbox_tenant_id ON domain_outbox (tenant_id);

-- Aggregate filtering for ordering guarantees
CREATE INDEX idx_outbox_aggregate ON domain_outbox (aggregate_type, aggregate_id, created_at);

-- Cleanup: published events older than 30 days
-- (handled by application-level job or pg_cron)
