package web

import (
	"github.com/curtisnewbie/goauth/domain"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/goauth"
	"github.com/curtisnewbie/miso/miso"
	"github.com/gin-gonic/gin"
)

const (
	codeMngResources  = "manage-resources"
	nameMngReesources = "Manage Resources Access"
)

type PathDoc struct {
	Desc   string
	Type   domain.PathType
	Method string
	Code   string
}

func RegisterWebEndpoints(rail miso.Rail) {

	goauth.ReportPathsOnBootstrapped(rail)
	goauth.ReportResourcesOnBootstrapped(rail, []goauth.AddResourceReq{
		{Code: codeMngResources, Name: nameMngReesources},
	})

	/*
		------------------------------

		public endpoints

		-------------------------------
	*/
	miso.BaseRoute("/open/api").Group(
		miso.Get("/resource/brief/user", ListAllResBriefsOfRole, goauth.Public("List resources of current user")),
		miso.Get("/resource/brief/all", ListAllResBriefs, goauth.Public("List all resource brief info")),
		miso.IPost("/open/api/role/info", GetRoleInfo, goauth.Public("Get role info")),
	)

	/*
		------------------------------

		protected endpoints

		-------------------------------
	*/
	miso.BaseRoute("/open/api").Group(
		miso.IPost("/resource/add", CreateResourceIfNotExist, goauth.Protected("Admin add resource", codeMngResources)),
		miso.IPost("/resource/remove", DeleteResource, goauth.Protected("Admin remove resource", codeMngResources)),
		miso.Get("/resource/brief/candidates", ListResourceCandidatesForRole, goauth.Protected("List all resource candidates for role", codeMngResources)),
		miso.IPost("/resource/list", ListResources, goauth.Protected("Admin list resources", codeMngResources)),
		miso.IPost("/role/resource/add", AddResToRoleIfNotExist, goauth.Protected("Admin add resource to role", codeMngResources)),
		miso.IPost("/role/resource/remove", RemoveResFromRole, goauth.Protected("Admin remove resource from role", codeMngResources)),
		miso.IPost("/role/add", AddRole, goauth.Protected("Admin add role", codeMngResources)),
		miso.IPost("/role/list", ListRoles, goauth.Protected("Admin list roles", codeMngResources)),
		miso.Get("/role/brief/all", ListAllRoleBriefs, goauth.Protected("Admin list role brief info", codeMngResources)),
		miso.IPost("/role/resource/list", ListRoleRes, goauth.Protected("Admin list resources of role", codeMngResources)),
		miso.IPost("/path/list", ListPaths, goauth.Protected("Admin list paths", codeMngResources)),
		miso.IPost("/path/resource/bind", BindPathRes, goauth.Protected("Admin bind resource to path", codeMngResources)),
		miso.IPost("/path/resource/unbind", UnbindPathRes, goauth.Protected("Admin unbind resource and path", codeMngResources)),
		miso.IPost("/path/delete", DeletePath, goauth.Protected("Admin delete path", codeMngResources)),
		miso.IPost("/path/update", UpdatePath, goauth.Protected("Admin update path", codeMngResources)),
	)

	// internal endpoints
	miso.BaseRoute("/remote").Group(
		miso.IPost("/resource/add",
			func(c *gin.Context, rail miso.Rail, req domain.CreateResReq) (any, error) {
				user := common.GetUser(rail)
				return nil, domain.CreateResourceIfNotExist(rail, req, user)
			}),
		miso.IPost("/path/resource/access-test",
			func(c *gin.Context, rail miso.Rail, req domain.TestResAccessReq) (any, error) {
				return domain.TestResourceAccess(rail, req)
			}),
		miso.IPost("/path/add",
			func(c *gin.Context, rail miso.Rail, req domain.CreatePathReq) (any, error) {
				user := common.GetUser(rail)
				return nil, domain.CreatePathIfNotExist(rail, req, user)
			}),
		miso.IPost("/role/info",
			func(c *gin.Context, rail miso.Rail, req domain.RoleInfoReq) (any, error) {
				return domain.GetRoleInfo(rail, req)
			}),
	)
}

func ListAllResBriefsOfRole(c *gin.Context, ec miso.Rail) (any, error) {
	u := common.GetUser(ec)
	if u.IsNil {
		return []domain.ResBrief{}, nil
	}
	return domain.ListAllResBriefsOfRole(ec, u.RoleNo)
}

func ListAllResBriefs(c *gin.Context, ec miso.Rail) (any, error) {
	return domain.ListAllResBriefs(ec)
}

func GetRoleInfo(c *gin.Context, ec miso.Rail, req domain.RoleInfoReq) (any, error) {
	return domain.GetRoleInfo(ec, req)
}

func CreateResourceIfNotExist(c *gin.Context, ec miso.Rail, req domain.CreateResReq) (any, error) {
	user := common.GetUser(ec)
	return nil, domain.CreateResourceIfNotExist(ec, req, user)
}

func DeleteResource(c *gin.Context, ec miso.Rail, req domain.DeleteResourceReq) (any, error) {
	return nil, domain.DeleteResource(ec, req)
}

func ListResourceCandidatesForRole(c *gin.Context, ec miso.Rail) (any, error) {
	roleNo := c.Query("roleNo")
	return domain.ListResourceCandidatesForRole(ec, roleNo)
}

func ListResources(c *gin.Context, ec miso.Rail, req domain.ListResReq) (any, error) {
	return domain.ListResources(ec, req)
}

func AddResToRoleIfNotExist(c *gin.Context, ec miso.Rail, req domain.AddRoleResReq) (any, error) {
	user := common.GetUser(ec)
	return nil, domain.AddResToRoleIfNotExist(ec, req, user)
}

func RemoveResFromRole(c *gin.Context, ec miso.Rail, req domain.RemoveRoleResReq) (any, error) {
	return nil, domain.RemoveResFromRole(ec, req)
}

func AddRole(c *gin.Context, ec miso.Rail, req domain.AddRoleReq) (any, error) {
	user := common.GetUser(ec)
	return nil, domain.AddRole(ec, req, user)
}

func ListRoles(c *gin.Context, ec miso.Rail, req domain.ListRoleReq) (any, error) {
	return domain.ListRoles(ec, req)
}

func ListAllRoleBriefs(c *gin.Context, ec miso.Rail) (any, error) {
	return domain.ListAllRoleBriefs(ec)
}

func ListRoleRes(c *gin.Context, ec miso.Rail, req domain.ListRoleResReq) (any, error) {
	return domain.ListRoleRes(ec, req)
}

func ListPaths(c *gin.Context, ec miso.Rail, req domain.ListPathReq) (any, error) {
	return domain.ListPaths(ec, req)
}

func BindPathRes(c *gin.Context, ec miso.Rail, req domain.BindPathResReq) (any, error) {
	return nil, domain.BindPathRes(ec, req)
}

func UnbindPathRes(c *gin.Context, ec miso.Rail, req domain.UnbindPathResReq) (any, error) {
	return nil, domain.UnbindPathRes(ec, req)
}

func DeletePath(c *gin.Context, ec miso.Rail, req domain.DeletePathReq) (any, error) {
	return nil, domain.DeletePath(ec, req)
}

func UpdatePath(c *gin.Context, ec miso.Rail, req domain.UpdatePathReq) (any, error) {
	return nil, domain.UpdatePath(ec, req)
}

func reportPathOnBootstrapped(ec miso.Rail, url string, doc PathDoc) {
	miso.PostServerBootstrapped(func(c miso.Rail) error {
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
