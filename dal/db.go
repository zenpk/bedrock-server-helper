package dal

import (
	"database/sql"
	_ "github.com/glebarez/go-sqlite"
)

type Db struct {
	Db       *sql.DB
	Backups  *Backups
	Versions *Versions
}

func (d *Db) Connect(path string) error {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return err
	}
	d.Db = db
	d.Backups = &Backups{db: db}
	d.Versions = &Versions{db: db}
	if err := d.Backups.Create(); err != nil {
		return err
	}
	if err := d.Versions.Create(); err != nil {
		return err
	}
	return nil
}
