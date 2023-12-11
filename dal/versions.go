package dal

import (
	"database/sql"
	"github.com/zenpk/bedrock-server-helper/util"
)

type Versions struct {
	db      *sql.DB
	Id      int64
	Name    string
	Deleted bool
}

func (v Versions) Create() error {
	_, err := v.db.Exec(`
	CREATE TABLE IF NOT EXISTS versions (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    name TEXT NOT NULL,
		deleted INTEGER NOT NULL DEFAULT 0
	);`)
	return err
}

func (v Versions) List() ([]Versions, error) {
	versions := make([]Versions, 0)
	rows, err := v.db.Query(`SELECT * FROM versions WHERE deleted = 0 ORDER BY id DESC;`)
	defer rows.Close()
	for rows.Next() {
		var version Versions
		err = rows.Scan(&version.Id, &version.Name)
		if err != nil {
			return versions, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}

func (v Versions) Insert(name string) error {
	if name == "" {
		name = util.UnixYyyyMmDd()
	}
	_, err := v.db.Exec("INSERT INTO versions (name) VALUES (?);", name)
	return err
}

func (v Versions) DeleteByName(name string) error {
	_, err := v.db.Exec("UPDATE versions SET deleted = 1 WHERE name = ?;", name)
	return err
}
