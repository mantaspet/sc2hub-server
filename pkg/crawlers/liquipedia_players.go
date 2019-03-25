package crawlers

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/mantaspet/sc2hub-server/pkg/models"
)

func LiquipediaPlayers(region string) ([]models.Player, error) {
	var players []models.Player
	url := fmt.Sprintf("https://liquipedia.net/starcraft2/Players_(%v)", region)
	c := colly.NewCollector(
		colly.AllowedDomains("liquipedia.net"),
		colly.Async(true),
	)

	c.OnHTML("table.sortable.wikitable", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, el *colly.HTMLElement) {
			var p models.Player
			if i > 0 {
				p.PlayerID = el.ChildText("td:first-child")
				p.LiquipediaURL = el.ChildAttr("td:first-child a", "href")
				p.Name = el.ChildText("td:nth-of-type(2)")
				p.Team = el.ChildText("td:nth-of-type(3) a")
				p.Country = el.ChildText("td:nth-of-type(4)")
				p.Race = el.ChildText("td:nth-of-type(5)")
				p.StreamURL = el.ChildAttr("td:nth-of-type(6) a", "href")
			}
			if p.PlayerID != "" {
				players = append(players, p)
			}
		})
	})

	err := c.Visit(url)

	c.Wait()
	return players, err
}
