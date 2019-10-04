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
	Duration        string
	ThumbnailURL    string
	ViewCount       uint
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
		Duration     string
		ThumbnailURL string
		ViewCount    uint
		CreatedAt    time.Time
	}{
		ID:           v.ID,
		PlatformID:   v.PlatformID,
		Title:        v.Title,
		Duration:     v.Duration,
		ThumbnailURL: v.ThumbnailURL,
		ViewCount:    v.ViewCount,
		CreatedAt:    v.CreatedAt,
	})
}
