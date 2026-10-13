-- name: UpsertGuild :one
INSERT INTO guilds (id)
VALUES ($1)
ON CONFLICT (id) DO UPDATE SET id = EXCLUDED.id -- If guild already exists, do a harmless update
RETURNING *;

-- name: GetGuild :one
SELECT * FROM guilds WHERE id = $1;

-- name: SetGuildPrefix :exec
UPDATE guilds SET prefix = $2 WHERE id = $1;

-- name: DeleteGuild :exec
DELETE FROM guilds WHERE id = $1;