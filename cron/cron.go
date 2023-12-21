package cron

import (
	"errors"
	"github.com/go-co-op/gocron/v2"
	"github.com/zenpk/bedrock-server-helper/dal"
	"github.com/zenpk/bedrock-server-helper/runner"
	"log"
	"strconv"
)

type Cron struct {
	Db        *dal.Db
	Runner    *runner.Runner
	scheduler gocron.Scheduler
}

func (c *Cron) RefreshCron() error {
	if c.scheduler != nil {
		if err := c.scheduler.Shutdown(); err != nil {
			return err
		}
	}
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return err
	}
	crons, err := c.Db.Crons.List()
	if err != nil {
		return err
	}
	for _, cron := range crons {
		var job gocron.Job
		if cron.Name == JobBackup {
			job, err = scheduler.NewJob(gocron.CronJob(cron.Cron, false),
				gocron.NewTask(
					c.Runner.Backup,
					"",
					cron.WorldId,
					nil,
				),
			)
			if err != nil {
				return err
			}
		}
		if cron.Name == JobClean {
			days, err := strconv.ParseInt(cron.Parameters, 10, 64)
			if err != nil {
				return err
			}
			if days <= 0 {
				return errors.New("cannot clean future backups")
			}
			job, err = scheduler.NewJob(gocron.CronJob(cron.Cron, false),
				gocron.NewTask(
					c.Runner.CleanOldBackups,
					cron.WorldId,
					days,
				),
			)
			if err != nil {
				return nil
			}
		}
		log.Printf("job: %v, world id: %v, cron name: %v, started\n", cron.Name, cron.WorldId, job.Name())
	}
	scheduler.Start()
	c.scheduler = scheduler
	return nil
}