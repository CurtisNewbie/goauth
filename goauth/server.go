package goauth

import (
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/miso/miso"
)

func BootstrapServer(args []string) {
	common.LoadBuiltinPropagationKeys()
	miso.PreServerBootstrap(ScheduleTasks)
	miso.PreServerBootstrap(SubEventBus)
	miso.PreServerBootstrap(RegisterWebEndpoints)
	miso.BootstrapServer(args)
}
