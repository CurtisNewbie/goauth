package main

import (
	"os"
	"strings"

	"github.com/curtisnewbie/goauth/domain"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/gin-gonic/gin"
)

type PathDoc struct {
	Desc string
	Type domain.PathType
}

func main() {
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

	// path doc
	pathDocs := map[string]PathDoc{}

	/*
		open-api endpoints
	*/

	/*
		public endpoints
	*/
	url := server.OpenApiPath("/resource/brief/user")
	pathDocs[url] = PathDoc{Type: domain.PT_PUBLIC, Desc: "List resources of current user"}
	server.Get(url, func(c *gin.Context, ec common.ExecContext) (any, error) {
		if !ec.Authenticated() {
			return []domain.ResBrief{}, nil
		}
		return domain.ListAllResBriefsOfRole(ec, ec.User.RoleNo)
	})

	/*
		protected endpoints
	*/
	url = server.OpenApiPath("/resource/brief/all")
	pathDocs[url] = PathDoc{Type: domain.PT_PUBLIC, Desc: "List all resource brief info"}
	server.Get(url, func(c *gin.Context, ec common.ExecContext) (any, error) {
		return domain.ListAllResBriefs(ec)
	})

	url = server.OpenApiPath("/resource/add")
	pathDocs[url] = PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin add resource"}
	server.PostJ(url, func(c *gin.Context, ec common.ExecContext, req domain.CreateResReq) (any, error) {
		return nil, domain.CreateResourceIfNotExist(ec, req)
	})

	url = server.OpenApiPath("/resource/remove")
	pathDocs[url] = PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin remove resource"}
	server.PostJ(url, func(c *gin.Context, ec common.ExecContext, req domain.DeleteResourceReq) (any, error) {
		return nil, domain.DeleteResource(ec, req)
	})

	url = server.OpenApiPath("/resource/brief/candidates")
	pathDocs[url] = PathDoc{Type: domain.PT_PROTECTED, Desc: "List all resource candidates for role"}
	server.Get(url, func(c *gin.Context, ec common.ExecContext) (any, error) {
		roleNo := c.Query("roleNo")
		return domain.ListResourceCandidatesForRole(ec, roleNo)
	})

	url = server.OpenApiPath("/resource/list")
	pathDocs[url] = PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list resources"}
	server.PostJ(url, func(c *gin.Context, ec common.ExecContext, req domain.ListResReq) (any, error) {
		return domain.ListResources(ec, req)
	})

	url = server.OpenApiPath("/role/resource/add")
	pathDocs[url] = PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin add resource to role"}
	server.PostJ(url, func(c *gin.Context, ec common.ExecContext, req domain.AddRoleResReq) (any, error) {
		return nil, domain.AddResToRoleIfNotExist(ec, req)
	})

	url = server.OpenApiPath("/role/resource/remove")
	pathDocs[url] = PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin remove resource from role"}
	server.PostJ(url, func(c *gin.Context, ec common.ExecContext, req domain.RemoveRoleResReq) (any, error) {
		return nil, domain.RemoveResFromRole(ec, req)
	})

	url = server.OpenApiPath("/role/add")
	pathDocs[url] = PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin add role"}
	server.PostJ(url, func(c *gin.Context, ec common.ExecContext, req domain.AddRoleReq) (any, error) {
		return nil, domain.AddRole(ec, req)
	})

	url = server.OpenApiPath("/role/list")
	pathDocs[url] = PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list roles"}
	server.PostJ(url, func(c *gin.Context, ec common.ExecContext, req domain.ListRoleReq) (any, error) {
		return domain.ListRoles(ec, req)
	})

	url = server.OpenApiPath("/role/brief/all")
	pathDocs[url] = PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list role brief info"}
	server.Get(url, func(c *gin.Context, ec common.ExecContext) (any, error) {
		return domain.ListAllRoleBriefs(ec)
	})

	url = server.OpenApiPath("/role/resource/list")
	pathDocs[url] = PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list resources of role"}
	server.PostJ(url, func(c *gin.Context, ec common.ExecContext, req domain.ListRoleResReq) (any, error) {
		return domain.ListRoleRes(ec, req)
	})

	url = server.OpenApiPath("/path/list")
	pathDocs[url] = PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list paths"}
	server.PostJ(url, func(c *gin.Context, ec common.ExecContext, req domain.ListPathReq) (any, error) {
		return domain.ListPaths(ec, req)
	})

	url = server.OpenApiPath("/path/resource/bind")
	pathDocs[url] = PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin bind resource to path"}
	server.PostJ(url, func(c *gin.Context, ec common.ExecContext, req domain.BindPathResReq) (any, error) {
		return nil, domain.BindPathRes(ec, req)
	})

	url = server.OpenApiPath("/path/resource/unbind")
	pathDocs[url] = PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin unbind resource and path"}
	server.PostJ(url, func(c *gin.Context, ec common.ExecContext, req domain.UnbindPathResReq) (any, error) {
		return nil, domain.UnbindPathRes(ec, req)
	})

	url = server.OpenApiPath("/path/delete")
	pathDocs[url] = PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin delete path"}
	server.PostJ(url, func(c *gin.Context, ec common.ExecContext, req domain.DeletePathReq) (any, error) {
		return nil, domain.DeletePath(ec, req)
	})

	url = server.OpenApiPath("/path/update")
	pathDocs[url] = PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin update path"}
	server.PostJ(url, func(c *gin.Context, ec common.ExecContext, req domain.UpdatePathReq) (any, error) {
		return nil, domain.UpdatePath(ec, req)
	})

	url = server.OpenApiPath("/role/info")
	pathDocs[url] = PathDoc{Type: domain.PT_PROTECTED, Desc: "Get role info"}
	server.PostJ(url, func(c *gin.Context, ec common.ExecContext, req domain.RoleInfoReq) (any, error) {
		return domain.GetRoleInfo(ec, req)
	})

	// internal endpoints
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
	server.Post(server.InternalApiPath("/path/cache/reload"),
		func(c *gin.Context, ec common.ExecContext) (any, error) {
			ec.Log.Info("Request to reload path cache")

			// asynchronously reload the cache of paths and resources
			go func() {
				if e := domain.LoadPathResCache(ec); e != nil {
					ec.Log.Errorf("Failed to load path resource, %v", e)
				}
			}()
			return nil, nil
		})

	// report paths (to itself) on bootstrap
	server.OnServerBootstrapped(func() {
		ec := common.EmptyExecContext()
		routes := server.GetRecordedServerRoutes()

		for _, u := range routes {
			if !strings.HasPrefix(u, "/remote") {

				ptype := domain.PT_PROTECTED
				desc := ""

				if doc, ok := pathDocs[u]; ok {
					ptype = doc.Type
					desc = doc.Desc
				}

				r := domain.CreatePathReq{
					Type:  ptype,
					Desc:  desc,
					Group: "goauth",
					Url:   "/goauth" + u,
				}
				domain.CreatePathIfNotExist(ec, r)
			}
		}
	})

	// bootstrap server
	server.DefaultBootstrapServer(os.Args)
}
