package models

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var CardTemplate *template.Template

//go:embed card.txt.tmpl
var cardTemplateText string

func init() {
	CardTemplate = template.Must(template.New("card").Parse(cardTemplateText))
}

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
		if entry.IsDir() {
			continue
		}
		cardPath := filepath.Join(deckPath, entry.Name())
		card, err := ReadCardFromFile(cardPath)
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
		name := GetCardFilename(card)
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
		cardPath := filepath.Join(deckPath, GetCardFilename(card))
		err := WriteCardToFile(cardPath, card)
		if err != nil {
			return fmt.Errorf("failed to write card %s: %s", card.ID, err)
		}
	}

	return nil
}

func GetCardFilename(card *Card) string {
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

func ReadCardFromFile(path string) (*Card, error) {
	// get id out of file name
	filename := filepath.Base(path)
	id := strings.Split(filename, ".")[0]

	// read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return &Card{}, fmt.Errorf("failed to read card file: %s", err)
	}

	// parse the card and return
	card, err := ParseCardFromString(string(data), id)
	if err != nil {
		return &Card{}, fmt.Errorf("failed to parse card file: %s", err)
	}
	return card, nil
}

func WriteCardToFile(path string, card *Card) error {
	// process card into a map
	outputCard := map[string]string{}
	outputCard["Version"] = fmt.Sprintf("%d", card.Version)
	outputCard["LastReview"] = card.LastReview.Format(dateLayout)
	outputCard["NextReview"] = card.NextReview.Format(dateLayout)
	outputCard["Active"] = fmt.Sprintf("%t", card.Active)
	outputCard["Question"] = card.Question
	outputCard["Answer"] = card.Answer

	// open/create file
	fd, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to open card file for writing: %s", err)
	}
	defer fd.Close()

	// fill template and write to file
	err = CardTemplate.Execute(fd, outputCard)
	if err != nil {
		return fmt.Errorf("failed to fill card template: %s", err)
	}

	return nil
}
