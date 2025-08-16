-- name: GetLastSnapshotFromCharacter :one
SELECT * FROM pobsnapshots WHERE character_id = ? ORDER BY created_at DESC LIMIT 1 AND deleted_at IS NULL;

-- name: ListSnapshotsByCharacter :many
SELECT * FROM pobsnapshots WHERE character_id = ? AND deleted_at IS NULL;

-- name: CreateSnapshot :one
INSERT INTO pobsnapshots (
  id, 
  character_id,
  export_string,
  created_at,
  updated_at
) 
VALUES (
  ?, ?, ?, ?, ?
)
RETURNING *;
