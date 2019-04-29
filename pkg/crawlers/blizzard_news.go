package crawlers

import (
	"github.com/gocolly/colly"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"strconv"
	"strings"
	"time"
)

func BlizzardNews() ([]models.Article, error) {
	var articles []models.Article
	url := "https://news.blizzard.com/en-us/starcraft2"
	c := colly.NewCollector()

	c.OnHTML(".ArticleList .ArticleListItem", func(el *colly.HTMLElement) {
		var article models.Article
		article.Source = "Blizzard"
		article.URL = "https://news.blizzard.com" + el.ChildAttr("a", "href")
		imageStyle := el.ChildAttr(".ArticleListItem-image", "style")
		article.ThumbnailURL = "https://" + imageStyle[23:len(imageStyle)-1]
		article.Title = el.ChildText(".ArticleListItem-title")
		article.Excerpt = el.ChildText(".ArticleListItem-description .h6")
		timestamp := el.ChildText(".ArticleListItem-footerTimestamp")

		var publishedAt time.Time
		var err error
		if strings.Contains(timestamp, " days ago") {
			publishedAt, err = turnTimestampIntoDate(timestamp, 9, "days")
		} else if strings.Contains(timestamp, " day ago") {
			publishedAt, err = turnTimestampIntoDate(timestamp, 8, "hours")
		} else if strings.Contains(timestamp, " hours ago") {
			publishedAt, err = turnTimestampIntoDate(timestamp, 10, "hours")
		} else if strings.Contains(timestamp, " hour ago") {
			publishedAt, err = turnTimestampIntoDate(timestamp, 9, "hours")
		} else {
			publishedAt, err = time.Parse("January 2, 2006", timestamp)
		}
		if err == nil {
			article.PublishedAt = publishedAt
		}

		articles = append(articles, article)
	})

	err := c.Visit(url)
	return articles, err
}

func turnTimestampIntoDate(timestamp string, offset int, unit string) (time.Time, error) {
	unitsAgo := timestamp[:len(timestamp)-offset]
	unitCount, err := strconv.Atoi(unitsAgo)
	if err != nil {
		return time.Now(), err
	} else {
		if unit == "days" {
			return time.Now().AddDate(0, 0, -unitCount), err
		} else {
			return time.Now().Add(-time.Duration(unitCount) * time.Hour), err
		}
	}
}
