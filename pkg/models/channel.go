package models

import "encoding/json"

type Channel struct {
	ID              string
	PlatformID      int
	Login           string
	Title           string
	ProfileImageURL string

	// Not stored in channels table
	Pattern         string
	EventCategoryID int
}

func (c *Channel) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID              string
		PlatformID      int
		Login           string
		Title           string
		ProfileImageURL string
	}{
		ID:              c.ID,
		PlatformID:      c.PlatformID,
		Login:           c.Login,
		Title:           c.Title,
		ProfileImageURL: c.ProfileImageURL,
	})
}
