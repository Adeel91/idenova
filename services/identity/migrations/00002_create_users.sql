
-- +goose Up

CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    organization_id VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'MEMBER',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_users_organization
        FOREIGN KEY (organization_id)
        REFERENCES organizations(id)
        ON DELETE CASCADE,

    CONSTRAINT valid_user_role
        CHECK (role IN ('ADMIN', 'MEMBER')),

    CONSTRAINT unique_organization_email
        UNIQUE (organization_id, email)
);

CREATE INDEX idx_users_organization_id
ON users (organization_id);

-- +goose Down

DROP TABLE users;
