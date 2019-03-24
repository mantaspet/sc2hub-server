package main

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) getEventCategories(w http.ResponseWriter, r *http.Request) {
	res, err := app.eventCategories.SelectAll()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.json(w, res)
}

func (app *application) createEventCategory(w http.ResponseWriter, r *http.Request) {
	var ec models.EventCategory
	err := json.NewDecoder(r.Body).Decode(&ec)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	errs := ec.Validate(app.db)
	if errs != nil {
		app.validationError(w, errs)
		return
	}

	res, err := app.eventCategories.Insert(ec)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.json(w, res)
}

func (app *application) updateEventCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var ec models.EventCategory
	err := json.NewDecoder(r.Body).Decode(&ec)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}
	ec.ID, err = strconv.Atoi(id)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	errs := ec.Validate(app.db)
	if len(errs) > 0 {
		app.validationError(w, errs)
		return
	}

	res, err := app.eventCategories.Update(id, ec)
	if err == models.ErrNotFound {
		app.clientError(w, http.StatusNotFound, errors.New("event category with a specified ID does not exist"))
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.json(w, res)
}

func (app *application) deleteEventCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := app.eventCategories.Delete(id)
	if err == models.ErrNotFound {
		app.clientError(w, http.StatusNotFound, errors.New("event category with a specified ID does not exist"))
		return
	} else if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint fails") {
			app.clientError(w, http.StatusConflict, err)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.json(w, "Event category was deleted")
}

func (app *application) reorderEventCategories(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		ID       int
		Priority int
	}
	var b reqBody
	err := json.NewDecoder(r.Body).Decode(&b)
	if b.Priority == 0 || b.ID == 0 {
		app.clientError(w, http.StatusBadRequest, errors.New("request body must contain ID and Priority fields"))
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.eventCategories.UpdatePriorities(b.ID, b.Priority)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.json(w, "Event category priorities were updated")
}

func eventCategoryPreflight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
}
