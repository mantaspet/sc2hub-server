package models

import (
	"encoding/json"
	"time"
)

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

func (v *Video) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		TwitchID  int
		Title     string
		Duration  string
		CreatedAt time.Time
	}{
		TwitchID:  v.TwitchID,
		Title:     v.Title,
		Duration:  v.Duration,
		CreatedAt: v.CreatedAt,
	})
}
