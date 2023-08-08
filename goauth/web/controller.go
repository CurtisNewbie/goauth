package web

import (
	"github.com/curtisnewbie/goauth/domain"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/gin-gonic/gin"
)

func ListAllResBriefsOfRole(c *gin.Context, ec common.Rail) (any, error) {
	u := server.ExtractUser(c)
	if u.IsNil {
		return []domain.ResBrief{}, nil
	}
	return domain.ListAllResBriefsOfRole(ec, u.RoleNo)
}

func ListAllResBriefs(c *gin.Context, ec common.Rail) (any, error) {
	return domain.ListAllResBriefs(ec)
}

func GetRoleInfo(c *gin.Context, ec common.Rail, req domain.RoleInfoReq) (any, error) {
	return domain.GetRoleInfo(ec, req)
}

func CreateResourceIfNotExist(c *gin.Context, ec common.Rail, req domain.CreateResReq) (any, error) {
	user := server.ExtractUser(c)
	return nil, domain.CreateResourceIfNotExist(ec, req, user)
}

func DeleteResource(c *gin.Context, ec common.Rail, req domain.DeleteResourceReq) (any, error) {
	return nil, domain.DeleteResource(ec, req)
}

func ListResourceCandidatesForRole(c *gin.Context, ec common.Rail) (any, error) {
	roleNo := c.Query("roleNo")
	return domain.ListResourceCandidatesForRole(ec, roleNo)
}

func ListResources(c *gin.Context, ec common.Rail, req domain.ListResReq) (any, error) {
	return domain.ListResources(ec, req)
}

func AddResToRoleIfNotExist(c *gin.Context, ec common.Rail, req domain.AddRoleResReq) (any, error) {
	user := server.ExtractUser(c)
	return nil, domain.AddResToRoleIfNotExist(ec, req, user)
}

func RemoveResFromRole(c *gin.Context, ec common.Rail, req domain.RemoveRoleResReq) (any, error) {
	return nil, domain.RemoveResFromRole(ec, req)
}

func AddRole(c *gin.Context, ec common.Rail, req domain.AddRoleReq) (any, error) {
	user := server.ExtractUser(c)
	return nil, domain.AddRole(ec, req, user)
}

func ListRoles(c *gin.Context, ec common.Rail, req domain.ListRoleReq) (any, error) {
	return domain.ListRoles(ec, req)
}

func ListAllRoleBriefs(c *gin.Context, ec common.Rail) (any, error) {
	return domain.ListAllRoleBriefs(ec)
}

func ListRoleRes(c *gin.Context, ec common.Rail, req domain.ListRoleResReq) (any, error) {
	return domain.ListRoleRes(ec, req)
}

func ListPaths(c *gin.Context, ec common.Rail, req domain.ListPathReq) (any, error) {
	return domain.ListPaths(ec, req)
}

func BindPathRes(c *gin.Context, ec common.Rail, req domain.BindPathResReq) (any, error) {
	return nil, domain.BindPathRes(ec, req)
}

func UnbindPathRes(c *gin.Context, ec common.Rail, req domain.UnbindPathResReq) (any, error) {
	return nil, domain.UnbindPathRes(ec, req)
}

func DeletePath(c *gin.Context, ec common.Rail, req domain.DeletePathReq) (any, error) {
	return nil, domain.DeletePath(ec, req)
}

func UpdatePath(c *gin.Context, ec common.Rail, req domain.UpdatePathReq) (any, error) {
	return nil, domain.UpdatePath(ec, req)
}
