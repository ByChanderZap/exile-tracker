-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS accounts (
  id              TEXT PRIMARY KEY,
  account_name    TEXT NOT NULL,
  player          TEXT,

  created_at      TIMESTAMP NOT NULL,
  updated_at      TIMESTAMP NOT NULL,
  deleted_at      TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS accounts;
-- +goose StatementEnd
