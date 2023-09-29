package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type discount struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	Value       float64 `json:"value"`
}

func main() {
	db, err := initDb()
	if err != nil {
		log.Fatal("error initiating db:", err.Error())
	}

	r := mux.NewRouter()

	r.HandleFunc("/discounts", func(w http.ResponseWriter, r *http.Request) {
		var d discount
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := insert(db, d); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}).Methods("POST")

	r.HandleFunc("/discounts", func(w http.ResponseWriter, r *http.Request) {

	}).Methods("GET")

	r.HandleFunc("/discounts/{id}", func(w http.ResponseWriter, r *http.Request) {

	}).Methods("GET")

	r.HandleFunc("/discounts/{id}", func(w http.ResponseWriter, r *http.Request) {

	}).Methods("PUT")

	r.HandleFunc("/discounts/{id}", func(w http.ResponseWriter, r *http.Request) {

	}).Methods("DELETE")

	if err := http.ListenAndServe("localhost:8080", r); err != nil {
		log.Fatal("error occured:", err.Error())
	}
}

func initDb() (*sql.DB, error) {
	const file = "discounts.db"
	const createDb = `
CREATE TABLE IF NOT EXISTS discounts (
  id INTEGER NOT NULL PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  type VARCHAR(25) NOT NULL,
  value DECIMAL(10,5) NOT NULL
);
`

	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(createDb); err != nil {
		return nil, err
	}

	return db, nil
}

func insert(db *sql.DB, d discount) error {
	_, err := db.Exec("INSERT INTO discounts VALUES(NULL, ?,?,?,?)", d.Title, d.Description, d.Type, d.Value)

	return err
}
