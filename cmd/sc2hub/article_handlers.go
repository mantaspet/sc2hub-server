package main

import (
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

func (app *application) getArticlesByCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	res, err := app.articles.SelectByCategory(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.json(w, res)
}
