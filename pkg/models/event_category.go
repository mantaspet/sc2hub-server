package models

import (
	"database/sql"
	"encoding/json"
	"github.com/mantaspet/sc2hub-server/pkg/validators"
)

type EventCategory struct {
	ID          int
	Name        string
	Pattern     string
	InfoURL     string
	ImageURL    string
	Description string
	Priority    int
}

type EventCategoryArticle struct {
	EventCategoryID int
	ArticleID       int
}

func (ec *EventCategory) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID          int
		Name        string
		Pattern     string
		InfoURL     string
		ImageURL    string
		Description string
	}{
		ID:          ec.ID,
		Name:        ec.Name,
		Pattern:     ec.Pattern,
		InfoURL:     ec.InfoURL,
		ImageURL:    ec.ImageURL,
		Description: ec.Description,
	})
}

func (ec EventCategory) Validate(db *sql.DB) map[string]string {
	errors := make(map[string]string)
	validators.SetError(errors, "Name",
		validators.Required(ec.Name),
		validators.SQLUnique(db, "event_categories", "name", ec.Name, ec.ID),
	)
	validators.SetError(errors, "Pattern",
		validators.Required(ec.Pattern),
		validators.MaxLength(ec.Pattern, 20),
		validators.SQLUnique(db, "event_categories", "pattern", ec.Pattern, ec.ID),
	)
	return validators.Errors(errors)
}
