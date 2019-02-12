package main

import (
	"./crawler"
	"fmt"
)

func main() {
	fmt.Println("Hello from sc2hub")
	crawler.Crawl("https://www.teamliquid.net/calendar/")
}
