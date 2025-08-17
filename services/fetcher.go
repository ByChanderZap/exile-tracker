package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ByChanderZap/exile-tracker/models"
	"github.com/ByChanderZap/exile-tracker/poeclient"
	"github.com/ByChanderZap/exile-tracker/repository"
	"github.com/ByChanderZap/exile-tracker/utils"
	"github.com/rs/zerolog"
)

type FetcherService struct {
	repo      *repository.Repository
	poeClient *poeclient.POEClient
	log       zerolog.Logger
	ticker    *time.Ticker
	done      chan bool
}

func NewFetcherService(repo *repository.Repository, poeClient *poeclient.POEClient, interval time.Duration) *FetcherService {
	return &FetcherService{
		repo:      repo,
		poeClient: poeClient,
		log:       utils.ChildLogger("fetcher"),
		ticker:    time.NewTicker(interval),
		done:      make(chan bool),
	}
}

func (fs *FetcherService) Start(ctx context.Context) {
	fs.log.Info().Msg("Starting fetcher service")

	// Run once first
	go fs.fetchAllData()

	go func() {
		for {
			select {
			case <-fs.done:
				fs.log.Info().Msg("Fetcher service stopped")
				return
			case <-fs.ticker.C:
				fs.fetchAllData()
			case <-ctx.Done():
				fs.log.Info().Msg("Context cancelled, stopping fetcher service")
			}
		}
	}()
}

func (fs *FetcherService) Stop() {
	fs.ticker.Stop()
	fs.done <- true
}

func (fs *FetcherService) fetchAllData() {
	fs.log.Info().Msg("Starting fetch cycle")

	charactersToFetch, err := fs.repo.GetCharactersToFetch()
	if err != nil {
		fs.log.Error().Err(err).Msg("Faile to get characters to fetch from database")
		return
	}

	fs.log.Info().Int("characters_to_fetch", len(charactersToFetch)).Msg("Found characters to process")

	for _, ctf := range charactersToFetch {
		fs.FetchCharacterData(ctf)
		time.Sleep(2 * time.Second)
	}
	fs.log.Info().Msg("Data fetch cycle completed")
}

func (fs *FetcherService) FetchCharacterData(ctf models.CharactersToFetch) {
	c, err := fs.repo.GetCharacterByID(ctf.CharacterId)
	if err != nil {
		fs.log.Error().Err(err).Str("character_id", ctf.CharacterId).Msg("Failed to fetch")
		return
	}

	acc, err := fs.repo.GetAccountByID(c.AccountId)
	if err != nil {
		fs.log.Error().Err(err).Str("account_id", c.AccountId).Msg("Failed to fetch")
		return
	}

	// Creating a logger with the "context" of the character that is beign processed
	log := fs.log.With().
		Str("account", acc.AccountName).
		Str("character", c.CharacterName).Logger()

	if c.Died {
		log.Warn().Msg("Character is dead, skipping fetch")
		fs.repo.SetShouldSkip(true, ctf.Id)
		return
	}

	items, err := fs.poeClient.GetItemsJson(acc.AccountName, c.CharacterName, "pc")
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch items")
		return
	}

	passives, err := fs.poeClient.GetPassiveSkillsJson(acc.AccountName, c.CharacterName, "pc")
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch passive skills")
		return
	}

	var itemsResponse models.ItemsResponse
	if err := json.Unmarshal(items, &itemsResponse); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshall items")
		return
	}

	var passivesResponse models.PassiveSkillsResponse
	if err := json.Unmarshal(passives, &passivesResponse); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshall passive skills")
		return
	}

	fs.CreateSnapshot(c.ID, itemsResponse, passivesResponse)
}

func (fs *FetcherService) CreateSnapshot(characterId string, items models.ItemsResponse, passives models.PassiveSkillsResponse) {
	fs.log.Info().Msg("Not implemented yet")
}
