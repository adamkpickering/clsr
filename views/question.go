package views

import (
	"fmt"

	"github.com/adamkpickering/clsr/models"
	"github.com/gdamore/tcell/v2"
)

type ViewState interface {
	HandleEvent(event tcell.Event) ViewState
}

type QuestionViewState struct {
	card *models.Card
}

func NewQuestionViewState(card *models.Card) QuestionViewState {
	return QuestionViewState{
		card: card,
	}
}

func (vs QuestionViewState) HandleEvent(event tcell.Event) ViewState {
	switch switchedEvent := event.(type) {
	case *tcell.EventKey:
		return vs.handleKey(switchedEvent)
	}

	return vs
}

func (vs QuestionViewState) handleKey(event *tcell.EventKey) ViewState {
	switch event.Key() {
	case tcell.KeyEscape:
		return nil
	}
	//if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
	//}
	fmt.Println("asdf")
	return nil
}
