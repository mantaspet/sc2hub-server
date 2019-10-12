package models

import (
	"encoding/json"
	"time"
)

var VideoPageLength = 24

type Video struct {
	ID              string
	EventCategoryID int
	PlatformID      int
	ChannelID       string
	Title           string
	Duration        int
	ViewCount       int
	ThumbnailURL    string
	Type            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type PaginatedVideos struct {
	Items  []*Video
	Cursor int
}

func (v *Video) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID           string
		PlatformID   int
		Title        string
		Duration     int
		ViewCount    int
		ThumbnailURL string
		CreatedAt    time.Time
	}{
		ID:           v.ID,
		PlatformID:   v.PlatformID,
		Title:        v.Title,
		Duration:     v.Duration,
		ViewCount:    v.ViewCount,
		ThumbnailURL: v.ThumbnailURL,
		CreatedAt:    v.CreatedAt,
	})
}
