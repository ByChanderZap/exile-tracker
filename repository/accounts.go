package repository

import (
	"database/sql"
	"time"

	"github.com/ByChanderZap/exile-tracker/models"
	"github.com/google/uuid"
)

func (r *Repository) GetAllAccounts() ([]models.Account, error) {
	query := "SELECT id, account_name, player FROM accounts WHERE deleted_at IS NULL"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var acc models.Account
		err := rows.Scan(&acc.ID, &acc.AccountName, &acc.Player)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (r Repository) CreateAccount(accountName string, player string) error {
	query := `
	INSERT INTO accounts (id, account_name, player, created_at, updated_at)
	VALUES(?, ?, ?, ?, ?)
	`

	now := time.Now()
	idString := uuid.New().String()

	_, err := r.db.Exec(query, idString, accountName, player, now, now)
	return err
}

func (r *Repository) GetAccountByID(id string) (models.Account, error) {
	query := `
	SELECT id, account_name, player
	FROM accounts
	WHERE id = ?
	`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return models.Account{}, err
	}
	defer rows.Close()

	var a models.Account
	if rows.Next() {
		err := rows.Scan(
			&a.ID,
			&a.AccountName,
			&a.Player,
		)
		if err != nil {
			return models.Account{}, nil
		}
		return a, nil
	}
	return models.Account{}, sql.ErrNoRows
}
