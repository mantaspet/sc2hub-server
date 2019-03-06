package api

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/mantaspet/sc2hub-server/crawlers"
	"github.com/mantaspet/sc2hub-server/models"
	"log"
	"net/http"
	"strings"
	"time"
)

func GetEvents(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	var events []models.Event
	rows, err := DB.Query("SELECT id, COALESCE(event_category_id, 0) as event_category_id, title, stage, starts_at, info FROM events")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&event.ID, &event.EventCategoryID, &event.Title, &event.Stage, &event.StartsAt, &event.Info)
		if err != nil {
			panic(err.Error())
		}
		events = append(events, event)
		//log.Println(event)
	}
	err = rows.Err()
	if err != nil {
		panic(err.Error())
	}
	respondWithJSON(w, http.StatusOK, events)
}

func CrawlEvents(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer fmt.Printf("Successfully crawled teamliquid.net events. Elapsed time: %v\n", time.Since(start))
	year := chi.URLParam(r, "year")
	month := chi.URLParam(r, "month")
	events := crawlers.TeamliquidEvents(year, month)
	var query strings.Builder
	if len(events) == 0 {
		return
	}
	query.WriteString("INSERT INTO events (title, stage, starts_at) VALUES ")
	for _, e := range events {
		fmt.Fprintf(&query, "(\"%v\", \"%v\", \"%v\"), ", e.Title, e.Stage, e.StartsAt)
	}
	q := query.String()
	if strings.HasSuffix(q, ", ") {
		q = q[:len(q)-2]
	}
	log.Println(q)
	res, err := DB.Exec(q)
	if err != nil {
		panic(err.Error())
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	response := "Rows inserted: " + string(rowCnt)
	respondWithJSON(w, http.StatusOK, response)
}
