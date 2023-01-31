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
