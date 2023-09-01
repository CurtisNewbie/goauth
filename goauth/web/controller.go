package web

import (
	"github.com/curtisnewbie/goauth/domain"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/miso/core"
	"github.com/curtisnewbie/miso/server"
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

func RegisterWebEndpoints(ec core.Rail) {
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
	server.Get(urlpath, ListAllResBriefsOfRole)

	urlpath = "/open/api/resource/brief/all"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PUBLIC, Desc: "List all resource brief info", Method: "GET"})
	server.Get(urlpath, ListAllResBriefs)

	urlpath = "/open/api/role/info"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PUBLIC, Desc: "Get role info", Method: "POST"})
	server.IPost(urlpath, GetRoleInfo)

	/*
		------------------------------

		protected endpoints

		-------------------------------
	*/
	urlpath = "/open/api/resource/add"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin add resource", Code: codeMngResources, Method: "POST"})
	server.IPost(urlpath, CreateResourceIfNotExist)

	urlpath = "/open/api/resource/remove"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin remove resource", Code: codeMngResources, Method: "POST"})
	server.IPost(urlpath, DeleteResource)

	urlpath = "/open/api/resource/brief/candidates"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "List all resource candidates for role", Code: codeMngResources,
		Method: "GET"})
	server.Get(urlpath, ListResourceCandidatesForRole)

	urlpath = "/open/api/resource/list"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list resources", Code: codeMngResources, Method: "POST"})
	server.IPost(urlpath, ListResources)

	urlpath = "/open/api/role/resource/add"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin add resource to role", Code: codeMngResources,
		Method: "POST"})
	server.IPost(urlpath, AddResToRoleIfNotExist)

	urlpath = "/open/api/role/resource/remove"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin remove resource from role", Code: codeMngResources,
		Method: "POST"})
	server.IPost(urlpath, RemoveResFromRole)

	urlpath = "/open/api/role/add"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin add role", Code: codeMngResources, Method: "POST"})
	server.IPost(urlpath, AddRole)

	urlpath = "/open/api/role/list"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list roles", Code: codeMngResources, Method: "POST"})
	server.IPost(urlpath, ListRoles)

	urlpath = "/open/api/role/brief/all"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list role brief info", Code: codeMngResources,
		Method: "GET"})
	server.Get(urlpath, ListAllRoleBriefs)

	urlpath = "/open/api/role/resource/list"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list resources of role", Code: codeMngResources,
		Method: "POST"})
	server.IPost(urlpath, ListRoleRes)

	urlpath = "/open/api/path/list"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin list paths", Code: codeMngResources, Method: "POST"})
	server.IPost(urlpath, ListPaths)

	urlpath = "/open/api/path/resource/bind"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin bind resource to path", Code: codeMngResources,
		Method: "POST"})
	server.IPost(urlpath, BindPathRes)

	urlpath = "/open/api/path/resource/unbind"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin unbind resource and path", Code: codeMngResources,
		Method: "POST"})
	server.IPost(urlpath, UnbindPathRes)

	urlpath = "/open/api/path/delete"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin delete path", Code: codeMngResources, Method: "POST"})
	server.IPost(urlpath, DeletePath)

	urlpath = "/open/api/path/update"
	reportPathOnBootstrapped(ec, urlpath, PathDoc{Type: domain.PT_PROTECTED, Desc: "Admin update path", Code: codeMngResources, Method: "POST"})
	server.IPost(urlpath, UpdatePath)

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

func ListAllResBriefsOfRole(c *gin.Context, ec core.Rail) (any, error) {
	u := common.GetUser(ec)
	if u.IsNil {
		return []domain.ResBrief{}, nil
	}
	return domain.ListAllResBriefsOfRole(ec, u.RoleNo)
}

func ListAllResBriefs(c *gin.Context, ec core.Rail) (any, error) {
	return domain.ListAllResBriefs(ec)
}

func GetRoleInfo(c *gin.Context, ec core.Rail, req domain.RoleInfoReq) (any, error) {
	return domain.GetRoleInfo(ec, req)
}

func CreateResourceIfNotExist(c *gin.Context, ec core.Rail, req domain.CreateResReq) (any, error) {
	user := common.GetUser(ec)
	return nil, domain.CreateResourceIfNotExist(ec, req, user)
}

func DeleteResource(c *gin.Context, ec core.Rail, req domain.DeleteResourceReq) (any, error) {
	return nil, domain.DeleteResource(ec, req)
}

func ListResourceCandidatesForRole(c *gin.Context, ec core.Rail) (any, error) {
	roleNo := c.Query("roleNo")
	return domain.ListResourceCandidatesForRole(ec, roleNo)
}

func ListResources(c *gin.Context, ec core.Rail, req domain.ListResReq) (any, error) {
	return domain.ListResources(ec, req)
}

func AddResToRoleIfNotExist(c *gin.Context, ec core.Rail, req domain.AddRoleResReq) (any, error) {
	user := common.GetUser(ec)
	return nil, domain.AddResToRoleIfNotExist(ec, req, user)
}

func RemoveResFromRole(c *gin.Context, ec core.Rail, req domain.RemoveRoleResReq) (any, error) {
	return nil, domain.RemoveResFromRole(ec, req)
}

func AddRole(c *gin.Context, ec core.Rail, req domain.AddRoleReq) (any, error) {
	user := common.GetUser(ec)
	return nil, domain.AddRole(ec, req, user)
}

func ListRoles(c *gin.Context, ec core.Rail, req domain.ListRoleReq) (any, error) {
	return domain.ListRoles(ec, req)
}

func ListAllRoleBriefs(c *gin.Context, ec core.Rail) (any, error) {
	return domain.ListAllRoleBriefs(ec)
}

func ListRoleRes(c *gin.Context, ec core.Rail, req domain.ListRoleResReq) (any, error) {
	return domain.ListRoleRes(ec, req)
}

func ListPaths(c *gin.Context, ec core.Rail, req domain.ListPathReq) (any, error) {
	return domain.ListPaths(ec, req)
}

func BindPathRes(c *gin.Context, ec core.Rail, req domain.BindPathResReq) (any, error) {
	return nil, domain.BindPathRes(ec, req)
}

func UnbindPathRes(c *gin.Context, ec core.Rail, req domain.UnbindPathResReq) (any, error) {
	return nil, domain.UnbindPathRes(ec, req)
}

func DeletePath(c *gin.Context, ec core.Rail, req domain.DeletePathReq) (any, error) {
	return nil, domain.DeletePath(ec, req)
}

func UpdatePath(c *gin.Context, ec core.Rail, req domain.UpdatePathReq) (any, error) {
	return nil, domain.UpdatePath(ec, req)
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
