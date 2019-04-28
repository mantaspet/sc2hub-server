package models

import (
	"encoding/json"
	"time"
)

type Video struct {
	ID              string
	EventID         int
	EventCategoryID int
	PlatformID      int
	ChannelID       string
	Title           string
	ThumbnailURL    string
	Duration        string
	Type            string
	CreatedAt       time.Time
}

func (v *Video) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID           string
		PlatformID   int
		Title        string
		Duration     string
		ThumbnailURL string
		CreatedAt    time.Time
	}{
		ID:           v.ID,
		PlatformID:   v.PlatformID,
		Title:        v.Title,
		Duration:     v.Duration,
		ThumbnailURL: v.ThumbnailURL,
		CreatedAt:    v.CreatedAt,
	})
}
