package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
)

func parsePaginationParam(param string) int {
	from, err := strconv.Atoi(param)
	if err != nil {
		from = 0
	}
	return from
}

func (app *application) parseIDParam(w http.ResponseWriter, r *http.Request) (int, error) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
	}
	return id, err
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	app.logTrace(err)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int, err error) {
	app.logTrace(err)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.Error(w, err.Error(), status)
}

func (app *application) validationError(w http.ResponseWriter, errors map[string]string) {
	res, err := json.Marshal(errors)
	if err != nil {
		app.errorLog.Println(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusUnprocessableEntity)
	_, err = w.Write(res)
	if err != nil {
		app.errorLog.Println(err.Error())
	}
}

func (app *application) json(w http.ResponseWriter, payload interface{}) {
	res, err := json.Marshal(payload)
	if err != nil {
		app.errorLog.Println(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, err = w.Write(res)
	if err != nil {
		app.errorLog.Println(err.Error())
	}
}

func (app *application) logTrace(err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = app.errorLog.Output(3, trace)
}

func (app *application) genericPreflightHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	if r.Header["Origin"][0] == app.adminOrigin {
		w.Header().Set("Access-Control-Allow-Origin", app.adminOrigin)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", app.appOrigin)
	}
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusOK)
}

func (app *application) matchesPattern(stringsToCheck []string, includePatterns string, excludePatterns string) bool {
	if len(includePatterns) == 0 {
		return true
	}
	// check the exclude patterns first. If one is found in any of the passed strings - abort
	if len(excludePatterns) > 0 {
		excludePatternsArray := strings.Split(excludePatterns, ",")
		for _, stringToCheck := range stringsToCheck {
			toCheck := strings.ToLower(stringToCheck)
			for _, excludePattern := range excludePatternsArray {
				if strings.Contains(toCheck, excludePattern) {
					return false
				}
			}
		}
	}

	// if the exclude patterns didn't find matches - start looking for include patterns
	includePatternsArray := strings.Split(includePatterns, ",")
	for _, stringToCheck := range stringsToCheck {
		toCheck := strings.ToLower(stringToCheck)
		for _, includePattern := range includePatternsArray {
			if strings.Contains(toCheck, includePattern) {
				return true
			}
		}
	}
	return false
}
