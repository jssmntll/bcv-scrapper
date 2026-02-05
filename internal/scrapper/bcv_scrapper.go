package scrapper

import (
	"crypto/tls"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/jssmntll/bcv-scrapper/internal/store"
)

func ScrapeBCV() {
	// create a new collector with allowed domains
	c := colly.NewCollector(
		colly.AllowedDomains("www.bcv.org.ve"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	c.SetRequestTimeout(30 * time.Second)

	var tempUSD, tempEUR, tempDate string

	c.OnHTML("#dolar strong", func(e *colly.HTMLElement) {
		tempUSD = strings.TrimSpace(e.Text)
	})

	c.OnHTML("#euro strong", func(e *colly.HTMLElement) {
		tempEUR = strings.TrimSpace(e.Text)
	})

	c.OnHTML("div.pull-right.dinpro.center span.date-display-single", func(e *colly.HTMLElement) {
		rawDate := strings.TrimSpace(e.Text)
		formattedDate, err := ParseBCVDate(rawDate)
		if err != nil {
			log.Printf("Failed to parse date [%s]: %v", rawDate, err)
			tempDate = rawDate // backup if it fails
		} else {
			tempDate = formattedDate // will be "2026-02-05"
		}
	})

	err := c.Visit("https://www.bcv.org.ve/")
	if err != nil {
		log.Printf("Error scraping: %v", err)
		return
	}

	// update the Store in a safe way
	store.GlobalStore.Lock()
	store.GlobalStore.USD = tempUSD
	store.GlobalStore.EUR = tempEUR
	store.GlobalStore.Date = tempDate
	store.GlobalStore.UpdatedAt = time.Now()
	store.GlobalStore.Unlock()

	log.Println("Data updated successfully")
}
