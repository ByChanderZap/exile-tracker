-- +goose Up
-- +goose StatementBegin
ALTER TABLE characters ADD COLUMN current_league TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- SQLite does not support DROP COLUMN directly.
-- To revert, you would need to recreate the table without the column if necessary.
SELECT 'down SQL query: cannot drop column in SQLite, manual intervention required if needed.';
-- +goose StatementEnd
