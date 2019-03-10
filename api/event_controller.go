package api

import (
	"fmt"
	"github.com/mantaspet/sc2hub-server/crawlers"
	"github.com/mantaspet/sc2hub-server/database"
	"net/http"
	"strconv"
	"time"
)

func getEvents(w http.ResponseWriter, r *http.Request) {
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
	events, err := database.SelectEvents(dateFromStr, dateToStr)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, events)
}

func crawlEvents(w http.ResponseWriter, r *http.Request) {
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
	if len(events) == 0 {
		respondWithJSON(w, http.StatusOK, "No events found")
		return
	}
	rowCnt, err := database.InsertEvents(events)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
	}
	rowCntStr := strconv.Itoa(int(rowCnt))
	res := "Rows affected: " + rowCntStr
	respondWithJSON(w, http.StatusOK, res)
}
