package api

import (
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/mantaspet/sc2hub-server/database"
	"net/http"
)

func getEventCategories(w http.ResponseWriter, r *http.Request) {
	eventCategories, err := database.SelectEventCategories()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, eventCategories)
}

func createEventCategory(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var ec database.EventCategory
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
	res, err := database.InsertEventCategory(ec)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, res)
}

func updateEventCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	decoder := json.NewDecoder(r.Body)
	var ec database.EventCategory
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
	res, err := database.UpdateEventCategory(id, ec)
	if err == sql.ErrNoRows {
		respondWithJSON(w, http.StatusNotFound, "Event category with a specified ID does not exist")
		return
	} else if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, res)
}

func deleteEventCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := database.DeleteEventCategory(id)
	if err == sql.ErrNoRows {
		respondWithJSON(w, http.StatusNotFound, "Event category with a specified ID does not exist")
		return
	} else if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, "Event category was deleted")
}

func reorderEventCategories(w http.ResponseWriter, r *http.Request) {
	var m map[int]int
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&m)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	err = database.UpdateEventCategoryPriorities(m)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, "Event category priorities were updated")
}
