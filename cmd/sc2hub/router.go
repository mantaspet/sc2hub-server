package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (app *application) router() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Route("/events", func(r chi.Router) {
		r.Get("/", app.getEvents)
		r.Get("/{id}", app.getEvent)
		r.Get("/{id}/videos", app.getVideosByCategory)     // TODO replace once video module is implemented
		r.Get("/{id}/articles", app.getArticlesByCategory) // TODO replace once article module is implemented
		r.Get("/crawl", app.crawlEvents)
	})

	r.Route("/event-categories", func(r chi.Router) {
		r.Get("/", app.getEventCategories)
		r.Get("/{id}", app.getEventCategory)
		r.Get("/{id}/videos", app.getVideosByCategory)
		r.Get("/{id}/articles", app.getArticlesByCategory)
		r.Post("/", app.createEventCategory)
		r.Put("/{id}", app.updateEventCategory)
		r.Put("/reorder", app.reorderEventCategories)
		r.Delete("/{id}", app.deleteEventCategory)
		r.Options("/*", app.genericPreflightHandler)
	})

	r.Route("/players", func(r chi.Router) {
		r.Get("/", app.getAllPlayers)
		r.Get("/crawl", app.crawlPlayers)
	})

	return r
}
