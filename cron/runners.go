package cron

// cleanOldBackups deletes backups older than days
func cleanOldBackups(days int) {

}

// backup current world
func backup(name string) {

}

func restore(name string, ifBackup bool) {
	if ifBackup {
		backup("")
	}

}
