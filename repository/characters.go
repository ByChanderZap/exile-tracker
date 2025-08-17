package repository

import (
	"database/sql"
	"time"

	"github.com/ByChanderZap/exile-tracker/models"
	"github.com/google/uuid"
)

func (r *Repository) GetCharactersByAccountId(accountId string) ([]models.Character, error) {
	query := `
	SELECT id, account_id, character_name, died, current_league, created_at, updated_at
	FROM characters
	WHERE account_id = ? AND deleted_at IS NULL
	`
	rows, err := r.db.Query(query, accountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []models.Character
	for rows.Next() {
		var char models.Character
		err := rows.Scan(
			&char.ID,
			&char.AccountId,
			&char.CharacterName,
			&char.Died,
			&char.CurrentLeague,
			&char.CreatedAt,
			&char.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		characters = append(characters, char)
	}

	return characters, nil
}

func (r *Repository) CreateCharacter(accountId string, characterName string) error {
	query := `
		INSERT INTO characters(id, account_id, character_name, created_at, updated_at)
		VALUES(?,?,?,?,?)
	`

	now := time.Now()
	idString := uuid.New().String()
	_, err := r.db.Exec(query, idString, accountId, characterName, now, now)

	return err
}

func (r *Repository) UpdateDiedStatus(characterId string, died bool) error {
	query := `
		UPDATE characters SET died = ?, updated_at = ? 
		WHERE id = ?
	`

	now := time.Now()

	_, err := r.db.Exec(query, died, now, characterId)

	return err
}

func (r *Repository) GetCharactersToFetch() ([]models.CharactersToFetch, error) {
	query := `
		SELECT id, character_id, last_fetch, should_skip
		FROM characters_to_fetch
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cToFetch []models.CharactersToFetch
	for rows.Next() {
		var char models.CharactersToFetch
		err := rows.Scan(
			&char.Id,
			&char.CharacterId,
			&char.LastFetch,
			&char.ShouldSkip,
		)
		if err != nil {
			return nil, err
		}
		cToFetch = append(cToFetch, char)
	}
	return cToFetch, nil
}

func (r *Repository) SetShouldSkip(shouldSkip bool, id string) error {
	query := `
		UPDATE characters_to_fetch
		SET should_skip = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, shouldSkip, time.Now(), id)
	return err
}

func (r *Repository) GetCharacterByID(id string) (models.Character, error) {
	query := `
	SELECT id, account_id, character_name, died, current_league
	FROM characters
	WHERE id = ?
	`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return models.Character{}, err
	}
	defer rows.Close()

	var c models.Character
	if rows.Next() {
		err := rows.Scan(
			&c.ID,
			&c.AccountId,
			&c.CharacterName,
			&c.Died,
			&c.CurrentLeague,
		)
		if err != nil {
			return models.Character{}, nil
		}
		return c, nil
	}
	return models.Character{}, sql.ErrNoRows
}
