# BCV Exchange Rate Scraper & API

A service written in **Go** that scrapes official exchange rates from the Central Bank of Venezuela (BCV). It provides a clean RESTful API for downstream consumption, featuring automated scheduling and data normalization.

## 🚀 Overview

The BCV website updates its exchange rates twice daily. This service automates the retrieval of these rates (USD/EUR) and the "Value Date," transforming unstructured HTML into a standardized JSON format.
