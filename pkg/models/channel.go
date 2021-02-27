package models

import "encoding/json"

type Channel struct {
	ID                string
	PlatformID        int
	Login             string
	Title             string
	ProfileImageURL   string
	IsCrawlingEnabled bool

	// Not stored in channels table
	IncludePatterns string
	ExcludePatterns string
	EventCategoryID int
}

func (c *Channel) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID                string
		PlatformID        int
		Login             string
		Title             string
		ProfileImageURL   string
		IsCrawlingEnabled bool
	}{
		ID:                c.ID,
		PlatformID:        c.PlatformID,
		Login:             c.Login,
		Title:             c.Title,
		ProfileImageURL:   c.ProfileImageURL,
		IsCrawlingEnabled: c.IsCrawlingEnabled,
	})
}
