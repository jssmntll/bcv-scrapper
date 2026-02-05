package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jssmntll/bcv-scrapper/internal/store"
)

func GetExchangeHandler(w http.ResponseWriter, r *http.Request) {
	// copy the data to a struct to avoid data corruption
	store.GlobalStore.RLock()
	rawUSD := store.GlobalStore.USD
	rawEUR := store.GlobalStore.EUR
	validityDate := store.GlobalStore.Date
	lastUpdate := store.GlobalStore.UpdatedAt
	store.GlobalStore.RUnlock()

	parse := func(s string) float64 {
		if s == "" {
			return 0
		}
		v, _ := strconv.ParseFloat(strings.Replace(s, ",", ".", 1), 64)
		return v
	}

	response := ExchangeResponse{
		Success: true,
	}
	response.Data.Bank = "Banco Central de Venezuela"
	response.Data.Date = validityDate
	response.Data.LastUpdate = lastUpdate
	response.Data.Rates = map[string]RateDetail{
		"usd": {Value: parse(rawUSD), Symbol: "$", Label: "Dólar"},
		"eur": {Value: parse(rawEUR), Symbol: "€", Label: "Euro"},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
