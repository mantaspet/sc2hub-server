package api

import (
	"github.com/go-chi/chi"
	"github.com/mantaspet/sc2hub-server/crawlers"
	"github.com/mantaspet/sc2hub-server/models"
	"log"
	"net/http"
)

func GetEvents(w http.ResponseWriter, r *http.Request) {
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
	respondWithJSON(w, http.StatusOK, events)
}

func CrawlEvents(w http.ResponseWriter, r *http.Request) {
	year := chi.URLParam(r, "year")
	month := chi.URLParam(r, "month")
	events := crawlers.TeamliquidEvents(year, month)
	respondWithJSON(w, http.StatusOK, events)
}
