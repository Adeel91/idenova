-- +goose Up

CREATE TABLE documents (
    id VARCHAR(255) PRIMARY KEY,
    verification_id VARCHAR(255) NOT NULL,
    document_type VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT valid_document_status
        CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED')),

    CONSTRAINT document_type_not_empty
        CHECK (LENGTH(TRIM(document_type)) > 0)
);

CREATE INDEX idx_documents_verification_id
ON documents (verification_id);

-- +goose Down

DROP TABLE documents;
