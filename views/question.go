package views

import (
	"fmt"

	"github.com/adamkpickering/clsr/models"
	"github.com/gdamore/tcell/v2"
)

type ViewState struct {
	card   *models.Card
	screen tcell.Screen
}

func NewViewState(screen tcell.Screen, card *models.Card) ViewState {
	return ViewState{
		card:   card,
		screen: screen,
	}
}

func (vs *ViewState) HandleEvent(event tcell.Event) bool {
	switch switchedEvent := event.(type) {
	case *tcell.EventResize:
		vs.screen.Sync()
	case *tcell.EventKey:
		return vs.handleKey(switchedEvent)
	}
	return false
}

func (vs *ViewState) handleKey(event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyEscape:
		return
	}
	//if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
	//}
	fmt.Println("asdf")
}

func (vs *ViewState) Draw() {
	vs.screen.Show()
}
