package models

import "fmt"

type Event struct {
	ID              int    `json:"id" db:"id"`
	EventCategoryID int    `json:"event_category_id" db:"event_category_id"`
	Title           string `json:"title" db:"title"`
	Stage           string `json:"stage" db:"stage"`
	StartsAt        string `json:"starts_at" db:"starts_at"`
	Info            string `json:"info" db:"info"`
}

func (e Event) String() string {
	return fmt.Sprintf("\nid: %d\nevent_category_id: %d\ntitle: %v\nstage: %v\nstarts_at: %v\ninfo: %v\n", e.ID, e.EventCategoryID, e.Title, e.Stage, e.StartsAt, e.Info)
}
