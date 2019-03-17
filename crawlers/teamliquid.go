package crawlers

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/mantaspet/sc2hub-server/database"
	"strconv"
)

func TeamliquidEvents(year string, month string) ([]database.Event, error) {
	var events []database.Event
	var day string
	url := fmt.Sprintf("https://www.teamliquid.net/calendar/?view=month&year=%v&month=%v&game=1", year, month)
	c := colly.NewCollector()

	c.OnHTML("td.evc-l:not(.mo_out) .ev-feed", func(e *colly.HTMLElement) {
		day = e.Attr("data-day")
		if len(day) == 1 {
			day = "0" + day
		}

		e.ForEach(".ev-block", func(i int, el *colly.HTMLElement) {
			var event database.Event
			event.Title = el.ChildText("span[data-event-id]")
			event.TeamLiquidID, _ = strconv.Atoi(el.ChildAttr("span[data-event-id]", "data-event-id"))
			event.Stage = el.ChildText(".ev-stage")
			event.StartsAt = year + "-" + month + "-" + day + " " + el.ChildText(".ev-timer")
			events = append(events, event)
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	err := c.Visit(url)
	return events, err
}
