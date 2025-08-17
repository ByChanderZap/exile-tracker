package models

type CreateAccountInput struct {
	AccountName string  `json:"account_name"`
	Player      *string `json:"player"`
}

type UpdateAccountInput struct {
	AccountName string  `json:"account_name"`
	Player      *string `json:"player"`
}
