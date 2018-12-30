package models

import (
	"time"
)

type Topup struct {
	ID         int
	Amount     float64
	CardNumber string
	TopupAt    time.Time
}
