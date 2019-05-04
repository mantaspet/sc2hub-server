package models

import "time"

var ArticlePageLength = 24

type Article struct {
	ID           int
	Title        string
	Source       string
	PublishedAt  time.Time
	Excerpt      string
	ThumbnailURL string
	URL          string
}

type PaginatedArticles struct {
	Items  []*Article
	Cursor int
}
