package main

import (
	"database/sql"
	"errors"
	"time"
)

func initDb(dbFile string) (*sql.DB, error) {
	const createDb = `
CREATE TABLE IF NOT EXISTS discounts (
  id INTEGER NOT NULL PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  type VARCHAR(25) NOT NULL,
  value DECIMAL(10,5) NOT NULL,
  created_at DATE NOT NULL,
  updated_at DATE DEFAULT NULL
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
	_, err := db.Exec("INSERT INTO discounts VALUES(NULL, ?,?,?,?,?, NULL)", d.Title, d.Description, d.Type, d.Value, time.Now())

	return err
}

func getOne(db *sql.DB, id int) (*discount, error) {
	row := db.QueryRow("SELECT * FROM discounts WHERE id = ?", id)

	var d discount
	var updated sql.NullTime
	if err := row.Scan(&d.ID, &d.Title, &d.Description, &d.Type, &d.Value, &d.CreatedAt, &updated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotFound
		}

		return nil, err
	}

	if updated.Valid {
		d.UpdatedAt = updated.Time
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
	r, err := db.Exec("UPDATE discounts SET title = ?, description = ?, type = ?, value = ?, updated_at = ? WHERE id = ?", title, description, discountType, value, time.Now(), id)
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
		var updated sql.NullTime
		if err := rows.Scan(&d.ID, &d.Title, &d.Description, &d.Type, &d.Value, &d.CreatedAt, &updated); err != nil {
			return nil, err
		}

		if updated.Valid {
			d.UpdatedAt = updated.Time
		}

		ds = append(ds, d)
	}

	return ds, nil
}

func countDiscounts(db *sql.DB) (int, error) {
	r, err := db.Query("SELECT COUNT(*) FROM discounts")
	if err != nil {
		return 0, nil
	}

	var count int
	for r.Next() {
		if err := r.Scan(&count); err != nil {
			return 0, err
		}
	}

	return count, nil
}

var errNotFound = errors.New("discount not found")
