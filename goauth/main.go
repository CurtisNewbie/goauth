package main

import (
	"os"

	"github.com/curtisnewbie/goauth/domain"
	"github.com/curtisnewbie/goauth/web"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/gin-gonic/gin"
)

const CODE_MNG_RESOURCES = "manage-resources"
const NAME_MNG_RESOURCES = "Manage Resources Access"

type PathDoc struct {
	Desc string
	Type domain.PathType
	Code string
}

func main() {
	ec := common.EmptyExecContext()
	scheduleJobs()                         // schedule cron jobs
	registerWebEndpoints(ec)               // register http server endpoints
	server.DefaultBootstrapServer(os.Args) // bootstrap server
}

func registerWebEndpoints(ec common.ExecContext) {
	/*
		------------------------------

		public endpoints

		-------------------------------
	*/
	urlpath := server.OpenApiPath("/resource/brief/user")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PUBLIC, Desc: "List resources of current user"})
	server.Get(urlpath, web.ListAllResBriefsOfRole)

	urlpath = server.OpenApiPath("/resource/brief/all")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PUBLIC, Desc: "List all resource brief info"})
	server.Get(urlpath, web.ListAllResBriefs)

	urlpath = server.OpenApiPath("/role/info")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PUBLIC, Desc: "Get role info"})
	server.PostJ(urlpath, web.GetRoleInfo)

	/*
		------------------------------

		protected endpoints

		-------------------------------
	*/
	urlpath = server.OpenApiPath("/resource/add")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin add resource", Code: CODE_MNG_RESOURCES})
	server.PostJ(urlpath, web.CreateResourceIfNotExist)

	urlpath = server.OpenApiPath("/resource/remove")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin remove resource", Code: CODE_MNG_RESOURCES})
	server.PostJ(urlpath, web.DeleteResource)

	urlpath = server.OpenApiPath("/resource/brief/candidates")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "List all resource candidates for role", Code: CODE_MNG_RESOURCES})
	server.Get(urlpath, web.ListResourceCandidatesForRole)

	urlpath = server.OpenApiPath("/resource/list")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list resources", Code: CODE_MNG_RESOURCES})
	server.PostJ(urlpath, web.ListResources)

	urlpath = server.OpenApiPath("/role/resource/add")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin add resource to role", Code: CODE_MNG_RESOURCES})
	server.PostJ(urlpath, web.AddResToRoleIfNotExist)

	urlpath = server.OpenApiPath("/role/resource/remove")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin remove resource from role", Code: CODE_MNG_RESOURCES})
	server.PostJ(urlpath, web.RemoveResFromRole)

	urlpath = server.OpenApiPath("/role/add")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin add role", Code: CODE_MNG_RESOURCES})
	server.PostJ(urlpath, web.AddRole)

	urlpath = server.OpenApiPath("/role/list")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list roles", Code: CODE_MNG_RESOURCES})
	server.PostJ(urlpath, web.ListRoles)

	urlpath = server.OpenApiPath("/role/brief/all")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list role brief info", Code: CODE_MNG_RESOURCES})
	server.Get(urlpath, web.ListAllRoleBriefs)

	urlpath = server.OpenApiPath("/role/resource/list")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list resources of role", Code: CODE_MNG_RESOURCES})
	server.PostJ(urlpath, web.ListRoleRes)

	urlpath = server.OpenApiPath("/path/list")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list paths", Code: CODE_MNG_RESOURCES})
	server.PostJ(urlpath, web.ListPaths)

	urlpath = server.OpenApiPath("/path/resource/bind")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin bind resource to path", Code: CODE_MNG_RESOURCES})
	server.PostJ(urlpath, web.BindPathRes)

	urlpath = server.OpenApiPath("/path/resource/unbind")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin unbind resource and path", Code: CODE_MNG_RESOURCES})
	server.PostJ(urlpath, web.UnbindPathRes)

	urlpath = server.OpenApiPath("/path/delete")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin delete path", Code: CODE_MNG_RESOURCES})
	server.PostJ(urlpath, web.DeletePath)

	urlpath = server.OpenApiPath("/path/update")
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin update path", Code: CODE_MNG_RESOURCES})
	server.PostJ(urlpath, web.UpdatePath)

	/*
		Generate resource scripts for production environment, for internal use only

		curl -X POST "http://localhost:8081/internal/resource/script/generate" \
			-H 'content-type:application/json' \
			-d '{ "resCodes" : ["basic-user", "manage-users"]}' \
			-o output.sql
	*/
	urlpath = "/internal/resource/script/generate"
	server.RawPost(urlpath, web.GenResourceScript)

	// internal endpoints
	server.PostJ(server.InternalApiPath("/resource/add"),
		func(c *gin.Context, ec common.ExecContext, req domain.CreateResReq) (any, error) {
			return nil, domain.CreateResourceIfNotExist(ec, req)
		})
	server.PostJ(server.InternalApiPath("/path/resource/access-test"),
		func(c *gin.Context, ec common.ExecContext, req domain.TestResAccessReq) (any, error) {
			return domain.TestResourceAccess(ec, req)
		})
	server.PostJ(server.InternalApiPath("/path/add"),
		func(c *gin.Context, ec common.ExecContext, req domain.CreatePathReq) (any, error) {
			return nil, domain.CreatePathIfNotExist(ec, req)
		})
	server.PostJ(server.InternalApiPath("/path/batch/add"),
		func(c *gin.Context, ec common.ExecContext, req domain.BatchCreatePathReq) (any, error) {
			return nil, domain.BatchCreatePathIfNotExist(ec, req)
		})
	server.PostJ(server.InternalApiPath("/role/info"),
		func(c *gin.Context, ec common.ExecContext, req domain.RoleInfoReq) (any, error) {
			return domain.GetRoleInfo(ec, req)
		})
}

func reportPathOnBootstrapped(ec common.ExecContext, url string, doc PathDoc) {
	server.OnServerBootstrapped(func() {
		ptype := doc.Type
		desc := doc.Desc
		resCode := doc.Code

		r := domain.CreatePathReq{
			Type:    ptype,
			Desc:    desc,
			Group:   "goauth",
			Url:     "/goauth" + url,
			ResCode: resCode,
		}
		if e := domain.CreatePathIfNotExist(ec, r); e != nil {
			ec.Log.Fatal(e)
		}
	})
}

func scheduleJobs() {
	// jobs (with single instance only)
	common.ScheduleCron("0 0/15 * * * *", func() {
		ec := common.EmptyExecContext()
		if e := domain.LoadRoleResCache(ec); e != nil {
			ec.Log.Errorf("Failed to load role resource, %v", e)
		}
	})
	common.ScheduleCron("0 0/15 * * * *", func() {
		ec := common.EmptyExecContext()
		if e := domain.LoadPathResCache(ec); e != nil {
			ec.Log.Errorf("Failed to load path resource, %v", e)
		}
	})

	// for the first time
	server.OnServerBootstrapped(func() {
		ec := common.EmptyExecContext()
		if e := domain.LoadRoleResCache(ec); e != nil {
			ec.Log.Errorf("Failed to load role resource, %v", e)
		}
	})
	server.OnServerBootstrapped(func() {
		ec := common.EmptyExecContext()
		if e := domain.LoadPathResCache(ec); e != nil {
			ec.Log.Errorf("Failed to load path resource, %v", e)
		}
	})
}
