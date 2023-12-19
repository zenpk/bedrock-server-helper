package dal

import (
	"database/sql"
)

type Crons struct {
	db         *sql.DB
	Id         int64
	Name       string // better to use id, but for simplicity
	Cron       string
	Parameters string // could be any format, for example, for the clean job, it should be a number
	WorldId    int64
	Deleted    bool
}

func (c Crons) Create() error {
	_, err := c.db.Exec(`
	CREATE TABLE IF NOT EXISTS crons (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    name TEXT NOT NULL,
		cron TEXT NOT NULL,
		parameters TEXT,
		world_id INTEGER NOT NULL,
		deleted INTEGER NOT NULL DEFAULT 0
	);`)
	return err
}

func (c Crons) List() ([]Crons, error) {
	crons := make([]Crons, 0)
	rows, err := c.db.Query(`SELECT * FROM crons WHERE deleted = 0;`)
	if err != nil {
		return crons, err
	}
	defer rows.Close()
	for rows.Next() {
		var cron Crons
		err = rows.Scan(&cron.Id, &cron.Name, &cron.Cron, &cron.Parameters, &cron.WorldId, &cron.Deleted)
		if err != nil {
			return crons, err
		}
		crons = append(crons, cron)
	}
	return crons, nil
}
