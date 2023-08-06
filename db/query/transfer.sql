-- name: CreateTransfer :one
INSERT INTO transfers(
    from_account_id,
    to_account_id,
    amount
)
values($1,$2,$3)
RETURNING *;

-- name: GetTransferFromId :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;



-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1;