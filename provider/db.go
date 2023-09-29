package main

import (
	"database/sql"
	"errors"
)

func initDb(dbFile string) (*sql.DB, error) {
	const createDb = `
CREATE TABLE IF NOT EXISTS discounts (
  id INTEGER NOT NULL PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  type VARCHAR(25) NOT NULL,
  value DECIMAL(10,5) NOT NULL
);
`

	db, err := sql.Open("sqlite3", dbFile)
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
