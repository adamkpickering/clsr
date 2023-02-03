package deck_source

import (
	"encoding/json"
	"fmt"
	"github.com/adamkpickering/clsr/internal/models"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type JSONFileDeckSource struct {
	baseDirectory string
}

func NewJSONFileDeckSource(baseDirectory string) (JSONFileDeckSource, error) {
	absoluteBaseDirectory, err := filepath.Abs(baseDirectory)
	if err != nil {
		return JSONFileDeckSource{}, fmt.Errorf("failed to get directory %q as absolute path: %w", baseDirectory, err)
	}
	// check that passed base directory is valid
	_, err = os.ReadDir(absoluteBaseDirectory)
	if err != nil {
		return JSONFileDeckSource{}, fmt.Errorf("problem with base directory %q: %w", baseDirectory, err)
	}

	deckSource := JSONFileDeckSource{
		baseDirectory: absoluteBaseDirectory,
	}
	return deckSource, nil
}

func (deckSource JSONFileDeckSource) ReadDeck(name string) (*models.Deck, error) {
	fileName := fmt.Sprintf("%s.json", name)
	deckPath := filepath.Join(deckSource.baseDirectory, fileName)

	// read deck file
	contents, err := os.ReadFile(deckPath)
	if err != nil {
		return &models.Deck{}, fmt.Errorf("failed to read deck: %w", err)
	}

	// decode contents into Deck struct
	deck := &models.Deck{}
	err = json.Unmarshal(contents, deck)
	if err != nil {
		return &models.Deck{}, fmt.Errorf("failed to parse contents of deck: %w", err)
	}

	// do any post-parse changes to cards that are needed
	for _, card := range deck.Cards {
		card.Deck = deck.Name
		sort.Stable(card.Reviews)
		for i := range card.Reviews {
			card.Reviews[i].Datetime = card.Reviews[i].Datetime.In(time.Local)
		}
	}

	return deck, nil
}

func (deckSource JSONFileDeckSource) WriteDeck(passedDeck *models.Deck) error {
	// copy deck and set location of datetimes to UTC
	deck := passedDeck.Copy()
	for _, card := range deck.Cards {
		for i := range card.Reviews {
			card.Reviews[i].Datetime = card.Reviews[i].Datetime.In(time.UTC)
		}
	}

	// marshal contents of deck file
	contents, err := json.MarshalIndent(deck, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Deck to JSON: %w", err)
	}

	// write deck file
	fileName := fmt.Sprintf("%s.json", deck.Name)
	deckPath := filepath.Join(deckSource.baseDirectory, fileName)
	err = os.WriteFile(deckPath, contents, 0644)
	if err != nil {
		return fmt.Errorf("failed to write deck to file: %w", err)
	}

	return nil
}

func (deckSource JSONFileDeckSource) ListDecks() ([]string, error) {
	dirEntries, err := os.ReadDir(deckSource.baseDirectory)
	if err != nil {
		return []string{}, fmt.Errorf("failed to read deck directory: %w", err)
	}

	deckNames := []string{}
	for _, dirEntry := range dirEntries {
		nodeName := dirEntry.Name()
		expectedExtension := ".json"
		if filepath.Ext(nodeName) == expectedExtension {
			deckName := strings.TrimSuffix(nodeName, expectedExtension)
			deckNames = append(deckNames, deckName)
		}
	}

	return deckNames, nil
}

func (deckSource JSONFileDeckSource) DeleteDeck(name string) error {
	fileName := fmt.Sprintf("%s.json", name)
	deckPath := filepath.Join(deckSource.baseDirectory, fileName)

	err := os.Remove(deckPath)
	if err != nil {
		return fmt.Errorf("failed to delete deck: %w", err)
	}

	return nil
}
