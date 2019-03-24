package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	app.logTrace(err)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
	_ = app.errorLog.Output(2, trace)
}
