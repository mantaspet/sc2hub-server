package crawlers

import (
	"github.com/gocolly/colly"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"strings"
	"time"
)

func TeamLiquidNews() ([]models.Article, error) {
	var articles []models.Article
	url := "https://tl.net/news/"
	c := colly.NewCollector()

	c.OnHTML(".findex_newsblock tr", func(el *colly.HTMLElement) {
		var article models.Article
		if strings.Contains(el.ChildText("td:nth-child(2) span"), "StarCraft 2") == false {
			return
		}
		article.Source = "TeamLiquid.net"
		article.URL = "https://tl.net" + el.ChildAttr("td:nth-child(2) a", "href")
		article.ThumbnailURL = "https://i.imgur.com/e2o9EmK.jpg"
		article.Title = el.ChildText("td:nth-child(2) a")
		article.Excerpt = ""
		publishedAt, err := time.Parse("2 Jan 2006", el.ChildText("td:nth-child(3)"))
		if err == nil {
			article.PublishedAt = publishedAt
		}
		if len(article.Title) > 1 {
			lastPageSymbol := article.Title[len(article.Title)-1:]
			if lastPageSymbol == ">" {
				article.Title = article.Title[:len(article.Title)-1]
			}
			articles = append(articles, article)
		}
	})

	err := c.Visit(url)
	return articles, err
}
