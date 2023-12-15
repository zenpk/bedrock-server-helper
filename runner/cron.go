package runner

import (
	"github.com/go-co-op/gocron/v2"
)

func (r Runner) StartCron() (gocron.Scheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}
	//backupJob, err := scheduler.NewJob(
	//	gocron.CronJob("55 06 * * *", false),
	//	gocron.NewTask(
	//		r.backup,
	//		"",
	//	),
	//)
	//if err != nil {
	//	return nil, err
	//}
	//log.Printf("backup job %v started", backupJob.ID())
	//cleanJob, err := scheduler.NewJob(
	//	gocron.CronJob("55 07 * * *", false),
	//	gocron.NewTask(
	//		r.cleanOldBackups,
	//		7,
	//	),
	//)
	//if err != nil {
	//	return nil, err
	//}
	//log.Printf("clean job %v started", cleanJob.ID())
	scheduler.Start()
	return scheduler, nil
}
