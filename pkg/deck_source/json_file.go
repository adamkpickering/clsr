package deck_source

import (
	"encoding/json"
	"filepath"
	"fmt"
	"github.com/adamkpickering/clsr/pkg/models"
	"os"
)

type JSONFileDeckSource struct {
	baseDirectory string
}

func NewJSONFileDeckSource(baseDirectory string) (JSONFileDeckSource, error) {
	// check that passed base directory is valid
	_, err := os.ReadDir(baseDirectory)
	if err != nil {
		return JSONFileDeckSource{}, fmt.Errorf("problem with base directory %s: %w", baseDirectory, err)
	}

	deckSource := JSONFileDeckSource{
		baseDirectory: baseDirectory,
	}
	return deckSource, nil
}

func (deckSource JSONFileDeckSource) ReadDeck(name string) (*models.Deck, error) {
	fileName := fmt.Sprintf("%s.json", name)
	deckPath := filepath.Join(deckSource.baseDirectory, fileName)

	// read deck file
	contents, err := os.ReadFile(deckPath)
	if err != nil {
		return &models.Deck{}, fmt.Errorf("failed to read deck %s: %w", name, err)
	}

	// decode contents into Deck struct
	deck := &models.Deck{}
	err = json.Unmarshal(contents, deck)
	if err != nil {
		return &models.Deck{}, fmt.Errorf("failed to parse contents of deck: %w", name, err)
	}
	return deck, nil
}

func (deckSource JSONFileDeckSource) WriteDeck(deck *models.Deck) error {
	fileName := fmt.Sprintf("%s.json", deck.Name)
	deckPath := filepath.Join(deckSource.baseDirectory, fileName)

	// marshal contents of deck file
	contents, err := json.MarshalIndent(deck, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Deck to JSON: %w")
	}

	// write deck file
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
		deckNames = append(deckNames, dirEntry.Name())
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
