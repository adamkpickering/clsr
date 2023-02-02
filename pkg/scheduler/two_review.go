package scheduler

import (
	"fmt"
	"sort"
	"time"

	"github.com/adamkpickering/clsr/pkg/config"
	"github.com/adamkpickering/clsr/pkg/models"
	"github.com/adamkpickering/clsr/pkg/utils"
)

// A scheduler that looks at the most recent two reviews
// of a card to determine when it should be reviewed next.
type TwoReviewScheduler struct {
	config *config.Config
}

func NewTwoReviewScheduler(config *config.Config) *TwoReviewScheduler {
	return &TwoReviewScheduler{
		config: config,
	}
}

func (scheduler *TwoReviewScheduler) IsDue(card *models.Card) (bool, error) {
	reviews := getSortedReviewsCopy(card)
	if len(reviews) == 0 {
		return true, nil
	}

	nextReview, err := scheduler.GetNextReview(card)
	if err != nil {
		return false, fmt.Errorf("failed to get next review: %w", err)
	}

	if utils.DatesEqual(reviews[0].Datetime, nextReview) {
		return time.Now().After(nextReview), nil
	} else {
		nextReviewYear, nextReviewMonth, nextReviewDay := nextReview.Date()
		midnightNextReview := time.Date(nextReviewYear, nextReviewMonth, nextReviewDay, 0, 0, 0, 0, nextReview.Location())
		return time.Now().After(midnightNextReview), nil
	}
}

// Returns the datetime that the card is next due.
func (scheduler *TwoReviewScheduler) GetNextReview(card *models.Card) (time.Time, error) {
	reviews := getSortedReviewsCopy(card)
	reviewsLength := len(reviews)
	if reviewsLength == 0 {
		return time.Now(), nil
	}

	if reviews[0].Result == models.Failed {
		interval := scheduler.config.FailedReviewInterval * uint(time.Hour)
		return reviews[0].Datetime.Add(time.Duration(interval)), nil
	}

	if reviewsLength == 1 || (reviewsLength == 2 && reviews[1].Result == models.Failed) {
		intervalInHours, err := getSecondReviewIntervalFor(reviews[0].Result, scheduler.config)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to get second review interval: %w", err)
		}
		interval := intervalInHours * uint(time.Hour)
		return reviews[0].Datetime.Add(time.Duration(interval)), nil
	}

	lastReview := reviews[0]
	lastLastReview := reviews[1]
	oldInterval := lastReview.Datetime.Sub(lastLastReview.Datetime)
	multiplier, err := getMultiplierFor(lastReview.Result, scheduler.config)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get interval multiplier: %w", err)
	}
	newInterval := float64(oldInterval) * multiplier
	return lastReview.Datetime.Add(time.Duration(newInterval)), nil
}

func getSecondReviewIntervalFor(result models.ReviewResult, config *config.Config) (uint, error) {
	switch result {
	case models.Hard:
		return config.SecondReviewIntervals.Hard, nil
	case models.Normal:
		return config.SecondReviewIntervals.Normal, nil
	case models.Easy:
		return config.SecondReviewIntervals.Easy, nil
	default:
		return 0.0, fmt.Errorf("got unexpected result %q", result)
	}
}

func getMultiplierFor(result models.ReviewResult, config *config.Config) (float64, error) {
	switch result {
	case models.Hard:
		return config.IntervalMultipliers.Hard, nil
	case models.Normal:
		return config.IntervalMultipliers.Normal, nil
	case models.Easy:
		return config.IntervalMultipliers.Easy, nil
	default:
		return 0.0, fmt.Errorf("got unexpected result %q", result)
	}
}

func getSortedReviewsCopy(card *models.Card) models.ReviewSlice {
	reviews := make(models.ReviewSlice, len(card.Reviews))
	copy(reviews, card.Reviews)
	sort.Stable(reviews)
	return reviews
}
