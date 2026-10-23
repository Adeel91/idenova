-- name: CreateVerification :exec
INSERT INTO verifications (id, customer_id, status)
VALUES ($1, $2, $3);

-- name: GetVerification :one
SELECT id, customer_id, status, created_at
FROM verifications
WHERE id = $1;