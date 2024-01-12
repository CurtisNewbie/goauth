package goauth

import (
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/miso/miso"
)

const (
	addPathEventBus       = "event.bus.goauth.add-path"
	addResourceEventBus   = "event.bus.goauth.add-resource"
	addPathEventBusV2     = "goauth.add-path"
	addResourceEventBusV2 = "goauth.add-resource"
)

func SubEventBus(rail miso.Rail) error {
	// event bus to report path asynchronously
	miso.SubEventBus(addPathEventBus, 2, ListenAddPathEvent)
	miso.SubEventBus(addPathEventBusV2, 2, ListenAddPathEvent)

	// event bus to report resource asynchronously
	miso.SubEventBus(addResourceEventBus, 2, ListenAddResourceEvent)
	miso.SubEventBus(addResourceEventBusV2, 2, ListenAddResourceEvent)

	return nil
}

func ListenAddResourceEvent(rail miso.Rail, req CreateResReq) error {
	rail.Debugf("receive %+v", req)
	return CreateResourceIfNotExist(rail, req, common.NilUser())
}

func ListenAddPathEvent(rail miso.Rail, req CreatePathReq) error {
	rail.Debugf("receive %+v", req)
	return CreatePathIfNotExist(rail, req, common.NilUser())
}
