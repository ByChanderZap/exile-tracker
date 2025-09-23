package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ByChanderZap/exile-tracker/buildsSitesClient"
	"github.com/ByChanderZap/exile-tracker/config"
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

	err = fs.CreateSnapshot(c.ID, itemsResponse, passivesResponse)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create snapshot")
		return
	}
}

func (fs *FetcherService) CreateSnapshot(characterId string, items models.ItemsResponse, passives models.PassiveSkillsResponse) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return errors.Join(err, errors.New("failed to get current directory"))
	}

	dirPath := filepath.Join(currentDir, characterId)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return errors.Join(err, errors.New("failed to create character directory"))
	}

	itemsPath := filepath.Join(dirPath, "items.json")
	file, err := os.OpenFile(itemsPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Join(err, errors.New("error trying to open items file"))
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(items)
	if err != nil {
		return errors.Join(err, errors.New("something went wrong while encoding json items"))
	}

	passivesPath := filepath.Join(dirPath, "passives.json")
	passivesFile, err := os.OpenFile(passivesPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Join(err, errors.New("error trying to open passives file"))
	}
	defer passivesFile.Close()

	encoder2 := json.NewEncoder(passivesFile)
	err = encoder2.Encode(passives)
	if err != nil {
		return errors.Join(err, errors.New("something went wrong while encoding json passives"))
	}

	result, err := fs.generatePoBBin(itemsPath, passivesPath)
	if err != nil {
		return errors.Join(err, errors.New("failed to execute PoB"))
	}

	// Clean up after execution
	os.RemoveAll(dirPath)

	dbSnapshot, err := fs.repo.GetLatestSnapshotByCharacter(characterId)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return errors.Join(err, errors.New("something went wrong while getting snapshots"))
		}
		fs.log.Warn().Msg("No previous snapshots found.")
	}

	if result == dbSnapshot.ExportString {
		return errors.New("no changes detected between latest and current snapshot")
	}

	err = fs.repo.CreatePOBSnapshot(repository.CreatePoBSnapshotParams{
		CharacterId:  characterId,
		ExportString: result,
	})
	if err != nil {
		return errors.Join(err, errors.New("something went wrong while trying to store snapshot"))
	}

	return nil
}

func (fs *FetcherService) generatePoBBin(itemsPath string, passivesPath string) (string, error) {
	fs.log.Info().Msg("Executing Path of Building in headless mode")
	pobRoot := config.Envs.POBRoot

	srcDir := filepath.Join(pobRoot, "src")

	runtimeLua := filepath.Join(pobRoot, "runtime", "lua")
	runtime := filepath.Join(pobRoot, "runtime")

	os.Setenv("LUA_PATH", runtimeLua+"/?.lua;"+runtimeLua+"/?/init.lua;;")
	os.Setenv("LUA_CPATH", runtime+"/?.so;"+runtime+"/?.dll;;")

	// Use absolute paths for JSON files to avoid any path resolution issues
	cmd := exec.Command("/usr/bin/luajit", "HeadlessWrapper.lua", itemsPath, passivesPath)
	cmd.Dir = srcDir // Set working directory for this command only

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Join(err, errors.New("the command execution failed"))
	}

	lines := strings.Split(string(output), "\n")
	last := lines[len(lines)-2]

	uploadedBuild, err := buildsSitesClient.UploadBuild(last, buildsSitesClient.SitesUrl.PoeNinja)
	if err != nil {
		return "", errors.Join(err, errors.New("failed when uploading build"))
	}
	fs.log.Debug().Msg(uploadedBuild)
	return uploadedBuild, nil
}
