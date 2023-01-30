package models

import (
	"time"
)

type Review struct {
	Version  int       `json:"version"`
	Result   int       `json:"result"`
	DateTime time.Time `json:"datetime"`
}