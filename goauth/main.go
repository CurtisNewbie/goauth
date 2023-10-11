package main

import (
	"fmt"
	"os"

	"github.com/curtisnewbie/goauth/domain"
	"github.com/curtisnewbie/goauth/web"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/miso/miso"
)

const (
	addPathEventBus     = "goauth.add-path"
	addResourceEventBus = "goauth.add-resource"
)

func main() {
	common.LoadBuiltinPropagationKeys()

	miso.PreServerBootstrap(func(rail miso.Rail) error {
		if err := scheduleTasks(rail); err != nil { // schedule cron jobs
			return err
		}
		web.RegisterWebEndpoints(rail) // register http server endpoints
		return SubEventBus(rail)       // subscribe to event bus
	})

	miso.BootstrapServer(os.Args) // bootstrap server
}

func SubEventBus(rail miso.Rail) error {

	// event bus to report path asynchronously
	miso.SubEventBus(addPathEventBus, 2, func(rail miso.Rail, req domain.CreatePathReq) error {
		rail.Debugf("receive %+v", req)
		return domain.CreatePathIfNotExist(rail, req, common.NilUser())
	})

	// event bus to report resource asynchronously
	miso.SubEventBus(addResourceEventBus, 2, func(rail miso.Rail, req domain.CreateResReq) error {
		rail.Debugf("receive %+v", req)
		return domain.CreateResourceIfNotExist(rail, req, common.NilUser())
	})

	return nil
}

func scheduleTasks(rail miso.Rail) error {
	// distributed tasks
	var err error = miso.ScheduleNamedDistributedTask("*/15 * * * *", false, "LoadRoleResCacheTask", func(ec miso.Rail) error {
		return domain.LoadRoleResCache(ec)
	})
	if err != nil {
		return err
	}
	err = miso.ScheduleNamedDistributedTask("*/15 * * * *", false, "LoadPathResCacheTask", func(ec miso.Rail) error {
		return domain.LoadPathResCache(ec)
	})
	if err != nil {
		return err
	}

	// for the first time
	miso.PostServerBootstrapped(func(c miso.Rail) error {
		ec := miso.EmptyRail()
		if e := domain.LoadRoleResCache(ec); e != nil {
			return fmt.Errorf("failed to load role resource, %v", e)
		}
		return nil
	})
	miso.PostServerBootstrapped(func(c miso.Rail) error {
		ec := miso.EmptyRail()
		if e := domain.LoadPathResCache(ec); e != nil {
			return fmt.Errorf("failed to load path resource, %v", e)
		}
		return nil
	})
	return nil
}
