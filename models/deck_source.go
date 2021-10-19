package models

import (
	"fmt"
	"os"
	"path/filepath"
)

type DeckSource interface {
	CreateDeck(name string) (*Deck, error)
	ListDecks() ([]string, error)
	DeleteDeck(name string) error
	LoadDeck(name string) (*Deck, error)
	SyncDeck(deck *Deck) error
}

type FlatFileDeckSource struct {
	baseDir string
}

func NewFlatFileDeckSource(baseDir string) (FlatFileDeckSource, error) {
	deckSource := FlatFileDeckSource{
		baseDir: baseDir,
	}
	return deckSource, nil
}

func (deckSource FlatFileDeckSource) CreateDeck(name string) (*Deck, error) {
	deckPath := filepath.Join(deckSource.baseDir, name)
	err := os.MkdirAll(deckPath, 0755)
	if err != nil {
		return &Deck{}, fmt.Errorf("failed to create new deck: %s", err)
	}
	deck, err := NewDeck(name)
	if err != nil {
		return &Deck{}, fmt.Errorf("failed to construct *Deck: %s", err)
	}
	return deck, nil
}

func (deckSource FlatFileDeckSource) ListDecks() ([]string, error) {
	deckNames := []string{}
	entries, err := os.ReadDir(deckSource.baseDir)
	if err != nil {
		newErr := fmt.Errorf("could not read base directory: %s", err)
		return []string{}, newErr
	}
	for _, entry := range entries {
		if entry.IsDir() {
			deckNames = append(deckNames, entry.Name())
		}
	}
	return deckNames, nil
}

func (deckSource FlatFileDeckSource) DeleteDeck(name string) error {
	deckPath := filepath.Join(deckSource.baseDir, name)
	err := os.Remove(deckPath)
	if err != nil {
		return fmt.Errorf("failed to remove deck: %s", err)
	}
	return nil
}

func (deckSource FlatFileDeckSource) LoadDeck(name string) (*Deck, error) {
	// get a list of card paths
	deckPath := filepath.Join(deckSource.baseDir, name)
	entries, err := os.ReadDir(deckPath)
	if err != nil {
		return &Deck{}, fmt.Errorf("failed to read deck %s: %s", name, err)
	}

	// load each card in the deck
	deck := &Deck{
		Name: name,
	}
	for _, entry := range entries {
		entryInfo, err := entry.Info()
		if err != nil {
			return &Deck{}, fmt.Errorf("failed to get info for entry %s: %s", entry.Name(), err)
		}
		if entryInfo.IsDir() {
			continue
		}
		cardPath := filepath.Join(deckPath, entry.Name())
		card, err := parseCardFromFile(cardPath)
		if err != nil {
			return &Deck{}, fmt.Errorf("failed to parse card %s: %s", cardPath, err)
		}
		deck.Cards = append(deck.Cards, card)
	}

	return deck, nil
}

func (deckSource FlatFileDeckSource) SyncDeck(deck *Deck) error {
	deckPath := filepath.Join(deckSource.baseDir, deck.Name)

	// delete any cards that are not in the Deck
	currentCardFilenames, err := getDirFilenames(deckPath)
	if err != nil {
		return fmt.Errorf("failed to read deck directory: %s", err)
	}
	newCardFilenames := map[string]struct{}{}
	for _, card := range deck.Cards {
		name := getCardFilename(card)
		newCardFilenames[name] = struct{}{}
	}
	for _, currentCardFilename := range currentCardFilenames {
		_, ok := newCardFilenames[currentCardFilename]
		if !ok {
			cardPath := filepath.Join(deckPath, currentCardFilename)
			err := os.Remove(cardPath)
			if err != nil {
				return fmt.Errorf("failed to remove %s: %s", cardPath, err)
			}
		}
	}

	// write remaining cards
	for _, card := range deck.Cards {
		err := card.WriteToDir(deckPath)
		if err != nil {
			return fmt.Errorf("failed to write card %s: %s", card.ID, err)
		}
	}

	return nil
}

func getCardFilename(card *Card) string {
	return fmt.Sprintf("%s.txt", card.ID)
}

func getDirFilenames(path string) ([]string, error) {
	// get files/dirs in directory
	entries, err := os.ReadDir(path)
	if err != nil {
		return []string{}, fmt.Errorf("failed to read directory: %s", err)
	}

	// get names of only files
	filenames := []string{}
	for _, entry := range entries {
		entryInfo, err := entry.Info()
		if err != nil {
			return []string{}, fmt.Errorf("failed to get info for entry %s: %s", entry.Name(), err)
		}
		if entryInfo.IsDir() {
			continue
		}
		filenames = append(filenames, entry.Name())
	}

	return filenames, nil
}