# BCV Exchange Rate Scraper & API

A service written in **Go** that scrapes official exchange rates from the Central Bank of Venezuela (BCV). It provides a clean RESTful API for downstream consumption, featuring automated scheduling and data normalization.

## Overview

The BCV website updates its exchange rates twice daily. This service automates the retrieval of these rates (USD/EUR) and the "Value Date," transforming unstructured HTML into a standardized JSON format.

### Key Features
* **Modular Architecture**: Follows the Standard Go Project Layout for high maintainability.
* **Concurrent Safety**: Implements `sync.RWMutex` with "Atomic Snapshots" to prevent data races during API reads.
* **Automated Scheduling**: Internal cron jobs execute updates twice daily (9:00 AM & 4:00 PM).
* **Data Normalization**: 
    * Converts Spanish dates (e.g., "Jueves, 05 Febrero 2026") to ISO 8601 (`YYYY-MM-DD`).
    * Parses string values ("36,52") into precise `float64` types for financial calculations.
* **Resilient Networking**: Custom TLS transport to handle the BCV's non-standard certificate authority.

---

## Project Structure

The project is organized into decoupled packages to separate concerns:

```text
/bcv-scraper
├── cmd/
│   └── server/          # Application entry point (main.go)
├── internal/
│   ├── api/             # HTTP Handlers and JSON Response models
│   ├── scraper/         # Scraping logic and string-to-date parsing
│   └── store/           # Thread-safe global memory state
├── go.mod               # Dependency management
└── Dockerfile           # Multi-stage optimized build (pending)
```

---

## Tech Stack

* **Language:** Go 1.21+
* **Scraping:** [Colly v2](http://go-colly.org/) (Event-driven framework)
* **Scheduling:** [robfig/cron/v3](https://github.com/robfig/cron)
* **Deployment:** Docker (Alpine-based for minimal footprint)

---

## 📋 API Specification

### Get Latest Rates
Obtains a JSON object with the latest verified exchange rates.

* **Endpoint:** `GET /api/bcv/exchanges`
* **Response Example:**

```json
{
  "success": true,
  "data": {
    "bank": "Banco Central de Venezuela",
    "rates": {
      "usd": { "value": 36.52, "symbol": "Bs", "label": "Dólar" },
      "eur": { "value": 39.15, "symbol": "Bs", "label": "Euro" }
    },
    "date": "2026-02-05",
    "last_update": "2026-02-04T22:45:00Z"
  }
}
```

---

## ⚙️ Installation & Setup

### 1. Prerequisites
Before getting started, make sure you have the following installed:

* **Go 1.21+**
* **Docker** (optional)

### 2. Setup
Execute the following commands in your terminal to clone the repository and prepare the dependencies:

Clone the repository
```bash
git clone https://github.com/jssmntll/bcv-scraper.git
```

Enter the project directory
```bash
cd bcv-scraper
```

Dependencies (Go Modules)
```bash
go mod tidy
```

### 3. Running Locally
To run the application locally, execute the following command:

```bash
go run cmd/server/main.go
```

### 4. Docker Build
Pending ...

---
