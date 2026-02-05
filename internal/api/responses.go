package api

import "time"

type RateDetail struct {
	Value  float64 `json:"value"`
	Symbol string  `json:"symbol"`
	Label  string  `json:"label"`
}

type ExchangeResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Bank       string                `json:"bank"`
		Rates      map[string]RateDetail `json:"rates"`
		Date       string                `json:"date"`
		LastUpdate time.Time             `json:"last_update"`
	} `json:"data"`
}
