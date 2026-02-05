package main

import (
	"log"
	"net/http"

	"github.com/jssmntll/bcv-scrapper/internal/api"
	"github.com/jssmntll/bcv-scrapper/internal/scrapper"
	"github.com/robfig/cron/v3"
)

func main() {
	// execute the first scrape at startup
	go scrapper.ScrapeBCV()

	// configure cron (2 times a day: 9am and 4pm)
	c := cron.New()
	c.AddFunc("0 9,16 * * *", func() {
		log.Println("Update scheduled started...")
		scrapper.ScrapeBCV()
	})
	c.Start()

	http.HandleFunc("/api/bcv/exchanges", api.GetExchangeHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
