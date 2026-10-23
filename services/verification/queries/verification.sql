-- name: CreateVerification :exec
INSERT INTO verifications (id, customer_id, status)
VALUES ($1, $2, $3);

-- name: GetVerification :one
SELECT id, customer_id, status, created_at
FROM verifications
WHERE id = $1;

-- name: UpdateVerificationStatus :execrows
UPDATE verifications
SET status = $2
WHERE id = $1 AND status = 'PENDING';

-- name: ListVerifications :many
SELECT id, customer_id, status, created_at
FROM verifications
WHERE (sqlc.narg('status')::text IS NULL
       OR status = sqlc.narg('status')::text)
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;