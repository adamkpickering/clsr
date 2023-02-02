package views

import (
	"errors"
	"fmt"
	"strings"

	"github.com/adamkpickering/clsr/pkg/models"
	"github.com/gdamore/tcell/v2"
)

var ErrExit error = errors.New("exit study session")

var StyleDefault tcell.Style

type studyState int

const (
	questionState studyState = iota
	questionAndAnswerState
)

var keyToReviewResult = map[rune]models.ReviewResult{
	'1': models.Failed,
	'2': models.Hard,
	'3': models.Normal,
	'4': models.Easy,
}

type StudySession struct {
	Screen tcell.Screen
	Cards  []*models.Card
}

func (ss StudySession) Run() error {
	totalCards := len(ss.Cards)
	for i, card := range ss.Cards {
		err := ss.studyCard(card, totalCards, i+1)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ss StudySession) studyCard(card *models.Card, totalCards, cardNumber int) error {
	state := questionState
	for {
		// render screen
		ss.Screen.Clear()
		err := ss.render(card, state, totalCards, cardNumber)
		if err != nil {
			return err
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
				return ErrExit
			}

			// handle different keys depending on different states
			switch state {
			case questionState:
				if keyRune == 'i' {
					card.Active = false
					card.Modified = true
					return nil
				} else if key == tcell.KeyEnter || keyRune == ' ' {
					state = questionAndAnswerState
				}
			case questionAndAnswerState:
				if keyRune == 'i' {
					card.Active = false
					card.Modified = true
					return nil
				}
				reviewResult, ok := keyToReviewResult[keyRune]
				if !ok {
					continue
				}
				newReview := models.NewReview(reviewResult)
				card.Reviews = append(models.ReviewSlice{newReview}, card.Reviews...)
				return nil
			}
		}
	}
}

func (ss StudySession) processString(rawString string) []string {
	trimmedString := strings.TrimSpace(rawString)
	stringLines := strings.Split(trimmedString, "\n")
	return stringLines
}

func (ss StudySession) render(card *models.Card, state studyState, totalCards, cardNumber int) error {
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
		// failedMultiplier, _ := getMultiplierFromRune(failedKey)
		// failed := card.GetMultipliedReviewInterval(failedMultiplier)
		// hardMultiplier, _ := getMultiplierFromRune(hardKey)
		// hard := card.GetMultipliedReviewInterval(hardMultiplier)
		// normalMultiplier, _ := getMultiplierFromRune(normalKey)
		// normal := card.GetMultipliedReviewInterval(normalMultiplier)
		// easyMultiplier, _ := getMultiplierFromRune(easyKey)
		// easy := card.GetMultipliedReviewInterval(easyMultiplier)
		// keyLineFmt := " <%c>: failed (%dd)\t\t <%c>: hard (%dd)\t\t" +
		// 	"<%c>: normal (%dd)\t\t<%c>: easy (%dd)"
		// keyLine := fmt.Sprintf(
		// 	keyLineFmt,
		// 	failedKey, failed,
		// 	hardKey, hard,
		// 	normalKey, normal,
		// 	easyKey, easy,
		// )
		// lines = append(lines, keyLine)
	}
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

// func getMultiplierFromRune(key rune) (float64, error) {
// 	switch key {
// 	case failedKey:
// 		return 0.0, nil
// 	case hardKey:
// 		return 1.0, nil
// 	case normalKey:
// 		return 1.5, nil
// 	case easyKey:
// 		return 2.0, nil
// 	default:
// 		return 0.0, fmt.Errorf(`invalid key "%c"`, key)
// 	}
// }
