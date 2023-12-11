package dal

import (
	"database/sql"
	"github.com/zenpk/bedrock-server-helper/util"
)

type Backups struct {
	db        *sql.DB
	Id        int64
	Name      string
	Timestamp int64
	Deleted   bool
}

func (b Backups) Create() error {
	_, err := b.db.Exec(`
	CREATE TABLE IF NOT EXISTS backups (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    name TEXT NOT NULL,
		timestamp INTEGER NOT NULL,
		deleted INTEGER NOT NULL DEFAULT 0
	);`)
	return err
}

func (b Backups) List() ([]Backups, error) {
	backups := make([]Backups, 0)
	rows, err := b.db.Query(`SELECT * FROM backups WHERE deleted = 0 ORDER BY id DESC`)
	defer rows.Close()
	for rows.Next() {
		var backup Backups
		err = rows.Scan(&backup.Id, &backup.Name, &backup.Timestamp)
		if err != nil {
			return backups, err
		}
		backups = append(backups, backup)
	}
	return backups, nil
}

// Insert default name YYYY-MM-DD
func (b Backups) Insert(name string) error {
	if name == "" {
		name = util.UnixYyyyMmDd()
	}
	// dealing with name collision
	for {
		rows, err := b.db.Query("SELECT * FROM backups WHERE (name = ? AND deleted != 0);", name)
		if err != nil {
			return err
		}
		if !rows.Next() {
			break
		}
		name += "1"
	}
	_, err := b.db.Exec("INSERT INTO backups (name, timestamp) VALUES (?, ?);", name, util.UnixSeconds())
	return err
}

func (b Backups) DeleteByName(name string) error {
	_, err := b.db.Exec("UPDATE backups SET deleted = 1 WHERE name = ?;", name)
	return err
}

func (b Backups) SelectDaysBefore(days int64) ([]Backups, error) {
	backups := make([]Backups, 0)
	beforeTimestamp := util.UnixSeconds() - days*24*60*60
	rows, err := b.db.Query(`SELECT * FROM backups WHERE (deleted = 0 AND timestamp < ?) ORDER BY id DESC;`,
		beforeTimestamp)
	defer rows.Close()
	for rows.Next() {
		var backup Backups
		err = rows.Scan(&backup.Id, &backup.Name, &backup.Timestamp)
		if err != nil {
			return backups, err
		}
		backups = append(backups, backup)
	}
	return backups, nil
}
