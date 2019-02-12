package main

import (
	"fmt"
	"github.com/mantaspet/sc2hub-server/crawler"
)

func main() {
	fmt.Println("Hello from sc2hub")
	crawler.Crawl("https://www.teamliquid.net/calendar/")
}
