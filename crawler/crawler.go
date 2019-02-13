package crawler

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type event struct {
	title string
	stage string
	date  string
	time  string
}

func Crawl() {
	fmt.Printf("Begin crawling\n")
	events := teamliquid()
	for _, e := range events {
		fmt.Printf("%v\n", e)
	}
}

func teamliquid() []event {
	var events []event
	year, month, day := time.Now().Date()
	url := fmt.Sprintf("https://www.teamliquid.net/calendar/?view=month&year=%v&month=%v&day=%v&game=1", year, int(month), day)
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
		var monthno string
		if month < 10 {
			monthno = fmt.Sprintf("0%v", int(month))
		} else {
			monthno = fmt.Sprintf("%v", int(month))
		}
		for i := range titles {
			if len(days[i]) < 2 {
				days[i] = fmt.Sprintf("0%v", days[i])
			}
			date := fmt.Sprintf("%v-%v-%v", year, monthno, days[i])
			events = append(events, event{title: titles[i], stage: stages[i], date: date, time: times[i]})
		}
	}

	return events
}
