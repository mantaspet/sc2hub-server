package main

import (
	"github.com/mantaspet/sc2hub-server/api"
	"github.com/mantaspet/sc2hub-server/database"
	"log"
	"net/http"
)

func main() {
	database.InitDatabase()
	defer database.Close()
	router := api.InitRouter()
	log.Fatal(http.ListenAndServe(":9000", router))
}
