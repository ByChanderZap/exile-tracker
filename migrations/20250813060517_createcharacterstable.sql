-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS characters (
  id              TEXT PRIMARY KEY,
  account_id      TEXT,
  character_name  TEXT NOT NULL,
  died            boolean NOT NULL DEFAULT false,
  
  created_at      TIMESTAMP NOT NULL,
  updated_at      TIMESTAMP NOT NULL,
  deleted_at      TIMESTAMP,
  FOREIGN KEY(account_id) REFERENCES accounts(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS characters;
-- +goose StatementEnd
