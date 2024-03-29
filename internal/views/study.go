package views

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/adamkpickering/clsr/internal/models"
	"github.com/adamkpickering/clsr/internal/scheduler"
	"github.com/adamkpickering/clsr/internal/utils"
	"github.com/gdamore/tcell/v2"
)

var ErrExit error = errors.New("exit study session")
var ErrEdit error = errors.New("edit card")

var StyleDefault tcell.Style

type studyState int

const (
	questionState studyState = iota
	questionAndAnswerState
)

const (
	failedKey rune = '1'
	hardKey   rune = '2'
	normalKey rune = '3'
	easyKey   rune = '4'
)

var keyToReviewResult = map[rune]models.ReviewResult{
	failedKey: models.Failed,
	hardKey:   models.Hard,
	normalKey: models.Normal,
	easyKey:   models.Easy,
}

type StudySession struct {
	Screen    tcell.Screen
	Cards     []*models.Card
	Scheduler scheduler.Scheduler
}

// Runs a study session. If err is ErrExit, the user requested to quit
// the study session and save changes to cards. If err is ErrEdit, the
// user wants to edit the card whose ID is included in the first argument.
func (ss StudySession) Run() (string, error) {
	totalCards := len(ss.Cards)
	for i, card := range ss.Cards {
		if cardID, err := ss.studyCard(card, totalCards, i+1); err != nil {
			return cardID, err
		}
	}
	return "", nil
}

func (ss StudySession) studyCard(card *models.Card, totalCards, cardNumber int) (string, error) {
	state := questionState
	for {
		// render screen
		ss.Screen.Clear()
		if err := ss.render(card, state, totalCards, cardNumber, ss.Scheduler); err != nil {
			return "", err
		}
		ss.Screen.Show()

		// poll for event and act on it
		eventInterface := ss.Screen.PollEvent()
		switch event := eventInterface.(type) {
		case *tcell.EventResize:
			ss.Screen.Sync()
		case *tcell.EventKey:
			key := event.Key()
			var keyRune rune
			if key == tcell.KeyRune {
				keyRune = event.Rune()
			}

			// allow user to exit cleanly and prematurely
			if key == tcell.KeyEscape || key == tcell.KeyCtrlC || keyRune == 'q' {
				return "", ErrExit
			}

			// handle different keys depending on different states
			switch state {
			case questionState:
				if keyRune == 'i' {
					card.Active = false
					card.Modified = true
					return "", nil
				} else if keyRune == 'e' {
					return card.ID, ErrEdit
				} else if key == tcell.KeyEnter || keyRune == ' ' {
					state = questionAndAnswerState
				}
			case questionAndAnswerState:
				if keyRune == 'i' {
					card.Active = false
					card.Modified = true
					return "", nil
				} else if keyRune == 'e' {
					return card.ID, ErrEdit
				}
				reviewResult, ok := keyToReviewResult[keyRune]
				if !ok {
					continue
				}
				newReview := models.NewReview(reviewResult)
				card.Reviews = append(models.ReviewSlice{newReview}, card.Reviews...)
				card.Modified = true
				return "", nil
			}
		}
	}
}

func (ss StudySession) processString(rawString string) []string {
	trimmedString := strings.TrimSpace(rawString)
	stringLines := strings.Split(trimmedString, "\n")
	return stringLines
}

func (ss StudySession) render(card *models.Card, state studyState, totalCards, cardNumber int, scheduler scheduler.Scheduler) error {
	var lines []string

	// add the status line
	statusFmtString := " Card %d/%d\t\t\tDeck: %s\t\t\tID: %s"
	statusLine := fmt.Sprintf(statusFmtString, cardNumber, totalCards, card.Deck, card.ID)
	lines = append(lines, statusLine)
	lines = append(lines, "")

	// add question, divider and (maybe) answer
	for _, questionLine := range ss.processString(card.Question) {
		lines = append(lines, " "+questionLine)
	}
	lines = append(lines, "\n")
	lines = append(lines, " ------")
	lines = append(lines, "\n")
	switch state {
	case questionState:
		for i := 0; i < len(ss.processString(card.Answer)); i++ {
			lines = append(lines, "")
		}
	case questionAndAnswerState:
		for _, answerLine := range ss.processString(card.Answer) {
			lines = append(lines, " "+answerLine)
		}
	}

	// add controls lines
	lines = append(lines, "")
	lines = append(lines, "")
	switch state {
	case questionState:
		lines = append(lines, " <space>/<enter>: show answer")
	case questionAndAnswerState:
		failed, err := getReadableDurationForResult(models.Failed, card, scheduler)
		if err != nil {
			return fmt.Errorf("failed to get readable duration for Failed: %w", err)
		}
		hard, err := getReadableDurationForResult(models.Hard, card, scheduler)
		if err != nil {
			return fmt.Errorf("failed to get readable duration for Hard: %w", err)
		}
		normal, err := getReadableDurationForResult(models.Normal, card, scheduler)
		if err != nil {
			return fmt.Errorf("failed to get readable duration for Normal: %w", err)
		}
		easy, err := getReadableDurationForResult(models.Easy, card, scheduler)
		if err != nil {
			return fmt.Errorf("failed to get readable duration for Easy: %w", err)
		}
		keyLineFmt := " <%c>: failed (%s)\t\t <%c>: hard (%s)\t\t" +
			"<%c>: normal (%s)\t\t<%c>: easy (%s)"
		keyLine := fmt.Sprintf(
			keyLineFmt,
			failedKey, failed,
			hardKey, hard,
			normalKey, normal,
			easyKey, easy,
		)
		lines = append(lines, keyLine)
	}
	lines = append(lines, " <e>: edit card")
	lines = append(lines, " <i>: set card to inactive")
	lines = append(lines, " <ctrl-C>/<escape>/<q>: save studied cards & exit")

	// print to screen
	if _, height := ss.Screen.Size(); len(lines) > height {
		return errors.New("screen is too small")
	}
	for lineIndex, line := range lines {
		for i, runeValue := range line {
			ss.Screen.SetContent(i, lineIndex, runeValue, nil, StyleDefault)
		}
	}
	return nil
}

func getReadableDurationForResult(result models.ReviewResult, card *models.Card, scheduler scheduler.Scheduler) (string, error) {
	nextReview, err := getHypotheticalNextReview(result, card, scheduler)
	if err != nil {
		return "", fmt.Errorf("failed to get hypothetical next review: %w", err)
	}
	return utils.GetReadableTimeDifference(time.Now(), nextReview), nil
}

func getHypotheticalNextReview(result models.ReviewResult, card *models.Card, scheduler scheduler.Scheduler) (time.Time, error) {
	newCard := card.Copy()
	newReview := models.NewReview(result)
	newCard.Reviews = append(models.ReviewSlice{newReview}, newCard.Reviews...)
	return scheduler.GetNextReview(newCard)
}
