package models

type Channel struct {
	ID              string
	EventCategoryID int
	PlatformID      int
	Login           string
	Title           string
	ProfileImageURL string
	Pattern         string
}
