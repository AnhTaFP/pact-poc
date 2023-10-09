package main

import (
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type discount struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Value       float64   `json:"value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

const discountsDB = "discounts.db"

func main() {
	db, err := initDb(discountsDB)
	if err != nil {
		log.Fatal("error initiating db:", err.Error())
	}

	address := "localhost:8080"
	startServer(address, db)
}
