package main

import (
	"fmt"
	"github.com/mantaspet/sc2hub-server/pkg/crawlers"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func getPaginatedArticlesResponse(articles []*models.Article, cursor int) models.PaginatedArticles {
	var res models.PaginatedArticles
	itemCount := len(articles)
	if itemCount < models.ArticlePageLength+1 {
		res = models.PaginatedArticles{
			Cursor: 0,
			Items:  articles,
		}
	} else {
		res = models.PaginatedArticles{
			Cursor: cursor,
			Items:  articles[:itemCount-1],
		}
	}
	return res
}

func (app *application) getAllArticles(w http.ResponseWriter, r *http.Request) {
	var articles []*models.Article
	var err error
	from := app.parsePaginationParam(r.URL.Query().Get("from"))
	if r.URL.Query().Get("recent") != "" {
		articles, err = app.articles.SelectPage(9, 0, "")
	} else {
		articles, err = app.articles.SelectPage(models.ArticlePageLength, from, r.URL.Query().Get("query"))
	}
	if err != nil {
		app.serverError(w, err)
		return
	}

	res := getPaginatedArticlesResponse(articles, from+models.ArticlePageLength)
	app.json(w, res)
}

func (app *application) getArticlesByPlayer(w http.ResponseWriter, r *http.Request) {
	id, err := app.parseIDParam(w, r)
	if err != nil {
		return
	}

	from := app.parsePaginationParam(r.URL.Query().Get("from"))

	articles, err := app.articles.SelectByPlayer(models.ArticlePageLength, from, r.URL.Query().Get("query"), id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	res := getPaginatedArticlesResponse(articles, from+models.ArticlePageLength)
	app.json(w, res)
}

func (app *application) getArticlesByCategory(w http.ResponseWriter, r *http.Request) {
	id, err := app.parseIDParam(w, r)
	if err != nil {
		return
	}

	from := app.parsePaginationParam(r.URL.Query().Get("from"))

	articles, err := app.articles.SelectByCategory(models.ArticlePageLength, from, r.URL.Query().Get("query"), id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	res := getPaginatedArticlesResponse(articles, from+models.ArticlePageLength)
	app.json(w, res)
}

func (app *application) crawlArticles(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer fmt.Printf("Successfully crawled articles. Elapsed time: %v\n", time.Since(start))

	crawledArticles, err := crawlers.BlizzardNews()
	if err != nil {
		app.serverError(w, err)
		return
	}

	rowCnt, err := app.articles.InsertMany(crawledArticles)
	if err != nil {
		app.serverError(w, err)
		return
	}
	rowCntStr := strconv.Itoa(int(rowCnt))
	res := "Rows affected: " + rowCntStr

	if rowCnt == 0 {
		app.json(w, res)
		return
	}

	articles, err := app.articles.SelectLastInserted(rowCnt)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Select all player IDs for matching against crawled article titles and excerpts
	players, err := app.players.SelectAllPlayerIDs()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Select all event category patterns for matching against crawled article titles and excerpts
	ecs, err := app.eventCategories.SelectAllPatterns()
	if err != nil {
		app.serverError(w, err)
		return
	}

	var playerArticles []models.PlayerArticle
	var ecArticles []models.EventCategoryArticle
	for _, a := range articles {
		fmt.Println(a.ID)
		for _, p := range players {
			if strings.Contains(a.Title, p.PlayerID) || strings.Contains(a.Excerpt, p.PlayerID) {
				playerArticle := models.PlayerArticle{
					PlayerID:  p.ID,
					ArticleID: a.ID,
				}
				playerArticles = append(playerArticles, playerArticle)
				break
			}
		}

		for _, ec := range ecs {
			if strings.Contains(strings.ToLower(a.Title), ec.Pattern) ||
				strings.Contains(strings.ToLower(a.Excerpt), ec.Pattern) {
				ecArticle := models.EventCategoryArticle{
					EventCategoryID: ec.ID,
					ArticleID:       a.ID,
				}
				ecArticles = append(ecArticles, ecArticle)
				break
			}
		}
	}

	_, err = app.players.InsertPlayerArticles(playerArticles)
	if err != nil {
		app.serverError(w, err)
	}

	_, err = app.eventCategories.InsertEventCategoryArticles(ecArticles)
	if err != nil {
		app.serverError(w, err)
	}

	app.json(w, res)
}
