package dal

import (
	"database/sql"
	_ "github.com/glebarez/go-sqlite"
)

type Db struct {
	Db      *sql.DB
	Worlds  *Worlds
	Backups *Backups
	Servers *Servers
	Crons   *Crons
}

func (d *Db) ConnectAndCreate(path string) error {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return err
	}
	d.Db = db
	d.Worlds = &Worlds{db: db}
	d.Backups = &Backups{db: db}
	d.Servers = &Servers{db: db}
	d.Crons = &Crons{db: db}
	if err := d.Worlds.Create(); err != nil {
		return err
	}
	if err := d.Backups.Create(); err != nil {
		return err
	}
	if err := d.Servers.Create(); err != nil {
		return err
	}
	if err := d.Crons.Create(); err != nil {
		return err
	}
	return nil
}
