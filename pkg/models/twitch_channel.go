package models

type TwitchChannel struct {
	ID              int
	EventCategoryID int
	TwitchUserID    int
	Login           string
	DisplayName     string
	ProfileImageURL string
}
