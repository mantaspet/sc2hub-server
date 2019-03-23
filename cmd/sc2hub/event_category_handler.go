package main

import (
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) getEventCategories(w http.ResponseWriter, r *http.Request) {
	eventCategories, err := app.eventCategories.SelectEventCategories()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, eventCategories)
}

func (app *application) createEventCategory(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var ec models.EventCategory
	err := decoder.Decode(&ec)
	if err != nil {
		respondWithJSON(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	validation := ec.Validate()
	if len(validation) > 0 {
		respondWithJSON(w, http.StatusUnprocessableEntity, validation)
		return
	}
	res, err := app.eventCategories.InsertEventCategory(ec)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, res)
}

func (app *application) updateEventCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	decoder := json.NewDecoder(r.Body)
	var ec models.EventCategory
	err := decoder.Decode(&ec)
	if err != nil {
		respondWithJSON(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	ec.ID, err = strconv.Atoi(id)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	validation := ec.Validate()
	if len(validation) > 0 {
		respondWithJSON(w, http.StatusUnprocessableEntity, validation)
		return
	}
	res, err := app.eventCategories.UpdateEventCategory(id, ec)
	if err == sql.ErrNoRows {
		respondWithJSON(w, http.StatusNotFound, "Event category with a specified ID does not exist")
		return
	} else if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, res)
}

func (app *application) deleteEventCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := app.eventCategories.DeleteEventCategory(id)
	if err == sql.ErrNoRows {
		respondWithJSON(w, http.StatusNotFound, "Event category with a specified ID does not exist")
		return
	} else if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint fails") {
			respondWithJSON(w, http.StatusConflict, err.Error())
		} else {
			respondWithJSON(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, "Event category was deleted")
}

func (app *application) reorderEventCategories(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		ID       int
		Priority int
	}
	var b reqBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&b)
	if b.Priority == 0 || b.ID == 0 {
		respondWithJSON(w, http.StatusBadRequest, "Request body must contain ID and Priority fields")
		return
	}
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	err = app.eventCategories.UpdateEventCategoryPriorities(b.ID, b.Priority)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, "Event category priorities were updated")
}

func eventCategoryPreflight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
}
