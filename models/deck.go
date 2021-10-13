package models

import (
	"fmt"
	"os"
	"path/filepath"
)

// A Deck is a collection of Cards. It is persisted to the filesystem
// as a directory containing files that represent its Cards.
type Deck struct {
	Name  string
	Cards []Card
}

func LoadDeck(path string) (*Deck, error) {
	// get a list of card paths
	entries, err := os.ReadDir(path)
	if err != nil {
		return &Deck{}, fmt.Errorf("failed to read directory %s: %s", path, err)
	}

	// load each card in the deck
	deck := &Deck{
		Name: filepath.Base(path),
	}
	for _, entry := range entries {
		entryInfo, err := entry.Info()
		if err != nil {
			return &Deck{}, fmt.Errorf("failed to get info for entry %s: %s", entry.Name(), err)
		}
		if entryInfo.IsDir() {
			continue
		}
		cardPath := filepath.Join(path, entry.Name())
		card, err := parseCardFromFile(cardPath)
		if err != nil {
			return &Deck{}, fmt.Errorf("failed to parse card %s: %s", cardPath, err)
		}
		deck.Cards = append(deck.Cards, card)
	}

	return deck, nil
}

func NewDeck(path string) (*Deck, error) {
	// check whether it exists
	if _, err := os.Stat(path); err == nil {
		return &Deck{}, fmt.Errorf("Location %s is already occupied", path)
	}

	// create any needed directories
	err := os.MkdirAll(path, os.FileMode(int(0755)))
	if err != nil {
		return &Deck{}, fmt.Errorf("Could not create directory %s: %s", path, err)
	}

	deck := &Deck{
		Name: filepath.Base(path),
	}
	return deck, nil
}

//func (deck *Deck) Write() error {
//	fmt.Prinln("this is deck")
//}
