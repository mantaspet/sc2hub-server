package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mantaspet/sc2hub-server/api"
	"log"
	"net/http"
)

func eventRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", api.GetEvents)
	r.Get("/crawl", api.CrawlEvents)
	return r
}

func eventCategoryRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", api.GetEventCategories)
	r.Post("/", api.CreateEventCategory)
	r.Put("/{id}", api.UpdateEventCategory)
	r.Delete("/{id}", api.DeleteEventCategory)
	return r
}

func main() {
	api.InitDatabase()
	defer api.DB.Close()

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Mount("/events", eventRouter())
	r.Mount("/event-categories", eventCategoryRouter())

	log.Fatal(http.ListenAndServe(":9000", r))
}
