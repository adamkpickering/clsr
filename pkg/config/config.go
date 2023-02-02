package config

type Config struct {
	// The time until a card is due after a review has been failed.
	FailedReviewInterval  uint
	SecondReviewIntervals SecondReviewIntervals
	IntervalMultipliers   IntervalMultipliers
}

// This applies when the card has been reviewed exactly once,
// or when the review before the last review was failed.
// The time until the next review is the value below that corresponds
// to the Result of the one review. Values are in hours.
type SecondReviewIntervals struct {
	Hard   uint
	Normal uint
	Easy   uint
}

// This applies when a card has been reviewed 2 or more times.
// The time until the next review is calculated by multiplying
// the difference between the times of the last two reviews,
// and multiplying that difference by the below multiplier that
// corresponds to the Result on the last review.
type IntervalMultipliers struct {
	Hard   float64
	Normal float64
	Easy   float64
}

var DefaultConfig = &Config{
	FailedReviewInterval: 4,
	SecondReviewIntervals: SecondReviewIntervals{
		Hard:   4,
		Normal: 24,
		Easy:   48,
	},
	IntervalMultipliers: IntervalMultipliers{
		Hard:   1.0,
		Normal: 1.5,
		Easy:   2.0,
	},
}
