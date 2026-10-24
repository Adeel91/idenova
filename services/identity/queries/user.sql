
-- name: CreateUser :exec
INSERT INTO users (id, organization_id, email, role)
VALUES ($1, $2, $3, $4);

-- name: GetUserByID :one
SELECT id, organization_id, email, role, created_at
FROM users
WHERE id = $1;

-- name: ListUsersByOrganization :many
SELECT id, organization_id, email, role, created_at
FROM users
WHERE organization_id = $1
ORDER BY created_at DESC, id DESC;
