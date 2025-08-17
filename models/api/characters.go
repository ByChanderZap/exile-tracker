package models

type CreateCharacterInput struct {
	AccountId     string `json:"account_id"`
	CharacterName string `json:"CharacterName"`
	Died          bool   `json:"died"`
	CurrentLeague string `json:"current_league"`
}

type UpdateCharacterInput struct {
	CharacterName string `json:"character_name"`
	CurrentLeague string `json:"current_league"`
}

type AddCharacterToFetchInput struct {
	CharacterId string `json:"character_id"`
}
