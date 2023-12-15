package dal

import (
	"database/sql"
	"errors"
)

type Servers struct {
	db      *sql.DB
	Id      int64
	Version string
	WorldId int64
	Deleted bool
}

func (s Servers) Create() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS servers (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    version TEXT NOT NULL,
	    world_id INTEGER NOT NULL,
		deleted INTEGER NOT NULL DEFAULT 0
	);`)
	return err
}

func (s Servers) ListByWorldId(worldId int64) ([]Servers, error) {
	servers := make([]Servers, 0)
	rows, err := s.db.Query(`SELECT * FROM servers WHERE (deleted = 0 AND world_id = ?) ORDER BY id DESC;`, worldId)
	defer rows.Close()
	for rows.Next() {
		var server Servers
		err = rows.Scan(&server.Id, &server.Version, &server.WorldId)
		if err != nil {
			return servers, err
		}
		servers = append(servers, server)
	}
	return servers, nil
}

func (s Servers) Insert(version string, worldId int64) error {
	if version == "" {
		return errors.New("version mustn't be empty")
	}
	if worldId <= 0 {
		return errors.New("world_id must be bigger than 0")
	}
	_, err := s.db.Exec("INSERT INTO versions (version, world_id) VALUES (?, ?);", version, worldId)
	return err
}

func (s Servers) SelectById(id int64) (Servers, error) {
	rows, err := s.db.Query("SELECT * FROM servers WHERE id = ?;", id)
	var server Servers
	for rows.Next() {
		err = rows.Scan(&server.Id, &server.Version, &server.WorldId)
		if err != nil {
			return server, err
		}
	}
	return server, nil
}

func (s Servers) IsInUse(id int64) (bool, error) {
	rows, err := s.db.Query("SELECT * FROM worlds WHERE (using_server = ? AND deleted = 0);", id)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func (s Servers) DeleteById(id int64) error {
	// check if in use
	inUse, err := s.IsInUse(id)
	if err != nil {
		return err
	}
	if inUse {
		return errors.New("server is in use")
	}
	_, err = s.db.Exec("UPDATE versions SET deleted = 1 WHERE id = ?;", id)
	return err
}
