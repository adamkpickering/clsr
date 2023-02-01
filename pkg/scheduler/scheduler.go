package scheduler

import (
	"github.com/adamkpickering/clsr/pkg/models"
	"time"
)

// A scheduler tells you whether cards are due, and when they will be due.
type Scheduler interface {
	IsDue(card *models.Card) (bool, error)
	GetNextReview(card *models.Card) (time.Time, error)
}