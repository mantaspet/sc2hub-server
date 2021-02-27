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
		r.Options("/*", app.genericPreflightHandler)
		r.Get("/", app.getAllArticles)
		r.Get("/crawl", isAuthenticated(app, app.initArticleCrawler))
	})

	r.Route("/channels", func(r chi.Router) {
		r.Options("/*", app.genericPreflightHandler)
		r.Get("/live", app.getLiveChannels)
		r.Get("/live-registered", app.getLiveRegisteredChannels)
		r.Put("/{id}", isAuthenticated(app, app.updateChannel))
	})

	r.Route("/events", func(r chi.Router) {
		r.Options("/*", app.genericPreflightHandler)
		r.Get("/", app.getEvents)
		r.Get("/{id}", app.getEvent)
		r.Get("/crawl", isAuthenticated(app, app.initEventCrawler))
	})

	r.Route("/event-categories", func(r chi.Router) {
		r.Options("/*", app.genericPreflightHandler)
		r.Get("/", app.getEventCategories)
		r.Get("/{id}", app.getEventCategory)
		r.Get("/{id}/videos", app.getVideosByCategory)
		r.Get("/{id}/articles", app.getArticlesByCategory)
		r.Get("/{id}/broadcasts", app.getEventBroadcasts)
		r.Get("/{id}/articles", app.getArticlesByCategory)
		r.Get("/{id}/channels", app.getChannelsByCategory)
		r.Post("/{id}/channels", isAuthenticated(app, app.addChannelToCategory))
		r.Delete("/{categoryID}/channels/{channelID}", isAuthenticated(app, app.deleteCategoryChannel))
		r.Post("/", isAuthenticated(app, app.createEventCategory))
		r.Put("/{id}", isAuthenticated(app, app.updateEventCategory))
		r.Put("/reorder", isAuthenticated(app, app.reorderEventCategories))
		r.Delete("/{id}", isAuthenticated(app, app.deleteEventCategory))
		r.Options("/*", app.genericPreflightHandler)
	})

	r.Route("/players", func(r chi.Router) {
		r.Options("/*", app.genericPreflightHandler)
		r.Get("/", app.getAllPlayers)
		r.Get("/{id}", app.getPlayer)
		r.Get("/ids", app.getAllPlayerIDs)
		r.Get("/crawl", isAuthenticated(app, app.initPlayerCrawler))
		r.Get("/{id}/videos", app.getVideosByPlayer)
		r.Get("/{id}/articles", app.getArticlesByPlayer)
		r.Options("/*", app.genericPreflightHandler)
	})

	r.Route("/videos", func(r chi.Router) {
		r.Options("/*", app.genericPreflightHandler)
		r.Get("/", app.getAllVideos)
		r.Get("/query-apis", isAuthenticated(app, app.initVideoQuerying))
	})

	r.Route("/twitch", func(r chi.Router) {
		r.Options("/*", app.genericPreflightHandler)
		r.Get("/app-access-token", app.getTwitchAppAccessToken)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Options("/*", app.genericPreflightHandler)
		r.Post("/token", app.getAccessToken)
		r.Options("/*", app.genericPreflightHandler)
	})

	return r
}
