package api

import (
	"github.com/mantaspet/sc2hub-server/models"
	"log"
)

func GetEvents() []models.Event {
	var event models.Event
	var events []models.Event
	rows, err := DB.Query("SELECT * FROM events")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&event.ID, &event.EventCategoryID, &event.Title, &event.Stage, &event.StartsAt, &event.Info)
		if err != nil {
			log.Fatal(err)
		}
		events = append(events, event)
		//log.Println(event)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return events
}
