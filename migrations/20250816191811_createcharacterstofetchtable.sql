-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS characters_to_fetch(
  id TEXT PRIMARY KEY,
  character_id TEXT NOT NULL,
  last_fetch TEXT,
  should_skip boolean NOT NULL DEFAULT false,

  FOREIGN KEY(character_id) REFERENCES characters(id) ON DELETE CASCADE 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS characters_to_fetch;
-- +goose StatementEnd
