package main

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/mantaspet/sc2hub-server/pkg/crawlers"
	"github.com/mantaspet/sc2hub-server/pkg/models"
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

	app.json(w, events)
}

func (app *application) getEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, err := strconv.Atoi(id); err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	res, err := app.events.SelectOne(id)
	if err == models.ErrNotFound {
		app.clientError(w, http.StatusNotFound, errors.New("event with a specified ID does not exist"))
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.json(w, res)
}

func (app *application) initEventCrawler(w http.ResponseWriter, r *http.Request) {
	year := r.URL.Query().Get("year")
	month := r.URL.Query().Get("month")
	if len(month) != 2 {
		month = time.Now().Format("01")
	}
	if len(year) != 4 {
		year = time.Now().Format("2006")
	}

	res, err := app.crawlEvents(year, month)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.json(w, res)
}

func (app *application) crawlEvents(year string, month string) (string, error) {
	start := time.Now()
	defer app.infoLog.Printf("Crawled teamliquid.net events. Elapsed time: %v\n", time.Since(start))

	events, err := crawlers.TeamliquidEvents(year, month)
	if err != nil {
		return "", err
	}

	if len(events) == 0 {
		return "No events found", nil
	}

	events, err = app.eventCategories.AssignToEvents(events)
	if err != nil {
		return "", err
	}
	if len(events) == 0 {
		return "", err
	}

	rowCnt, err := app.events.InsertMany(events)
	if err != nil {
		return "", err
	}
	rowCntStr := strconv.Itoa(int(rowCnt))
	res := "Rows affected: " + rowCntStr

	return res, nil
}
