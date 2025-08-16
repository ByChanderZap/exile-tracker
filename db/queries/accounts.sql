-- name: GetAccountById :one
SELECT * FROM accounts WHERE id = ? AND deleted_at IS NULL;

-- name: ListAccounts :many
SELECT * FROM accounts WHERE deleted_at IS NULL;

-- name: CreateAccount :one
INSERT INTO accounts (
  id, 
  account_name,
  created_at,
  updated_at
) 
VALUES (
  ?, ?, ?, ?
)
RETURNING *;
