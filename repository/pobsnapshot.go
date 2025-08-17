package repository

import (
	"time"

	"github.com/google/uuid"
)

func (r *Repository) CreatePOBSnapshot(characterId string, exportString string) error {
	query := `
	INSERT INTO pobsnapshots (id, character_id, export_string, created_at, updated_at)
	VALUES(?, ?, ?, ?, ?)
	`

	now := time.Now()
	idString := uuid.New().String()
	_, err := r.db.Exec(query, idString, characterId, exportString, now, now)
	return err
}
