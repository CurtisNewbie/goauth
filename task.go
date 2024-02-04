package goauth

import (
	"github.com/curtisnewbie/miso/miso"
)

func ScheduleTasks(rail miso.Rail) error {
	// distributed tasks
	var err error = miso.ScheduleDistributedTask(miso.Job{
		Cron:                   "*/15 * * * *",
		CronWithSeconds:        false,
		Name:                   "LoadRoleResCacheTask",
		TriggeredOnBoostrapped: true,
		Run: func(ec miso.Rail) error {
			return LoadRoleResCache(ec)
		}})
	if err != nil {
		return err
	}
	err = miso.ScheduleDistributedTask(miso.Job{
		Cron:                   "*/15 * * * *",
		CronWithSeconds:        false,
		Name:                   "LoadPathResCacheTask",
		TriggeredOnBoostrapped: true,
		Run: func(ec miso.Rail) error {
			return LoadPathResCache(ec)
		}})
	if err != nil {
		return err
	}
	err = miso.ScheduleDistributedTask(miso.Job{
		Cron:                   "*/15 * * * *",
		CronWithSeconds:        false,
		Name:                   "LoadResCodeCacheTask",
		TriggeredOnBoostrapped: true,
		Run: func(ec miso.Rail) error {
			return LoadResCodeCache(ec)
		}})
	if err != nil {
		return err
	}
	return nil
}
