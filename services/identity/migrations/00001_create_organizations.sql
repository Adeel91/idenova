
-- +goose Up

CREATE TABLE organizations (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT organization_name_not_empty
        CHECK (LENGTH(TRIM(name)) > 0)
);

-- +goose Down

DROP TABLE organizations;
