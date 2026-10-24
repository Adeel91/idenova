-- +goose Up

CREATE TABLE idempotency_keys (
    operation TEXT NOT NULL,
    idempotency_key TEXT NOT NULL,
    request_hash TEXT NOT NULL,
    resource_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (operation, idempotency_key)
);

CREATE INDEX idx_idempotency_keys_created_at
    ON idempotency_keys (created_at);

-- +goose Down

DROP TABLE IF EXISTS idempotency_keys;
