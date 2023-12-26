package goauth

import (
	"fmt"

	"github.com/curtisnewbie/miso/miso"
)

func ScheduleTasks(rail miso.Rail) error {
	// distributed tasks
	var err error = miso.ScheduleDistributedTask(miso.Job{
		Cron:            "*/15 * * * *",
		CronWithSeconds: false,
		Name:            "LoadRoleResCacheTask",
		Run: func(ec miso.Rail) error {
			return LoadRoleResCache(ec)
		}})
	if err != nil {
		return err
	}
	err = miso.ScheduleDistributedTask(miso.Job{
		Cron:            "*/15 * * * *",
		CronWithSeconds: false,
		Name:            "LoadPathResCacheTask",
		Run: func(ec miso.Rail) error {
			return LoadPathResCache(ec)
		}})
	if err != nil {
		return err
	}

	// for the first time
	miso.PostServerBootstrapped(func(c miso.Rail) error {
		ec := miso.EmptyRail()
		if e := LoadRoleResCache(ec); e != nil {
			return fmt.Errorf("failed to load role resource, %v", e)
		}
		return nil
	})
	miso.PostServerBootstrapped(func(c miso.Rail) error {
		ec := miso.EmptyRail()
		if e := LoadPathResCache(ec); e != nil {
			return fmt.Errorf("failed to load path resource, %v", e)
		}
		return nil
	})
	return nil
}
