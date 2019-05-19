package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (app *application) router() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Route("/articles", func(r chi.Router) {
		r.Get("/", app.getAllArticles)
		r.Get("/crawl", app.isAuthenticated(app.initArticleCrawler))
	})

	r.Route("/channels", func(r chi.Router) {
		r.Get("/twitch", app.getAllTwitchChannels)
	})

	r.Route("/events", func(r chi.Router) {
		r.Get("/", app.getEvents)
		r.Get("/{id}", app.getEvent)
		r.Get("/crawl", app.isAuthenticated(app.initEventCrawler))
	})

	r.Route("/event-categories", func(r chi.Router) {
		r.Get("/", app.getEventCategories)
		r.Get("/{id}", app.getEventCategory)
		r.Get("/{id}/videos", app.getVideosByCategory)
		r.Get("/{id}/articles", app.getArticlesByCategory)
		r.Get("/{id}/broadcasts", app.getEventBroadcasts)
		r.Get("/{id}/articles", app.getArticlesByCategory)
		r.Get("/{id}/channels", app.getChannelsByCategory)
		r.Post("/{id}/channels", app.isAuthenticated(app.addChannelToCategory))
		r.Delete("/{categoryID}/channels/{channelID}", app.isAuthenticated(app.deleteCategoryChannel))
		r.Post("/", app.isAuthenticated(app.createEventCategory))
		r.Put("/{id}", app.isAuthenticated(app.updateEventCategory))
		r.Put("/reorder", app.isAuthenticated(app.reorderEventCategories))
		r.Delete("/{id}", app.isAuthenticated(app.deleteEventCategory))
		r.Options("/*", app.genericPreflightHandler)
	})

	r.Route("/players", func(r chi.Router) {
		r.Get("/", app.getAllPlayers)
		r.Get("/{id}", app.getPlayer)
		r.Get("/crawl", app.isAuthenticated(app.initPlayerCrawler))
		r.Get("/{id}/videos", app.getVideosByPlayer)
		r.Get("/{id}/articles", app.getArticlesByPlayer)
	})

	r.Route("/videos", func(r chi.Router) {
		r.Get("/", app.getAllVideos)
		r.Get("/query-apis", app.isAuthenticated(app.initVideoQuerying))
	})

	r.Route("/twitch", func(r chi.Router) {
		r.Get("/app-access-token", app.getTwitchAppAccessToken)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/token", app.getAccessToken)
		r.Options("/*", app.genericPreflightHandler)
	})

	return r
}
