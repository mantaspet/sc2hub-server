package models

type Article struct {
	ID              int
	EventID         int
	EventCategoryID int
	SourceID        int
	Title           string
	Author          string
	PublishedAt     string
	Excerpt         string
}
