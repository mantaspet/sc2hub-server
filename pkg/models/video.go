package models

import "time"

type Video struct {
	ID              int
	EventID         int
	EventCategoryID int
	ChannelID       int
	TwitchID        int
	Title           string
	Duration        string
	CreatedAt       time.Time
}
