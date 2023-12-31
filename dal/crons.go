package dal

import (
	"database/sql"
	"errors"
)

type Crons struct {
	db         *sql.DB
	Id         int64
	JobName    string // better to use id, but for simplicity
	Cron       string
	Parameters string // could be any format, for example, for the clean job, it should be a number
	WorldId    int64
	Deleted    bool
}

func (c Crons) Create() error {
	_, err := c.db.Exec(`
	CREATE TABLE IF NOT EXISTS crons (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    job_name TEXT NOT NULL,
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
		err = rows.Scan(&cron.Id, &cron.JobName, &cron.Cron, &cron.Parameters, &cron.WorldId, &cron.Deleted)
		if err != nil {
			return crons, err
		}
		crons = append(crons, cron)
	}
	return crons, nil
}

func (c Crons) Insert(jobName string, worldId int64, parameters string, cron string) error {
	if jobName == "" {
		return errors.New("job name mustn't be empty")
	}
	if worldId <= 0 {
		return errors.New("world id must be bigger than 0")
	}
	if cron == "" {
		return errors.New("cron mustn't be empty")
	}
	_, err := c.db.Exec("INSERT INTO crons (job_name, cron, parameters, world_id) VALUES (?, ?, ?, ?);", jobName, cron, parameters, worldId)
	return err
}

func (c Crons) SelectByWorldId(worldId int64) ([]Crons, error) {
	crons := make([]Crons, 0)
	rows, err := c.db.Query(`SELECT * FROM crons WHERE ( world_id = ? AND deleted = 0) ORDER BY id DESC;`, worldId)
	if err != nil {
		return crons, err
	}
	defer rows.Close()
	for rows.Next() {
		var cron Crons
		err = rows.Scan(&cron.Id, &cron.JobName, &cron.Cron, &cron.Parameters, &cron.WorldId, &cron.Deleted)
		if err != nil {
			return crons, err
		}
		crons = append(crons, cron)
	}
	return crons, nil
}

func (c Crons) DeleteById(id int64) error {
	_, err := c.db.Exec("UPDATE crons SET deleted = 1 WHERE id = ?;", id)
	return err
}
