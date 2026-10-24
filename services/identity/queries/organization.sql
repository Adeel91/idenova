
-- name: CreateOrganization :exec
INSERT INTO organizations (id, name)
VALUES ($1, $2);

-- name: GetOrganizationByID :one
SELECT id, name, created_at
FROM organizations
WHERE id = $1;
