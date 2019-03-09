package api

import (
	"github.com/go-chi/chi"
	"net/http"
)

func GetEventCategories(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, "get events")
}

func CreateEventCategory(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, "create event")
}

func UpdateEventCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	respondWithJSON(w, http.StatusOK, "update event "+id)
}

func DeleteEventCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	respondWithJSON(w, http.StatusOK, "delete event "+id)
}
