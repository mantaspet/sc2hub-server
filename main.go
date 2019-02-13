package main

import (
	"encoding/json"
	"fmt"
	"github.com/mantaspet/sc2hub-server/crawler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "sc2hub api\nAvailable endpoints:\nGET /events")
		if err != nil {
			log.Fatal("Unable to return response ", err)
		}
	})

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		events := crawler.TeamliquidEvents()
		eventsJson, err := json.Marshal(events)
		if err != nil {
			log.Fatal("Cannot encode to JSON ", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(eventsJson)
		if err != nil {
			log.Fatal("Unable to return response ", err)
		}
	})

	log.Fatal(http.ListenAndServe(":9000", nil))
}
