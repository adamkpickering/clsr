package models

import (
	"time"
)

type ReviewResult string

const (
	Failed ReviewResult = "failed"
	Hard   ReviewResult = "hard"
	Normal ReviewResult = "normal"
	Easy   ReviewResult = "easy"
)

type Review struct {
	Version  int          `json:"version"`
	Result   ReviewResult `json:"result"`
	Datetime time.Time    `json:"datetime"`
}

type ReviewSlice []*Review

func NewReview(result ReviewResult) *Review {
	return &Review{
		Version:  0,
		Result:   result,
		Datetime: time.Now(),
	}
}

func (reviews ReviewSlice) Len() int {
	return len(reviews)
}

func (reviews ReviewSlice) Less(i, j int) bool {
	return reviews[i].Datetime.After(reviews[j].Datetime)
}

func (reviews ReviewSlice) Swap(i, j int) {
	reviews[i], reviews[j] = reviews[j], reviews[i]
}
