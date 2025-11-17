package dto

import "time"

type BuisnessTripDTO struct {
	ID          uint
	Destination string
	StartAt     time.Time
	EndAt       time.Time
}
