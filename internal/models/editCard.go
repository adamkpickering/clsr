package models

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	tempFileQuestion = "# Write the question here. This line, as well as the divider below, will be removed.\n"
	tempFileDivider  = "--------------------\n"
	tempFileAnswer   = "# Write the answer here. This line, as well as the above divider, will be removed.\n"
)

var ErrNotModified error = errors.New("temporary file not modified")

// Returns the editor specified in the EDITOR env var. If EDITOR is not specified,
// or has zero length, defaults to "nano".
func getPreferredEditor() (string, error) {
	value, ok := os.LookupEnv("EDITOR")
	if (!ok || len(value) == 0) && runtime.GOOS == "windows" {
		return "", errors.New("EDITOR environment variable is not set")
	}
	if ok && len(value) > 0 {
		return value, nil
	}
	return "nano", nil
}

// editCardViaEditor allows the user to edit a card by execing into
// their preferred editor. Any changes they make are made to the
// passed *models.Card. If the user exits without writing any changes,
// error is set to ErrNotModified.
func EditCardViaEditor(card *Card) error {
	// create temp directory
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// write temp file into the temp directory
	initialText := fmt.Sprintf("%s%s\n%s%s%s\n", tempFileQuestion, card.Question, tempFileDivider, tempFileAnswer, card.Answer)
	tempFilePath := filepath.Join(tempDir, "clsr_create_card.txt")
	err = os.WriteFile(tempFilePath, []byte(initialText), 0644)
	if err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	defer os.Remove(tempFilePath)

	// get last modified time of temp file
	info, err := os.Stat(tempFilePath)
	if err != nil {
		return fmt.Errorf("failed to get temp file info: %w", err)
	}
	firstModified := info.ModTime()

	// call the user's editor to let them edit the card
	editor, err := getPreferredEditor()
	if err != nil {
		return fmt.Errorf("failed to get editor: %w", err)
	}
	cmd := exec.Command(editor, tempFilePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("editor error: %w", err)
	}

	// return if the user did not write the temp file
	info, err = os.Stat(tempFilePath)
	if err != nil {
		return fmt.Errorf("failed to get temp file info after potential write: %w", err)
	}
	if !info.ModTime().After(firstModified) {
		return ErrNotModified
	}

	// read the contents of the temp file and parse into a Card
	contents, err := os.ReadFile(tempFilePath)
	if err != nil {
		return fmt.Errorf("failed to read temp file: %w", err)
	}
	elements := strings.Split(string(contents), tempFileDivider)
	if len(elements) != 2 {
		return fmt.Errorf(`splitting on "%s" did not produce exactly 2 elements`, tempFileDivider)
	}
	card.Question = strings.TrimSpace(strings.ReplaceAll(elements[0], tempFileQuestion, ""))
	card.Answer = strings.TrimSpace(strings.ReplaceAll(elements[1], tempFileAnswer, ""))

	return nil
}
