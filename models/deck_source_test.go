package models

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadDeck(t *testing.T) {
	deckSource, err := NewFlatFileDeckSource("testdata")
	if err != nil {
		t.Errorf("failed to instantiate FlatFileDeckSource: %s", err)
	}
	deck, err := deckSource.LoadDeck("deck1")
	if err != nil {
		t.Errorf("failed to load deck: %s", err)
	}
	if length := len(deck.Cards); length != 2 {
		t.Errorf("read %d, not 2, cards for deck", length)
	}
}

func TestSyncDeck(t *testing.T) {
	// create a deck in temporary directory
	testDeckName := "test_deck"
	tempDir := t.TempDir()
	deckSource, err := NewFlatFileDeckSource(tempDir)
	if err != nil {
		t.Fatalf("failed to create deck source: %s", err)
	}
	initialDeck, err := deckSource.CreateDeck(testDeckName)
	if err != nil {
		t.Fatalf("failed to create deck: %s", err)
	}
	card1 := NewCard("card1 question", "card1 answer")
	card2 := NewCard("card2 question", "card2 answer")
	initialDeck.Cards = []*Card{card1, card2}
	err = deckSource.SyncDeck(initialDeck)
	if err != nil {
		t.Fatalf("failed to do initial deck sync: %s", err)
	}

	// get map of card ID to file modification times
	IDToOldModTime := map[string]time.Time{}
	dirEntries, err := os.ReadDir(filepath.Join(tempDir, testDeckName))
	if err != nil {
		t.Fatalf("failed to read test deck dir: %s", err)
	}
	for _, dirEntry := range dirEntries {
		id := strings.Split(dirEntry.Name(), ".")[0]
		info, err := dirEntry.Info()
		if err != nil {
			t.Fatalf("failed to get file info for card %s: %s", id, err)
		}
		modTime := info.ModTime()
		IDToOldModTime[id] = modTime
	}

	// read the deck from the tempdir
	deck, err := deckSource.LoadDeck(testDeckName)
	if err != nil {
		t.Fatalf("failed to load deck: %s", err)
	}
	if len(deck.Cards) < 2 {
		t.Fatal("failed to read the expected number of cards")
	}
	t.Logf("%#v", deck.Cards[0])
	t.Logf("%#v", deck.Cards[1])

	// modify one card and sync the deck; note that we need to wait
	// for a short time before writing otherwise the file modification
	// times will be the same
	time.Sleep(10 * time.Millisecond)
	deck.Cards[0].Answer = "here is a new answer"
	deck.Cards[0].Modified = true
	deckSource.SyncDeck(deck)
	if deck.Cards[0].Modified {
		t.Fatalf("card %q should not have .Modified = true anymore", deck.Cards[0].ID)
	}

	// check that only one card has a changed last modified field
	changedCardCount := 0
	newDirEntries, err := os.ReadDir(filepath.Join(tempDir, testDeckName))
	if err != nil {
		t.Fatalf("failed to read new directory entries: %s", err)
	}
	for _, newDirEntry := range newDirEntries {
		id := strings.Split(newDirEntry.Name(), ".")[0]
		newInfo, err := newDirEntry.Info()
		if err != nil {
			t.Fatalf("failed to get file info for card %s: %s", id, err)
		}
		newModTime := newInfo.ModTime()
		oldModTime, ok := IDToOldModTime[id]
		if !ok {
			t.Fatal("could not find old modified time")
		}
		t.Logf("\nold: %v\nnew: %v", oldModTime, newModTime)
		if newModTime != oldModTime {
			changedCardCount += 1
		}
	}
	if changedCardCount != 1 {
		t.Fatalf("expected 1 card changed but got %d", changedCardCount)
	}
}
