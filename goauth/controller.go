package goauth

import (
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/goauth"
	"github.com/curtisnewbie/miso/miso"
	"github.com/gin-gonic/gin"
)

const (
	ResourceManageResources = "manage-resources"
)

var (
	resourceAccessCheckHisto = miso.NewPromHisto("goauth_resource_access_check_duration")
)

type PathDoc struct {
	Desc   string
	Type   PathType
	Method string
	Code   string
}

func RegisterWebEndpoints(rail miso.Rail) error {

	goauth.ReportOnBoostrapped(rail, []goauth.AddResourceReq{
		{Code: ResourceManageResources, Name: "Manage Resources Access"},
	})

	miso.BaseRoute("/open/api/resource").Group(
		miso.IPost("/add", CreateResourceIfNotExistEp).
			Desc("Admin add resource").
			Resource(ResourceManageResources),

		miso.IPost("/remove", DeleteResourceEp).
			Desc("Admin remove resource").
			Resource(ResourceManageResources),

		miso.Get("/brief/candidates", ListResourceCandidatesForRoleEp).
			Desc("List all resource candidates for role").
			Resource(ResourceManageResources),

		miso.IPost("/list", ListResourcesEp).
			Desc("Admin list resources").
			Resource(ResourceManageResources),

		miso.Get("/brief/user", ListAllResBriefsOfRoleEp).
			Desc("List resources of current user").
			Public(),

		miso.Get("/brief/all", ListAllResBriefsEp).
			Desc("List all resource brief info").
			Public(),
	)

	miso.BaseRoute("/open/api/role").Group(
		miso.IPost("/resource/add", AddResToRoleIfNotExistEp).
			Desc("Admin add resource to role").
			Resource(ResourceManageResources),

		miso.IPost("/resource/remove", RemoveResFromRoleEp).
			Desc("Admin remove resource from role").
			Resource(ResourceManageResources),

		miso.IPost("/add", AddRoleEp).
			Desc("Admin add role").
			Resource(ResourceManageResources),

		miso.IPost("/list", ListRolesEp).
			Desc("Admin list roles").
			Resource(ResourceManageResources),

		miso.Get("/brief/all", ListAllRoleBriefsEp).
			Desc("Admin list role brief info").
			Resource(ResourceManageResources),

		miso.IPost("/resource/list", ListRoleResEp).
			Desc("Admin list resources of role").
			Resource(ResourceManageResources),

		miso.IPost("/info", GetRoleInfoEp).
			Desc("Get role info").
			Public(),
	)

	miso.BaseRoute("/open/api/path").Group(
		miso.IPost("/list", ListPathsEp).
			Desc("Admin list paths").
			Resource(ResourceManageResources),

		miso.IPost("/resource/bind", BindPathResEp).
			Desc("Admin bind resource to path").
			Resource(ResourceManageResources),

		miso.IPost("/resource/unbind", UnbindPathResEp).
			Desc("Admin unbind resource and path").
			Resource(ResourceManageResources),

		miso.IPost("/delete", DeletePathEp).
			Desc("Admin delete path").
			Resource(ResourceManageResources),

		miso.IPost("/update", UpdatePathEp).
			Desc("Admin update path").
			Resource(ResourceManageResources),
	)

	// internal endpoints
	miso.BaseRoute("/remote").Group(

		miso.IPost("/resource/add",
			func(c *gin.Context, rail miso.Rail, req CreateResReq) (any, error) {
				user := common.GetUser(rail)
				return nil, CreateResourceIfNotExist(rail, req, user)
			}),
		miso.IPost("/path/resource/access-test",
			func(c *gin.Context, rail miso.Rail, req TestResAccessReq) (any, error) {
				timer := miso.NewHistTimer(resourceAccessCheckHisto)
				defer timer.ObserveDuration()

				return TestResourceAccess(rail, req)
			}),
		miso.IPost("/path/add",
			func(c *gin.Context, rail miso.Rail, req CreatePathReq) (any, error) {
				user := common.GetUser(rail)
				return nil, CreatePathIfNotExist(rail, req, user)
			}),
		miso.IPost("/role/info",
			func(c *gin.Context, rail miso.Rail, req RoleInfoReq) (any, error) {
				return GetRoleInfo(rail, req)
			}),
	)
	return nil
}

func ListAllResBriefsOfRoleEp(c *gin.Context, ec miso.Rail) (any, error) {
	u := common.GetUser(ec)
	if u.IsNil {
		return []ResBrief{}, nil
	}
	return ListAllResBriefsOfRole(ec, u.RoleNo)
}

func ListAllResBriefsEp(c *gin.Context, ec miso.Rail) (any, error) {
	return ListAllResBriefs(ec)
}

func GetRoleInfoEp(c *gin.Context, ec miso.Rail, req RoleInfoReq) (any, error) {
	return GetRoleInfo(ec, req)
}

func CreateResourceIfNotExistEp(c *gin.Context, ec miso.Rail, req CreateResReq) (any, error) {
	user := common.GetUser(ec)
	return nil, CreateResourceIfNotExist(ec, req, user)
}

func DeleteResourceEp(c *gin.Context, ec miso.Rail, req DeleteResourceReq) (any, error) {
	return nil, DeleteResource(ec, req)
}

func ListResourceCandidatesForRoleEp(c *gin.Context, ec miso.Rail) (any, error) {
	roleNo := c.Query("roleNo")
	return ListResourceCandidatesForRole(ec, roleNo)
}

func ListResourcesEp(c *gin.Context, ec miso.Rail, req ListResReq) (any, error) {
	return ListResources(ec, req)
}

func AddResToRoleIfNotExistEp(c *gin.Context, ec miso.Rail, req AddRoleResReq) (any, error) {
	user := common.GetUser(ec)
	return nil, AddResToRoleIfNotExist(ec, req, user)
}

func RemoveResFromRoleEp(c *gin.Context, ec miso.Rail, req RemoveRoleResReq) (any, error) {
	return nil, RemoveResFromRole(ec, req)
}

func AddRoleEp(c *gin.Context, ec miso.Rail, req AddRoleReq) (any, error) {
	user := common.GetUser(ec)
	return nil, AddRole(ec, req, user)
}

func ListRolesEp(c *gin.Context, ec miso.Rail, req ListRoleReq) (any, error) {
	return ListRoles(ec, req)
}

func ListAllRoleBriefsEp(c *gin.Context, ec miso.Rail) (any, error) {
	return ListAllRoleBriefs(ec)
}

func ListRoleResEp(c *gin.Context, ec miso.Rail, req ListRoleResReq) (any, error) {
	return ListRoleRes(ec, req)
}

func ListPathsEp(c *gin.Context, ec miso.Rail, req ListPathReq) (any, error) {
	return ListPaths(ec, req)
}

func BindPathResEp(c *gin.Context, ec miso.Rail, req BindPathResReq) (any, error) {
	return nil, BindPathRes(ec, req)
}

func UnbindPathResEp(c *gin.Context, ec miso.Rail, req UnbindPathResReq) (any, error) {
	return nil, UnbindPathRes(ec, req)
}

func DeletePathEp(c *gin.Context, ec miso.Rail, req DeletePathReq) (any, error) {
	return nil, DeletePath(ec, req)
}

func UpdatePathEp(c *gin.Context, ec miso.Rail, req UpdatePathReq) (any, error) {
	return nil, UpdatePath(ec, req)
}
