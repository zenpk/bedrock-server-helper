package dal

import (
	"database/sql"
	"errors"
)

type Worlds struct {
	db          *sql.DB
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Properties  string `json:"properties"`
	AllowList   string `json:"allowList"`
	HasSaveData bool
	UsingServer int64 `json:"usingServer"`
	Deleted     bool
}

func (w Worlds) Create() error {
	_, err := w.db.Exec(`
	CREATE TABLE IF NOT EXISTS worlds (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    name TEXT NOT NULL UNIQUE,
	    properties TEXT NOT NULL,
	    allow_list TEXT,
	    has_save_data INTEGER NOT NULL DEFAULT 0,
	    using_server INTEGER NOT NULL DEFAULT 0,
		deleted INTEGER NOT NULL DEFAULT 0
	);`)
	return err
}

func (w Worlds) List() ([]Worlds, error) {
	worlds := make([]Worlds, 0)
	rows, err := w.db.Query(`SELECT * FROM worlds WHERE deleted = 0 ORDER BY id DESC;`)
	if err != nil {
		return worlds, err
	}
	defer rows.Close()
	for rows.Next() {
		var world Worlds
		err = rows.Scan(&world.Id, &world.Name, &world.Properties, &world.AllowList, &world.HasSaveData, &world.UsingServer, &world.Deleted)
		if err != nil {
			return worlds, err
		}
		worlds = append(worlds, world)
	}
	return worlds, nil
}

func (w Worlds) Insert(name, properties, allowList string, serverId int64) error {
	_, err := w.db.Exec("INSERT INTO worlds (name, properties, allow_list, using_server) VALUES (?, ?, ?, ?);", name, properties, allowList, serverId)
	return err
}

func (w Worlds) DeleteById(id int64) error {
	_, err := w.db.Exec("UPDATE worlds SET deleted = 1 WHERE id = ?;", id)
	return err
}

func (w Worlds) SelectById(id int64) (Worlds, error) {
	rows, err := w.db.Query("SELECT * FROM worlds WHERE (id = ? AND deleted = 0);", id)
	if err != nil {
		return Worlds{}, err
	}
	defer rows.Close()
	var world Worlds
	if rows.Next() {
		if err := rows.Scan(&world.Id, &world.Name, &world.Properties, &world.AllowList, &world.HasSaveData, &world.UsingServer, &world.Deleted); err != nil {
			return Worlds{}, err
		}
	} else {
		return Worlds{}, errors.New("world not found")
	}
	return world, err
}

func (w Worlds) SetHasSaveData(id int64, hasSaveData bool) error {
	_, err := w.db.Exec("UPDATE worlds SET has_save_data = ? WHERE (id = ? AND deleted = 0);", hasSaveData, id)
	return err
}

func (w Worlds) SetUsingServer(id, serverId int64) error {
	_, err := w.db.Exec("UPDATE worlds SET using_server = ? WHERE (id = ? AND deleted = 0);", serverId, id)
	return err
}
