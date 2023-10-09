package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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

func startServer(address string, db *sql.DB) {
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
		typeParam := r.URL.Query().Get("type")

		discounts, err := queryDiscounts(db, typeParam)
		if err != nil {
			if errors.Is(err, errNotFound) {
				http.Error(w, "", http.StatusNotFound)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var resp struct {
			Discounts []discount `json:"discounts"`
		}

		resp.Discounts = discounts
		b, _ := json.Marshal(resp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}).Methods("GET")

	r.HandleFunc("/discounts/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		discountID, _ := strconv.Atoi(id)
		d, err := getOne(db, discountID)
		if err != nil {
			if errors.Is(err, errNotFound) {
				http.Error(w, "", http.StatusNotFound)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b, _ := json.Marshal(d)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}).Methods("GET")

	r.HandleFunc("/discounts/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		discountID, _ := strconv.Atoi(id)

		var body struct {
			Title       string  `json:"title"`
			Description string  `json:"description"`
			Type        string  `json:"type"`
			Value       float64 `json:"value"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := update(db, discountID, body.Title, body.Description, body.Type, body.Value); err != nil {
			if errors.Is(err, errNotFound) {
				http.Error(w, "", http.StatusNotFound)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}).Methods("PUT")

	r.HandleFunc("/discounts/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		discountID, _ := strconv.Atoi(id)

		if err := deleteOne(db, discountID); err != nil {
			if errors.Is(err, errNotFound) {
				http.Error(w, "", http.StatusNotFound)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}).Methods("DELETE")

	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatal("error occurred:", err.Error())
	}
}
