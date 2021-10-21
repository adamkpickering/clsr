package views

import (
	"errors"
	"fmt"
	"strings"

	"github.com/adamkpickering/clsr/models"
	"github.com/gdamore/tcell/v2"
)

var ErrExit error = errors.New("exit study session")

var StyleDefault tcell.Style

type studyState int

const (
	questionState studyState = iota
	questionAndAnswerState
)

type StudySession struct {
	Screen tcell.Screen
	Cards  []*models.Card
}

func (ss StudySession) Run() error {
	for _, card := range ss.Cards {
		err := ss.studyCard(card)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ss StudySession) studyCard(card *models.Card) error {
	state := questionState
	for {
		// render screen
		ss.Screen.Clear()
		switch state {
		case questionState:
			ss.renderQuestionOnly(card)
		case questionAndAnswerState:
			ss.renderQuestionAndAnswer(card)
		}
		ss.Screen.Show()

		// poll for event and act on it
		eventInterface := ss.Screen.PollEvent()
		switch event := eventInterface.(type) {
		case *tcell.EventResize:
			ss.Screen.Sync()
		case *tcell.EventKey:
			key := event.Key()

			// allow user to exit cleanly and prematurely
			if key == tcell.KeyEscape || key == tcell.KeyCtrlC {
				return ErrExit
			}

			// handle different keys depending on different states
			switch state {
			case questionState:
				if key == tcell.KeyRune {
					keyRune := event.Rune()
					if keyRune == 'i' {
						card.Active = false
						return nil
					} else if key == tcell.KeyEnter || keyRune == ' ' {
						state = questionAndAnswerState
					}
				}
			case questionAndAnswerState:
				if key == tcell.KeyRune {
					keyRune := event.Rune()
					multiplier, err := getMultiplierFromRune(keyRune)
					if err != nil {
						continue
					}
					card.SetNextReview(multiplier)
					return nil
				}
			}
		}
	}
}

func (ss StudySession) processString(rawString string) []string {
	trimmedString := strings.TrimSpace(rawString)
	stringLines := strings.Split(trimmedString, "\n")
	return stringLines
}

func (ss StudySession) renderQuestionOnly(card *models.Card) {
	lines := ss.processString(card.Question)
	for lineIndex, line := range lines {
		for i, runeValue := range line {
			ss.Screen.SetContent(i, lineIndex, runeValue, nil, StyleDefault)
		}
	}
}

func (ss StudySession) renderQuestionAndAnswer(card *models.Card) {
	// build lines var, which represents lines to print to screen
	lines := []string{}
	for _, questionLine := range ss.processString(card.Question) {
		lines = append(lines, questionLine)
	}
	lines = append(lines, "\n")
	lines = append(lines, "------")
	lines = append(lines, "\n")
	for _, answerLine := range ss.processString(card.Answer) {
		lines = append(lines, answerLine)
	}

	// print to screen
	for lineIndex, line := range lines {
		for i, runeValue := range line {
			ss.Screen.SetContent(i, lineIndex, runeValue, nil, StyleDefault)
		}
	}
}

func getMultiplierFromRune(key rune) (float64, error) {
	valueMap := map[rune]float64{
		'1': 0.0,
		'2': 1.0,
		'3': 1.5,
		'4': 2.0,
	}
	multiplier, ok := valueMap[key]
	if ok {
		return multiplier, nil
	}
	return 0.0, fmt.Errorf(`invalid key "%c"`, key)
}
