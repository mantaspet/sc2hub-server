package models

import "time"

type Video struct {
	ID              int
	EventID         int
	EventCategoryID int
	TwitchID        int
	Title           string
	URL             string
	Duration        string
	CreatedAt       time.Time
}
