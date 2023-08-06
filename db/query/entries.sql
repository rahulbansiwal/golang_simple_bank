-- name: CreateEntry :one 
INSERT INTO entries(
    account_id,amount
)
values($1,$2)
RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: GetEntriesFromAccountId :many
SELECT * FROM entries
LIMIT $1
OFFSET $2;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;

