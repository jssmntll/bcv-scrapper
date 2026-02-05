package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/robfig/cron/v3"
)

type DataStore struct {
	sync.RWMutex
	USD       string    `json:"usd"`
	EUR       string    `json:"eur"`
	Validity  string    `json:"validity"`
	UpdatedAt time.Time `json:"last_update"`
}

var store = &DataStore{}

func main() {
	// execute the first scrape at startup
	go scrapeBCV()

	// configure cron (2 times a day: 9am and 4pm)
	c := cron.New()
	c.AddFunc("0 9,16 * * *", func() {
		log.Println("Update scheduled started...")
		scrapeBCV()
	})
	c.Start()

	http.HandleFunc("/api/bcv/exchanges", getExchangeHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func scrapeBCV() {
	// create a new collector with allowed domains
	c := colly.NewCollector(
		colly.AllowedDomains("www.bcv.org.ve"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	c.SetRequestTimeout(30 * time.Second)

	var tempUSD, tempEUR, tempValidity string

	c.OnHTML("#dolar strong", func(e *colly.HTMLElement) {
		tempUSD = strings.TrimSpace(e.Text)
	})

	c.OnHTML("#euro strong", func(e *colly.HTMLElement) {
		tempEUR = strings.TrimSpace(e.Text)
	})

	c.OnHTML("div.pull-right.dinpro.center span.date-display-single", func(e *colly.HTMLElement) {
		rawDate := strings.TrimSpace(e.Text)
		formattedDate, err := parseBCVDate(rawDate)
		if err != nil {
			log.Printf("No se pudo parsear la fecha [%s]: %v", rawDate, err)
			tempValidity = rawDate // backup if it fails
		} else {
			tempValidity = formattedDate // will be "2026-02-05"
		}
	})

	err := c.Visit("https://www.bcv.org.ve/")
	if err != nil {
		log.Printf("Error scraping: %v", err)
		return
	}

	// update the Store in a safe way
	store.Lock()
	store.USD = tempUSD
	store.EUR = tempEUR
	store.Validity = tempValidity
	store.UpdatedAt = time.Now()
	store.Unlock()

	log.Println("Data updated successfully")
}

func getExchangeHandler(w http.ResponseWriter, r *http.Request) {
	store.RLock()
	// copy the data to a struct to avoid data corruption
	data := struct {
		USD       string    `json:"usd"`
		EUR       string    `json:"eur"`
		Validity  string    `json:"validity"`
		UpdatedAt time.Time `json:"last_update"`
	}{
		USD:       store.USD,
		EUR:       store.EUR,
		Validity:  store.Validity,
		UpdatedAt: store.UpdatedAt,
	}
	store.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func parseBCVDate(rawDate string) (string, error) {
	// remove the day name and the comma
	parts := strings.Split(rawDate, ",")
	datePart := rawDate
	if len(parts) > 1 {
		datePart = strings.TrimSpace(parts[1])
	}

	// map of month names in Spanish
	meses := map[string]string{
		"Enero":      "01",
		"Febrero":    "02",
		"Marzo":      "03",
		"Abril":      "04",
		"Mayo":       "05",
		"Junio":      "06",
		"Julio":      "07",
		"Agosto":     "08",
		"Septiembre": "09",
		"Octubre":    "10",
		"Noviembre":  "11",
		"Diciembre":  "12",
	}

	// replace the month by its numeric value
	// "05 Febrero 2026" -> "05 02 2026"
	for mesNom, mesNum := range meses {
		if strings.Contains(datePart, mesNom) {
			datePart = strings.Replace(datePart, mesNom, mesNum, 1)
			break
		}
	}

	// parse using the Go layout
	// the layout should match "05 02 2026" -> "02 01 2006"
	t, err := time.Parse("02 01 2006", datePart)
	if err != nil {
		return "", fmt.Errorf("Error parsing date: %v", err)
	}

	// return in YYYY-MM-DD format
	return t.Format("2006-01-02"), nil
}
