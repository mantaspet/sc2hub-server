package main

import (
	"errors"
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
		app.clientError(w, http.StatusBadRequest, errors.New("wrong date format. Must be yyyy-mm-dd"))
		return
	}
	dateTo, err := time.Parse("2006-01-02", dateToStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, errors.New("wrong date format. Must be yyyy-mm-dd"))
		return
	}

	dayDiff := dateTo.Sub(dateFrom).Hours() / 24
	if int(dayDiff) > allowedDayDiff {
		app.clientError(w, http.StatusBadRequest, errors.New("max allowed date range is "+strconv.Itoa(allowedDayDiff)+" days"))
		return
	}

	events, err := app.events.SelectInDateRange(dateFromStr, dateToStr)
	if err != nil {
		app.serverError(w, err)
		return
	}

	events, err = app.eventCategories.LoadOnEvents(events)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.json(w, events)
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
		app.serverError(w, err)
		return
	}

	events, err = app.eventCategories.AssignToEvents(events)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if len(events) == 0 {
		app.json(w, "No events found")
		return
	}

	rowCnt, err := app.events.Insert(events)
	if err != nil {
		app.serverError(w, err)
		return
	}
	rowCntStr := strconv.Itoa(int(rowCnt))
	res := "Rows affected: " + rowCntStr

	app.json(w, res)
}
