package api

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/mantaspet/sc2hub-server/models"
	"net/http"
)

func GetEventCategories(w http.ResponseWriter, r *http.Request) {
	var ec models.EventCategory
	var eventCategories []models.EventCategory
	rows, err := DB.Query("SELECT id, name, pattern, COALESCE(info_url, '') as info_url, `order` FROM event_categories")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&ec.ID, &ec.Name, &ec.Pattern, &ec.InfoURL, &ec.Order)
		if err != nil {
			panic(err)
		}
		eventCategories = append(eventCategories, ec)
	}
	err = rows.Err()
	if err != nil {
		panic(err.Error())
	}
	respondWithJSON(w, http.StatusOK, eventCategories)
}

func CreateEventCategory(w http.ResponseWriter, r *http.Request) {
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
	_, sqlErr := DB.Exec("INSERT INTO event_categories (name, pattern, info_url, `order`) VALUES (?, ?, ?, ?)", ec.Name, ec.Pattern, ec.InfoURL, ec.Order)
	if sqlErr != nil {
		respondWithJSON(w, http.StatusInternalServerError, sqlErr)
		return
	}
	respondWithJSON(w, http.StatusOK, "success")
}

func UpdateEventCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
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
	res, sqlErr := DB.Exec("UPDATE event_categories SET name=?, pattern=?, info_url=?, `order`=? WHERE id=?", ec.Name, ec.Pattern, ec.InfoURL, ec.Order, id)
	if sqlErr != nil {
		respondWithJSON(w, http.StatusInternalServerError, sqlErr)
		return
	}
	rowCnt, _ := res.RowsAffected()
	if rowCnt == 0 {
		respondWithJSON(w, http.StatusNotFound, "Event category with specified ID does not exist")
		return
	}
	respondWithJSON(w, http.StatusOK, "success")
}

func DeleteEventCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	res, err := DB.Exec("DELETE FROM event_categories WHERE id=?", id)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	rowCnt, _ := res.RowsAffected()
	if rowCnt == 0 {
		respondWithJSON(w, http.StatusNotFound, "Event category with specified ID does not exist")
		return
	}
	respondWithJSON(w, http.StatusOK, "success")
}
