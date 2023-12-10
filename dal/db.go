package dal

import (
	"database/sql"
	_ "github.com/glebarez/go-sqlite"
)

type Db struct {
	Db *sql.DB
}

func (d *Db) Connect(path string) error {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return err
	}
	d.Db = db
	return nil
}
