package models

import (
	"encoding/json"
	"fmt"
)

type Event struct {
	ID              int
	EventCategoryID int
	TeamLiquidID    int
	Title           string
	Stage           string
	StartsAt        string
}

func (e *Event) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID              int
		EventCategoryID int
		Title           string
		Stage           string
		StartsAt        string
	}{
		ID:              e.ID,
		EventCategoryID: e.EventCategoryID,
		Title:           e.Title,
		Stage:           e.Stage,
		StartsAt:        e.StartsAt,
	})
}

func (e Event) String() string {
	return fmt.Sprintf("\nid: %d\nevent_category_id: %d\nteam_liquid_id: %d\ntitle: %v\nstage: %v\nstarts_at: %v\n", e.ID, e.EventCategoryID, e.TeamLiquidID, e.Title, e.Stage, e.StartsAt)
}
