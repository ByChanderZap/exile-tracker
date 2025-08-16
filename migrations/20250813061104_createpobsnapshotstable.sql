-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pobsnapshots (
  id              TEXT PRIMARY KEY,
  character_id    TEXT NOT NULL,
  export_string   TEXT NOT NULL,

  created_at      TIMESTAMP NOT NULL,
  updated_at      TIMESTAMP NOT NULL,
  deleted_at      TIMESTAMP,
  FOREIGN KEY(character_id) REFERENCES characters(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pobsnapshots;
-- +goose StatementEnd
