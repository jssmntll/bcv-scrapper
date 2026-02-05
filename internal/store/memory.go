package store

import (
	"sync"
	"time"
)

type DataStore struct {
	sync.RWMutex
	USD       string    `json:"usd"`
	EUR       string    `json:"eur"`
	Date      string    `json:"date"`
	UpdatedAt time.Time `json:"last_update"`
}

var GlobalStore = &DataStore{}
