-- name: CreateSession :one
INSERT INTO sessions(
    id,username,refresh_token ,user_agent ,client_ip ,is_blocked,expired_at
    ) 
values( $1,$2,$3,$4,$5,$6,$7)
RETURNING *;

-- name: GetSessionFromId :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;