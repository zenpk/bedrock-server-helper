package dal

import (
	"database/sql"
	"errors"
	"github.com/zenpk/bedrock-server-helper/util"
)

type Backups struct {
	db        *sql.DB
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Timestamp int64  `json:"timestamp"`
	WorldId   int64  `json:"worldId"`
	Deleted   bool
}

func (b Backups) Create() error {
	_, err := b.db.Exec(`
	CREATE TABLE IF NOT EXISTS backups (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    name TEXT NOT NULL,
		timestamp INTEGER NOT NULL,
		world_id INTEGER NOT NULL,
		deleted INTEGER NOT NULL DEFAULT 0
	);`)
	return err
}

func (b Backups) ListByWorldId(worldId int64) ([]Backups, error) {
	backups := make([]Backups, 0)
	rows, err := b.db.Query(`SELECT * FROM backups WHERE ( world_id = ? AND deleted = 0) ORDER BY id DESC`, worldId)
	if err != nil {
		return backups, err
	}
	defer rows.Close()
	for rows.Next() {
		var backup Backups
		err = rows.Scan(&backup.Id, &backup.Name, &backup.Timestamp, &backup.WorldId, &backup.Deleted)
		if err != nil {
			return backups, err
		}
		backups = append(backups, backup)
	}
	return backups, nil
}

// Insert default name YYYY-MM-DD
func (b Backups) Insert(name string, worldId int64) error {
	if name == "" {
		return errors.New("backup name mustn't be empty")
	}
	if worldId <= 0 {
		return errors.New("world_id must be bigger than 0")
	}
	_, err := b.db.Exec("INSERT INTO backups (name, timestamp, world_id) VALUES (?, ?, ?);", name, util.UnixSeconds(), worldId)
	return err
}

func (b Backups) DeleteById(id int64) error {
	_, err := b.db.Exec("UPDATE backups SET deleted = 1 WHERE id = ?;", id)
	return err
}

func (b Backups) SelectDaysBefore(worldId, days int64) ([]Backups, error) {
	backups := make([]Backups, 0)
	beforeTimestamp := util.UnixSeconds() - days*24*60*60
	rows, err := b.db.Query(`SELECT * FROM backups WHERE (world_id = ? AND timestamp < ? AND deleted = 0) ORDER BY id DESC;`,
		worldId, beforeTimestamp)
	if err != nil {
		return backups, err
	}
	defer rows.Close()
	for rows.Next() {
		var backup Backups
		err = rows.Scan(&backup.Id, &backup.Name, &backup.Timestamp, &backup.WorldId, &backup.Deleted)
		if err != nil {
			return backups, err
		}
		backups = append(backups, backup)
	}
	return backups, nil
}

func (b Backups) SelectById(id int64) (Backups, error) {
	rows, err := b.db.Query("SELECT * FROM backups WHERE (id = ? AND deleted = 0);", id)
	if err != nil {
		return Backups{}, err
	}
	defer rows.Close()
	var backup Backups
	if rows.Next() {
		err = rows.Scan(&backup.Id, &backup.Name, &backup.Timestamp, &backup.WorldId, &backup.Deleted)
		if err != nil {
			return backup, err
		}
	} else {
		return backup, errors.New("no backup found")
	}
	return backup, nil
}

// ResolveName ensures that the backup name is legal and unique
func (b Backups) ResolveName(name string) (string, error) {
	if name == "" {
		name = util.UnixYyyyMmDd()
	}
	// dealing with name collision
	for {
		rows, err := b.db.Query("SELECT * FROM backups WHERE (name = ? AND deleted = 0);", name)
		if err != nil {
			return "", err
		}
		if !rows.Next() {
			break
		}
		if err := rows.Close(); err != nil {
			return "", err
		}
		name += "1"
	}
	return name, nil
}
