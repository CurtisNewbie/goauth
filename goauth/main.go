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
	"github.com/gin-gonic/gin"
)

const (
	codeMngResources    = "manage-resources"
	nameMngReesources   = "Manage Resources Access"
	addPathEventBus     = "goauth.add-path"
	addResourceEventBus = "goauth.add-resource"
)

type PathDoc struct {
	Desc   string
	Type   domain.PathType
	Method string
	Code   string
}

func main() {
	server.PreServerBootstrap(func(rail core.Rail) error {
		if err := scheduleTasks(rail); err != nil { // schedule cron jobs
			return err
		}
		registerWebEndpoints(rail)     // register http server endpoints
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

func registerWebEndpoints(ec core.Rail) {
	server.PostServerBootstrapped(func(c core.Rail) error {
		return domain.CreateResourceIfNotExist(ec, domain.CreateResReq{
			Code: codeMngResources,
			Name: nameMngReesources,
		}, common.NilUser())
	})

	/*
		------------------------------

		public endpoints

		-------------------------------
	*/
	urlpath := "/open/api/resource/brief/user"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PUBLIC, Desc: "List resources of current user", Method: "GET"})
	server.Get(urlpath, web.ListAllResBriefsOfRole)

	urlpath = "/open/api/resource/brief/all"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PUBLIC, Desc: "List all resource brief info", Method: "GET"})
	server.Get(urlpath, web.ListAllResBriefs)

	urlpath = "/open/api/role/info"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PUBLIC, Desc: "Get role info", Method: "POST"})
	server.IPost(urlpath, web.GetRoleInfo)

	/*
		------------------------------

		protected endpoints

		-------------------------------
	*/
	urlpath = "/open/api/resource/add"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin add resource", Code: codeMngResources, Method: "POST"})
	server.IPost(urlpath, web.CreateResourceIfNotExist)

	urlpath = "/open/api/resource/remove"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin remove resource", Code: codeMngResources, Method: "POST"})
	server.IPost(urlpath, web.DeleteResource)

	urlpath = "/open/api/resource/brief/candidates"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "List all resource candidates for role", Code: codeMngResources,
		Method: "GET"})
	server.Get(urlpath, web.ListResourceCandidatesForRole)

	urlpath = "/open/api/resource/list"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list resources", Code: codeMngResources, Method: "POST"})
	server.IPost(urlpath, web.ListResources)

	urlpath = "/open/api/role/resource/add"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin add resource to role", Code: codeMngResources,
		Method: "POST"})
	server.IPost(urlpath, web.AddResToRoleIfNotExist)

	urlpath = "/open/api/role/resource/remove"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin remove resource from role", Code: codeMngResources,
		Method: "POST"})
	server.IPost(urlpath, web.RemoveResFromRole)

	urlpath = "/open/api/role/add"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin add role", Code: codeMngResources, Method: "POST"})
	server.IPost(urlpath, web.AddRole)

	urlpath = "/open/api/role/list"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list roles", Code: codeMngResources, Method: "POST"})
	server.IPost(urlpath, web.ListRoles)

	urlpath = "/open/api/role/brief/all"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list role brief info", Code: codeMngResources,
		Method: "GET"})
	server.Get(urlpath, web.ListAllRoleBriefs)

	urlpath = "/open/api/role/resource/list"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list resources of role", Code: codeMngResources,
		Method: "POST"})
	server.IPost(urlpath, web.ListRoleRes)

	urlpath = "/open/api/path/list"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list paths", Code: codeMngResources, Method: "POST"})
	server.IPost(urlpath, web.ListPaths)

	urlpath = "/open/api/path/resource/bind"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin bind resource to path", Code: codeMngResources,
		Method: "POST"})
	server.IPost(urlpath, web.BindPathRes)

	urlpath = "/open/api/path/resource/unbind"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin unbind resource and path", Code: codeMngResources,
		Method: "POST"})
	server.IPost(urlpath, web.UnbindPathRes)

	urlpath = "/open/api/path/delete"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin delete path", Code: codeMngResources, Method: "POST"})
	server.IPost(urlpath, web.DeletePath)

	urlpath = "/open/api/path/update"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin update path", Code: codeMngResources, Method: "POST"})
	server.IPost(urlpath, web.UpdatePath)

	// internal endpoints
	server.IPost("/remote/resource/add",
		func(c *gin.Context, rail core.Rail, req domain.CreateResReq) (any, error) {
			user := common.GetUser(rail)
			return nil, domain.CreateResourceIfNotExist(rail, req, user)
		})
	server.IPost("/remote/path/resource/access-test",
		func(c *gin.Context, rail core.Rail, req domain.TestResAccessReq) (any, error) {
			return domain.TestResourceAccess(rail, req)
		})
	server.IPost("/remote/path/add",
		func(c *gin.Context, rail core.Rail, req domain.CreatePathReq) (any, error) {
			user := common.GetUser(rail)
			return nil, domain.CreatePathIfNotExist(rail, req, user)
		})
	server.IPost("/remote/role/info",
		func(c *gin.Context, rail core.Rail, req domain.RoleInfoReq) (any, error) {
			return domain.GetRoleInfo(rail, req)
		})
}

func reportPathOnBootstrapped(ec core.Rail, url string, doc PathDoc) {
	server.PostServerBootstrapped(func(c core.Rail) error {
		ptype := doc.Type
		desc := doc.Desc
		resCode := doc.Code
		method := doc.Method

		r := domain.CreatePathReq{
			Type:    ptype,
			Desc:    desc,
			Method:  method,
			Group:   "goauth",
			Url:     "/goauth" + url,
			ResCode: resCode,
		}
		return domain.CreatePathIfNotExist(ec, r, common.NilUser())
	})
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
