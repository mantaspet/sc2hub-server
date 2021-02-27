package models

import (
	"database/sql"
	"encoding/json"
	"github.com/mantaspet/sc2hub-server/pkg/validators"
)

type EventCategory struct {
	ID              int
	Name            string
	IncludePatterns string
	ExcludePatterns string
	InfoURL         string
	ImageURL        string
	Description     string
	Priority        int
}

type EventCategoryArticle struct {
	EventCategoryID int
	ArticleID       int
}

func (ec *EventCategory) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID              int
		Name            string
		IncludePatterns string
		ExcludePatterns string
		InfoURL         string
		ImageURL        string
		Description     string
	}{
		ID:              ec.ID,
		Name:            ec.Name,
		IncludePatterns: ec.IncludePatterns,
		ExcludePatterns: ec.ExcludePatterns,
		InfoURL:         ec.InfoURL,
		ImageURL:        ec.ImageURL,
		Description:     ec.Description,
	})
}

func (ec EventCategory) Validate(db *sql.DB) map[string]string {
	errors := make(map[string]string)
	validators.SetError(errors, "Name",
		validators.Required(ec.Name),
		validators.MaxLength(ec.Name, 255),
		validators.SQLUnique(db, "event_categories", "name", ec.Name, ec.ID),
	)
	validators.SetError(errors, "IncludePatterns",
		validators.Required(ec.IncludePatterns),
		validators.MaxLength(ec.IncludePatterns, 255),
	)
	return validators.Errors(errors)
}
