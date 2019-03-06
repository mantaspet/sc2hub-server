package crawlers

import (
	"fmt"
	"github.com/mantaspet/sc2hub-server/models"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func TeamliquidEvents(year string, month string) []models.Event {
	var events []models.Event
	url := fmt.Sprintf("https://www.teamliquid.net/calendar/?view=month&year=%v&month=%v&game=1", year, month)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.Fatal(err)
	}
	var stages, titles, times, days []string
	var currentday string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "span" {
			for _, a := range n.Attr {
				if a.Key == "data-event-id" {
					titles = append(titles, n.FirstChild.Data)
					break
				}
				if a.Val == "ev-timer" && n.FirstChild != nil {
					times = append(times, n.FirstChild.Data)
					days = append(days, currentday)
					break
				}
			}
		}
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Val == "ev-stage" {
					stages = append(stages, n.FirstChild.Data)
					break
				}
				if a.Key == "data-day" {
					currentday = a.Val
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if len(titles) == len(stages) && len(times) == len(days) && len(titles) == len(days) {
		for i := range titles {
			if len(days[i]) < 2 {
				days[i] = fmt.Sprintf("0%v", days[i])
			}
			datetime := fmt.Sprintf("%v-%v-%v %v", year, month, days[i], times[i])

			events = append(events, models.Event{Title: &titles[i], Stage: &stages[i], StartsAt: &datetime})
		}
	}
	return events
}
