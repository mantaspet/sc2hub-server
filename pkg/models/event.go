package models

import (
	"encoding/json"
	"fmt"
)

type Event struct {
	ID              int
	EventCategoryID int
	EventCategory   EventCategory
	TeamLiquidID    int
	Title           string
	Stage           string
	StartsAt        string
	Info            string
}

// TODO find a more elegant solution
func (e *Event) MarshalJSON() ([]byte, error) {
	if e.EventCategoryID == 0 {
		return json.Marshal(&struct {
			ID              int
			EventCategoryID interface{}
			EventCategory   interface{}
			Title           string
			Stage           string
			StartsAt        string
			Info            string
		}{
			ID:              e.ID,
			EventCategoryID: nil,
			EventCategory:   nil,
			Title:           e.Title,
			Stage:           e.Stage,
			StartsAt:        e.StartsAt,
			Info:            e.Info,
		})
	} else {
		return json.Marshal(&struct {
			ID              int
			EventCategoryID int
			EventCategory   EventCategory
			Title           string
			Stage           string
			StartsAt        string
			Info            string
		}{
			ID:              e.ID,
			EventCategoryID: e.EventCategoryID,
			EventCategory:   e.EventCategory,
			Title:           e.Title,
			Stage:           e.Stage,
			StartsAt:        e.StartsAt,
			Info:            e.Info,
		})
	}
}

func (e Event) String() string {
	return fmt.Sprintf("\nid: %d\nevent_category_id: %d\nteam_liquid_id: %d\ntitle: %v\nstage: %v\nstarts_at: %v\ninfo: %v\n", e.ID, e.EventCategoryID, e.TeamLiquidID, e.Title, e.Stage, e.StartsAt, e.Info)
}
