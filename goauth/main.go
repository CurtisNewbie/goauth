package main

import (
	"fmt"
	"os"

	"github.com/curtisnewbie/goauth/domain"
	"github.com/curtisnewbie/goauth/web"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/miso/bus"
	"github.com/curtisnewbie/miso/core"
	"github.com/curtisnewbie/miso/server"
	"github.com/curtisnewbie/miso/task"
)

const (
	addPathEventBus     = "goauth.add-path"
	addResourceEventBus = "goauth.add-resource"
)

func main() {
	server.PreServerBootstrap(func(rail core.Rail) error {
		if err := scheduleTasks(rail); err != nil { // schedule cron jobs
			return err
		}
		web.RegisterWebEndpoints(rail) // register http server endpoints
		return SubscribeEventBus(rail) // subscribe to event bus
	})

	server.BootstrapServer(os.Args) // bootstrap server
}

func SubscribeEventBus(rail core.Rail) error {

	// event bus to report path asynchronously
	bus.SubscribeEventBus(addPathEventBus, 2, func(rail core.Rail, req domain.CreatePathReq) error {
		rail.Debugf("receive %+v", req)
		return domain.CreatePathIfNotExist(rail, req, common.NilUser())
	})

	// event bus to report resource asynchronously
	bus.SubscribeEventBus(addResourceEventBus, 2, func(rail core.Rail, req domain.CreateResReq) error {
		rail.Debugf("receive %+v", req)
		return domain.CreateResourceIfNotExist(rail, req, common.NilUser())
	})

	return nil
}

func scheduleTasks(rail core.Rail) error {
	// distributed tasks
	var err error = task.ScheduleNamedDistributedTask("*/15 * * * *", false, "LoadRoleResCacheTask", func(ec core.Rail) error {
		return domain.LoadRoleResCache(ec)
	})
	if err != nil {
		return err
	}
	err = task.ScheduleNamedDistributedTask("*/15 * * * *", false, "LoadPathResCacheTask", func(ec core.Rail) error {
		return domain.LoadPathResCache(ec)
	})
	if err != nil {
		return err
	}

	// for the first time
	server.PostServerBootstrapped(func(c core.Rail) error {
		ec := core.EmptyRail()
		if e := domain.LoadRoleResCache(ec); e != nil {
			return fmt.Errorf("failed to load role resource, %v", e)
		}
		return nil
	})
	server.PostServerBootstrapped(func(c core.Rail) error {
		ec := core.EmptyRail()
		if e := domain.LoadPathResCache(ec); e != nil {
			return fmt.Errorf("failed to load path resource, %v", e)
		}
		return nil
	})
	return nil
}
