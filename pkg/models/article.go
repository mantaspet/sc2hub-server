package models

import "time"

type Article struct {
	ID           int
	Title        string
	Source       string
	PublishedAt  time.Time
	Excerpt      string
	ThumbnailURL string
	URL          string
}
