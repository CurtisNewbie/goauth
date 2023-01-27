package main

import (
	"os"

	"github.com/curtisnewbie/goauth/domain"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/gin-gonic/gin"
)

func main() {
	// jobs
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

	// open-api endpoints
	server.PostJ(server.OpenApiPath("/role/resource/add"),
		func(c *gin.Context, ec common.ExecContext, req domain.AddRoleResReq) (any, error) {
			return nil, domain.AddResToRole(ec, req)
		})
	server.PostJ(server.OpenApiPath("/role/resource/remove"),
		func(c *gin.Context, ec common.ExecContext, req domain.RemoveRoleResReq) (any, error) {
			return nil, domain.RemoveResFromRole(ec, req)
		})
	server.PostJ(server.OpenApiPath("/role/add"),
		func(c *gin.Context, ec common.ExecContext, req domain.AddRoleReq) (any, error) {
			return nil, domain.AddRole(ec, req)
		})
	server.PostJ(server.OpenApiPath("/role/list"),
		func(c *gin.Context, ec common.ExecContext, req domain.ListRoleReq) (any, error) {
			return domain.ListRoles(ec, req)
		})
	server.PostJ(server.OpenApiPath("/role/resource/list"),
		func(c *gin.Context, ec common.ExecContext, req domain.ListRoleResReq) (any, error) {
			return domain.ListRoleRes(ec, req)
		})
	server.PostJ(server.OpenApiPath("/path/list"),
		func(c *gin.Context, ec common.ExecContext, req domain.ListPathReq) (any, error) {
			return domain.ListPaths(ec, req)
		})
	server.PostJ(server.OpenApiPath("/path/resource/bind"),
		func(c *gin.Context, ec common.ExecContext, req domain.BindPathResReq) (any, error) {
			return nil, domain.BindPathRes(ec, req)
		})
	server.PostJ(server.OpenApiPath("/path/resource/unbind"),
		func(c *gin.Context, ec common.ExecContext, req domain.UnbindPathResReq) (any, error) {
			return nil, domain.UnbindPathRes(ec, req)
		})

	// internal endpoints
	server.PostJ(server.InternalApiPath("/path/resource/access-test"),
		func(c *gin.Context, ec common.ExecContext, req domain.TestResAccessReq) (any, error) {
			return nil, domain.TestResourceAccess(ec, req)
		})

	server.DefaultBootstrapServer(os.Args)
}