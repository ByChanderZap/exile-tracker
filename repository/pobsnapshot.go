package repository

import (
	"time"

	"github.com/ByChanderZap/exile-tracker/models"
	"github.com/google/uuid"
)

func (r *Repository) CreatePOBSnapshot(characterId string, exportString string) error {
	query := `
	INSERT INTO pobsnapshots (id, character_id, export_string, created_at, updated_at)
	VALUES(?, ?, ?, ?, ?)
	`

	now := time.Now().UTC().Format(time.RFC3339)
	idString := uuid.New().String()

	_, err := r.db.Exec(query, idString, characterId, exportString, now, now)
	return err
}

func (r *Repository) GetSnapshotsByCharacter(characterId string) ([]models.POBSnapshot, error) {
	query := `
	SELECT id, character_id, export_string, created_at, updated_at, deleted_at
	FROM pobsnapshots
	WHERE character_id = ?
	ORDER BY created_at ASC
	`
	rows, err := r.db.Query(query, characterId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snapshots []models.POBSnapshot
	for rows.Next() {
		var s models.POBSnapshot
		err := rows.Scan(
			&s.ID,
			&s.CharacterId,
			&s.ExportString,
			&s.CreatedAt,
			&s.UpdatedAt,
			&s.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return snapshots, nil
}

func (r *Repository) GetLatestSnapshotByCharacter(characterId string) (models.POBSnapshot, error) {
	query := `
	SELECT id, character_id, export_string, created_at, updated_at, deleted_at
	FROM pobsnapshots 
	WHERE character_id = ?
	ORDER BY created_at DESC
	LIMIT 1
	`
	var s models.POBSnapshot
	err := r.db.QueryRow(query, characterId).Scan(
		&s.ID,
		&s.CharacterId,
		&s.ExportString,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.DeletedAt,
	)
	if err != nil {
		return models.POBSnapshot{}, err
	}
	return s, nil
}

func (r *Repository) GetSnapshotByID(id string) (models.POBSnapshot, error) {
	query := `
	SELECT id, character_id, export_string, created_at, updated_at, deleted_at
	FROM pobsnapshots
	WHERE id = ?
	`
	var s models.POBSnapshot
	err := r.db.QueryRow(query, id).Scan(
		&s.ID,
		&s.CharacterId,
		&s.ExportString,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.DeletedAt,
	)
	if err != nil {
		return models.POBSnapshot{}, err
	}
	return s, nil
}
