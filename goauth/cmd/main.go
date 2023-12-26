package main

import (
	"os"

	"github.com/curtisnewbie/goauth"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/miso/miso"
)

func main() {
	common.LoadBuiltinPropagationKeys()
	miso.PreServerBootstrap(goauth.ScheduleTasks)
	miso.PreServerBootstrap(goauth.SubEventBus)
	miso.PreServerBootstrap(goauth.RegisterWebEndpoints)
	miso.BootstrapServer(os.Args)
}
