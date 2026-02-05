package scrapper

import (
	"fmt"
	"strings"
	"time"
)

func ParseBCVDate(rawDate string) (string, error) {
	// remove the day name and the comma
	parts := strings.Split(rawDate, ",")
	datePart := rawDate
	if len(parts) > 1 {
		datePart = strings.TrimSpace(parts[1])
	}

	// map of month names in Spanish
	months := map[string]string{
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
	for monthName, monthNumber := range months {
		if strings.Contains(datePart, monthName) {
			datePart = strings.Replace(datePart, monthName, monthNumber, 1)
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
