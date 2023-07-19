package main

import (
	"fmt"
	"os"

	"github.com/curtisnewbie/goauth/domain"
	"github.com/curtisnewbie/goauth/web"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/curtisnewbie/gocommon/task"
	"github.com/gin-gonic/gin"
)

const CODE_MNG_RESOURCES = "manage-resources"
const NAME_MNG_RESOURCES = "Manage Resources Access"

type PathDoc struct {
	Desc   string
	Type   domain.PathType
	Method string
	Code   string
}

func main() {
	server.PreServerBootstrap(func(c common.ExecContext) error {
		scheduleTasks()         // schedule cron jobs
		registerWebEndpoints(c) // register http server endpoints
		return nil
	})

	server.BootstrapServer(os.Args) // bootstrap server
}

func registerWebEndpoints(ec common.ExecContext) {
	server.PostServerBootstrapped(func(c common.ExecContext) error {
		return domain.CreateResourceIfNotExist(ec, domain.CreateResReq{
			Code: CODE_MNG_RESOURCES,
			Name: NAME_MNG_RESOURCES,
		})
	})

	/*
		------------------------------

		public endpoints

		-------------------------------
	*/
	urlpath := server.OpenApiPath("/resource/brief/user")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PUBLIC, Desc: "List resources of current user", Method: "GET"})
	server.Get(urlpath, web.ListAllResBriefsOfRole)

	urlpath = server.OpenApiPath("/resource/brief/all")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PUBLIC, Desc: "List all resource brief info", Method: "GET"})
	server.Get(urlpath, web.ListAllResBriefs)

	urlpath = server.OpenApiPath("/role/info")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PUBLIC, Desc: "Get role info", Method: "POST"})
	server.IPost(urlpath, web.GetRoleInfo)

	/*
		------------------------------

		protected endpoints

		-------------------------------
	*/
	urlpath = server.OpenApiPath("/resource/add")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin add resource", Code: CODE_MNG_RESOURCES, Method: "POST"})
	server.IPost(urlpath, web.CreateResourceIfNotExist)

	urlpath = server.OpenApiPath("/resource/remove")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin remove resource", Code: CODE_MNG_RESOURCES, Method: "POST"})
	server.IPost(urlpath, web.DeleteResource)

	urlpath = server.OpenApiPath("/resource/brief/candidates")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "List all resource candidates for role", Code: CODE_MNG_RESOURCES,
		Method: "GET"})
	server.Get(urlpath, web.ListResourceCandidatesForRole)

	urlpath = server.OpenApiPath("/resource/list")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list resources", Code: CODE_MNG_RESOURCES, Method: "POST"})
	server.IPost(urlpath, web.ListResources)

	urlpath = server.OpenApiPath("/role/resource/add")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin add resource to role", Code: CODE_MNG_RESOURCES,
		Method: "POST"})
	server.IPost(urlpath, web.AddResToRoleIfNotExist)

	urlpath = server.OpenApiPath("/role/resource/remove")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin remove resource from role", Code: CODE_MNG_RESOURCES,
		Method: "POST"})
	server.IPost(urlpath, web.RemoveResFromRole)

	urlpath = server.OpenApiPath("/role/add")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin add role", Code: CODE_MNG_RESOURCES, Method: "POST"})
	server.IPost(urlpath, web.AddRole)

	urlpath = server.OpenApiPath("/role/list")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list roles", Code: CODE_MNG_RESOURCES, Method: "POST"})
	server.IPost(urlpath, web.ListRoles)

	urlpath = server.OpenApiPath("/role/brief/all")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list role brief info", Code: CODE_MNG_RESOURCES,
		Method: "GET"})
	server.Get(urlpath, web.ListAllRoleBriefs)

	urlpath = server.OpenApiPath("/role/resource/list")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list resources of role", Code: CODE_MNG_RESOURCES,
		Method: "POST"})
	server.IPost(urlpath, web.ListRoleRes)

	urlpath = server.OpenApiPath("/path/list")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list paths", Code: CODE_MNG_RESOURCES, Method: "POST"})
	server.IPost(urlpath, web.ListPaths)

	urlpath = server.OpenApiPath("/path/resource/bind")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin bind resource to path", Code: CODE_MNG_RESOURCES,
		Method: "POST"})
	server.IPost(urlpath, web.BindPathRes)

	urlpath = server.OpenApiPath("/path/resource/unbind")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin unbind resource and path", Code: CODE_MNG_RESOURCES,
		Method: "POST"})
	server.IPost(urlpath, web.UnbindPathRes)

	urlpath = server.OpenApiPath("/path/delete")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin delete path", Code: CODE_MNG_RESOURCES, Method: "POST"})
	server.IPost(urlpath, web.DeletePath)

	urlpath = server.OpenApiPath("/path/update")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin update path", Code: CODE_MNG_RESOURCES, Method: "POST"})
	server.IPost(urlpath, web.UpdatePath)

	// internal endpoints
	server.IPost(server.InternalApiPath("/resource/add"),
		func(c *gin.Context, ec common.ExecContext, req domain.CreateResReq) (any, error) {
			return nil, domain.CreateResourceIfNotExist(ec, req)
		})
	server.IPost(server.InternalApiPath("/path/resource/access-test"),
		func(c *gin.Context, ec common.ExecContext, req domain.TestResAccessReq) (any, error) {
			return domain.TestResourceAccess(ec, req)
		})
	server.IPost(server.InternalApiPath("/path/add"),
		func(c *gin.Context, ec common.ExecContext, req domain.CreatePathReq) (any, error) {
			return nil, domain.CreatePathIfNotExist(ec, req)
		})
	server.IPost(server.InternalApiPath("/role/info"),
		func(c *gin.Context, ec common.ExecContext, req domain.RoleInfoReq) (any, error) {
			return domain.GetRoleInfo(ec, req)
		})
}

func reportPathOnBootstrapped(ec common.ExecContext, url string, doc PathDoc) {
	server.PostServerBootstrapped(func(c common.ExecContext) error {
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
		return domain.CreatePathIfNotExist(ec, r)
	})
}

func scheduleTasks() {
	// distributed tasks
	task.ScheduleNamedDistributedTask("0 0/15 * * * *", "LoadRoleResCacheTask", func(ec common.ExecContext) error {
		return domain.LoadRoleResCache(ec)
	})
	task.ScheduleNamedDistributedTask("0 0/15 * * * *", "LoadPathResCacheTask", func(ec common.ExecContext) error {
		return domain.LoadPathResCache(ec)
	})

	// for the first time
	server.PostServerBootstrapped(func(c common.ExecContext) error {
		ec := common.EmptyExecContext()
		if e := domain.LoadRoleResCache(ec); e != nil {
			return fmt.Errorf("failed to load role resource, %v", e)
		}
		return nil
	})
	server.PostServerBootstrapped(func(c common.ExecContext) error {
		ec := common.EmptyExecContext()
		if e := domain.LoadPathResCache(ec); e != nil {
			return fmt.Errorf("failed to load path resource, %v", e)
		}
		return nil
	})
}
