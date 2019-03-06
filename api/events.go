package api

import (
	"fmt"
	"github.com/mantaspet/sc2hub-server/crawlers"
	"github.com/mantaspet/sc2hub-server/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetEvents(w http.ResponseWriter, r *http.Request) {
	year := r.URL.Query().Get("year")
	month := r.URL.Query().Get("month")
	var event models.Event
	events := []models.Event{}
	if len(month) != 2 {
		month = time.Now().UTC().Format("01")
	}
	if len(year) != 4 {
		year = time.Now().UTC().Format("2006")
	}
	dateFromStr := year + month + "01"
	dateFrom, _ := time.Parse("20060102", dateFromStr)
	dateTo := dateFrom.AddDate(0, 1, 0)
	rows, err := DB.Query("SELECT id, COALESCE(event_category_id, 0) as event_category_id, COALESCE(team_liquid_id, 0) as team_liquid_id, title, stage, starts_at, info FROM events WHERE starts_at BETWEEN ? AND ?", dateFrom, dateTo)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&event.ID, &event.EventCategoryID, &event.TeamLiquidID, &event.Title, &event.Stage, &event.StartsAt, &event.Info)
		if err != nil {
			panic(err.Error())
		}
		events = append(events, event)
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
	year := r.URL.Query().Get("year")
	month := r.URL.Query().Get("month")
	if len(month) != 2 {
		month = time.Now().UTC().Format("01")
	}
	if len(year) != 4 {
		year = time.Now().UTC().Format("2006")
	}
	events := crawlers.TeamliquidEvents(year, month)
	eventCnt := len(events)
	if eventCnt == 0 {
		respondWithJSON(w, http.StatusOK, "No events found")
		return
	}
	valueStrings := make([]string, 0, len(events))
	valueArgs := make([]interface{}, 0, len(events)*4)
	for _, e := range events {
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, e.Title)
		valueArgs = append(valueArgs, e.TeamLiquidID)
		valueArgs = append(valueArgs, e.Stage)
		valueArgs = append(valueArgs, e.StartsAt)
	}
	q := fmt.Sprintf("INSERT INTO events (title, team_liquid_id, stage, starts_at) VALUES %s "+
		"ON DUPLICATE KEY UPDATE title=VALUES(title), stage=VALUES(stage), starts_at=VALUES(starts_at);",
		strings.Join(valueStrings, ","))
	res, err := DB.Exec(q, valueArgs...)

	log.Println(q)
	if err != nil {
		panic(err.Error())
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	rowCntStr := strconv.Itoa(int(rowCnt))
	response := "Rows affected: " + rowCntStr
	_, _ = DB.Exec("ALTER TABLE events AUTO_INCREMENT=1") // to prevent ON DUPLICATE KEY triggers from inflating next ID
	respondWithJSON(w, http.StatusOK, response)
}
