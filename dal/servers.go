package dal

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

type Servers struct {
	db            *sql.DB
	Id            int64
	Version       string
	VersionNumber int64 // for sorting, calculated from Version. E.g. 1.20.50.01 = 1*1000000 + 20*10000 + 50*100 + 1 = 12050001
	Deleted       bool
}

func (s Servers) Create() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS servers (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    version TEXT NOT NULL,
	    version_number INTEGER NOT NULL DEFAULT 0,
		deleted INTEGER NOT NULL DEFAULT 0
	);`)
	return err
}

func (s Servers) List() ([]Servers, error) {
	servers := make([]Servers, 0)
	rows, err := s.db.Query(`SELECT * FROM servers WHERE deleted = 0 ORDER BY version_number DESC;`)
	if err != nil {
		return servers, err
	}
	defer rows.Close()
	for rows.Next() {
		var server Servers
		err = rows.Scan(&server.Id, &server.Version, &server.VersionNumber, &server.Deleted)
		if err != nil {
			return servers, err
		}
		servers = append(servers, server)
	}
	return servers, nil
}

func (s Servers) Insert(version string) error {
	if version == "" {
		return errors.New("version mustn't be empty")
	}
	// make sure there is no server with the same version
	rows, err := s.db.Query("SELECT * FROM servers WHERE (version = ? AND deleted = 0);", version)
	if err != nil {
		return err
	}
	if rows.Next() {
		return errors.New("the world already has this version of server")
	}
	if err := rows.Close(); err != nil {
		return err
	}
	// calculate the version number
	versionNumberStrs := strings.Split(version, ".")
	versionNumbers := make([]int64, 4)
	for i, versionNumberStr := range versionNumberStrs {
		versionNumbers[i], err = strconv.ParseInt(versionNumberStr, 10, 64)
		if err != nil {
			return err
		}
	}
	versionNumber := versionNumbers[0]*1000000 + versionNumbers[1]*10000 + versionNumbers[2]*100 + versionNumbers[3]
	_, err = s.db.Exec("INSERT INTO servers (version, version_number) VALUES (?, ?);", version, versionNumber)
	return err
}

func (s Servers) SelectById(id int64) (Servers, error) {
	rows, err := s.db.Query("SELECT * FROM servers WHERE (id = ? AND deleted = 0);", id)
	if err != nil {
		return Servers{}, err
	}
	defer rows.Close()
	var server Servers
	for rows.Next() {
		err = rows.Scan(&server.Id, &server.Version, &server.VersionNumber, &server.Deleted)
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
	defer rows.Close()
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
	_, err = s.db.Exec("UPDATE servers SET deleted = 1 WHERE id = ?;", id)
	return err
}

func (s Servers) SelectLatest() (Servers, error) {
	rows, err := s.db.Query("SELECT * FROM servers WHERE deleted = 0 ORDER BY version_number DESC LIMIT 1;")
	if err != nil {
		return Servers{}, err
	}
	defer rows.Close()
	var server Servers
	for rows.Next() {
		err = rows.Scan(&server.Id, &server.Version, &server.VersionNumber, &server.Deleted)
		if err != nil {
			return server, err
		}
	}
	return server, nil
}
