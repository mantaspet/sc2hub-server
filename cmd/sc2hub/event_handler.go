package main

import (
	"fmt"
	"github.com/mantaspet/sc2hub-server/pkg/crawlers"
	"net/http"
	"strconv"
	"time"
)

func (app *application) getEvents(w http.ResponseWriter, r *http.Request) {
	allowedDayDiff := 90
	dateFromStr := r.URL.Query().Get("date_from")
	dateToStr := r.URL.Query().Get("date_to")
	dateFrom, err := time.Parse("2006-01-02", dateFromStr)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "Wrong date format. Must be yyyy-mm-dd")
		return
	}
	dateTo, err := time.Parse("2006-01-02", dateToStr)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "Wrong date format. Must be yyyy-mm-dd")
		return
	}
	dayDiff := dateTo.Sub(dateFrom).Hours() / 24
	if int(dayDiff) > allowedDayDiff {
		respondWithJSON(w, http.StatusBadRequest, "Max allowed date range is "+strconv.Itoa(allowedDayDiff)+" days")
		return
	}
	events, err := app.events.Select(dateFromStr, dateToStr)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	events, err = app.eventCategories.LoadCategories(events)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	events, err = app.eventCategories.AssignCategories(events)
	respondWithJSON(w, http.StatusOK, events)
}

func (app *application) crawlEvents(w http.ResponseWriter, r *http.Request) {
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
	events, err := crawlers.TeamliquidEvents(year, month)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	events, err = app.eventCategories.AssignCategories(events)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(events) == 0 {
		respondWithJSON(w, http.StatusOK, "No events found")
		return
	}
	rowCnt, err := app.events.InsertEvents(events)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	rowCntStr := strconv.Itoa(int(rowCnt))
	res := "Rows affected: " + rowCntStr
	respondWithJSON(w, http.StatusOK, res)
}
