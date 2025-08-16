-- name: GetCharacterByID :one
SELECT * FROM characters WHERE id = ? AND deleted_at IS NULL;

-- name: ListCharactersByAccount :many
SELECT * FROM characters WHERE account_id = ? AND deleted_at IS NULL;

-- name: CreateCharacter :one
INSERT INTO characters (
  id, 
  account_id,
  character_name,
  died,
  created_at,
  updated_at
) 
VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: KillCharacter :one
UPDATE characters
SET died = TRUE
where id = ?
