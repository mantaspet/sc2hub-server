package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func InitRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Route("/events", func(r chi.Router) {
		r.Get("/", getEvents)
		r.Get("/crawl", crawlEvents)
	})

	r.Route("/event-categories", func(r chi.Router) {
		r.Get("/", getEventCategories)
		r.Post("/", createEventCategory)
		r.Put("/{id}", updateEventCategory)
		r.Delete("/{id}", deleteEventCategory)
		r.Put("/reorder", reorderEventCategories)
	})

	return r
}
