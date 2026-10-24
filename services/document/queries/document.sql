-- name: CreateDocument :exec
INSERT INTO documents (id, verification_id, document_type, status)
VALUES ($1, $2, $3, $4);

-- name: GetDocumentByID :one
SELECT id, verification_id, document_type, status, created_at
FROM documents
WHERE id = $1;

-- name: ListDocumentsByVerification :many
SELECT id, verification_id, document_type, status, created_at
FROM documents
WHERE verification_id = $1
ORDER BY created_at DESC, id DESC;

-- name: UpdateDocumentStatus :execrows
UPDATE documents
SET status = $2
WHERE id = $1 AND status = 'PENDING';
