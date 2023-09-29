package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

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

func getOne(db *sql.DB, id int) (*discount, error) {
	row := db.QueryRow("SELECT * FROM discounts WHERE id = ?", id)

	var d discount
	if err := row.Scan(&d.ID, &d.Title, &d.Description, &d.Type, &d.Value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotFound
		}

		return nil, err
	}

	return &d, nil
}

func deleteOne(db *sql.DB, id int) error {
	r, err := db.Exec("DELETE FROM discounts WHERE id = ?", id)
	if err != nil {
		return err
	}

	affected, _ := r.RowsAffected()

	if affected == 0 {
		return errNotFound
	}

	return nil
}

func update(db *sql.DB, id int, title string, description string, discountType string, value float64) error {
	r, err := db.Exec("UPDATE discounts SET title = ?, description = ?, type = ?, value = ? WHERE id = ?", title, description, discountType, value, id)
	if err != nil {
		return err
	}

	affected, _ := r.RowsAffected()
	if affected == 0 {
		return errNotFound
	}

	return nil
}

func queryDiscounts(db *sql.DB, typeParam string) ([]discount, error) {
	rows, err := db.Query("SELECT * FROM discounts WHERE type = ?", typeParam)
	if err != nil {
		return nil, err
	}

	ds := make([]discount, 0)
	for rows.Next() {
		var d discount
		if err := rows.Scan(&d.ID, &d.Title, &d.Description, &d.Type, &d.Value); err != nil {
			return nil, err
		}

		ds = append(ds, d)
	}

	return ds, nil
}

var errNotFound = errors.New("discount not found")
