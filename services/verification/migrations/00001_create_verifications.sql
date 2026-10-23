-- +goose Up
CREATE TABLE verifications (
    id VARCHAR(255) PRIMARY KEY,
    customer_id VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_status
        CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED'))
);

-- +goose Down
DROP TABLE verifications;