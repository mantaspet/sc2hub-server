package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (app *application) initRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Route("/events", func(r chi.Router) {
		r.Get("/", app.getEvents)
		r.Get("/crawl", app.crawlEvents)
	})

	r.Route("/event-categories", func(r chi.Router) {
		r.Get("/", app.getEventCategories)
		r.Post("/", app.createEventCategory)
		r.Put("/{id}", app.updateEventCategory)
		r.Delete("/{id}", app.deleteEventCategory)
		r.Put("/reorder", app.reorderEventCategories)
		r.Options("/*", eventCategoryPreflight)
	})

	return r
}
