package dal

import (
	"database/sql"
)

type Worlds struct {
	db *sql.DB
	Id int64
}

func (w *Worlds) Create() error {
	_, err := w.db.Exec(`
	CREATE TABLE IF NOT EXISTS worlds (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    name TEXT NOT NULL
	);`)
	return err
}

func (w *Worlds) List() error {
	rows, err := w.db.Exec(`
	SELECT * FROM worlds
	`)
	rows.
	return err
}
