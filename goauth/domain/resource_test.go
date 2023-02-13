package domain

import (
	"fmt"
	"testing"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/mysql"
	"github.com/curtisnewbie/gocommon/redis"
	"github.com/curtisnewbie/gocommon/server"
)

func before(t *testing.T) {
	common.LoadConfigFromFile("../app-conf-dev.yml")
	if _, e := redis.InitRedisFromProp(); e != nil {
		t.Fatal(e)
	}
	if e := mysql.InitMySqlFromProp(); e != nil {
		t.Fatal(e)
	}
}

func TestUpdatePath(t *testing.T) {
	before(t)

	req := UpdatePathReq{
		PathNo: "path_578477630062592208429",
		Type:   PT_PUBLIC,
		Group:  "goauth",
	}
	e := UpdatePath(common.EmptyExecContext(), req)
	if e != nil {
		t.Fatal(e)
	}
}

func TestGetRoleInfo(t *testing.T) {
	before(t)

	req := RoleInfoReq{
		RoleNo: "role_554107924873216177918",
	}
	resp, e := GetRoleInfo(common.EmptyExecContext(), req)
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("%v", resp)
}

func TestCreatePathIfNotExist(t *testing.T) {
	before(t)

	req := CreatePathReq{
		Type:  PT_PROTECTED,
		Url:   "/goauth/open/api/role/resource/add",
		Group: "goauth",
	}
	e := CreatePathIfNotExist(common.EmptyExecContext(), req)
	if e != nil {
		t.Fatal(e)
	}
}

func TestDeletePath(t *testing.T) {
	before(t)

	req := DeletePathReq{
		PathNo: "path_555305367076864208429",
	}

	e := DeletePath(common.EmptyExecContext(), req)
	if e != nil {
		t.Fatal(e)
	}
}

func TestCreateRes(t *testing.T) {
	before(t)

	req := CreateResReq{
		Name: "GoAuth Test  ",
	}

	e := CreateResourceIfNotExist(common.EmptyExecContext(), req)
	if e != nil {
		t.Fatal(e)
	}
}

func TestBindPathRes(t *testing.T) {
	before(t)

	req := BindPathResReq{
		PathNo: "path_555326806016000208429",
		ResNo:  "res_555323073019904208429",
	}

	e := BindPathRes(common.EmptyExecContext(), req)
	if e != nil {
		t.Fatal(e)
	}
}

func TestUnbindPathRes(t *testing.T) {
	before(t)

	req := UnbindPathResReq{
		PathNo: "path_555326806016000208429",
	}

	e := UnbindPathRes(common.EmptyExecContext(), req)
	if e != nil {
		t.Fatal(e)
	}
}

func TestAddRole(t *testing.T) {
	before(t)

	req := AddRoleReq{
		Name: "Guest",
	}

	e := AddRole(common.EmptyExecContext(), req)
	if e != nil {
		t.Fatal(e)
	}
}

func TestAddResToRole(t *testing.T) {
	before(t)

	req := AddRoleResReq{
		RoleNo: "role_555329954676736208429",
		ResNo:  "res_555323073019904208429",
	}

	e := AddResToRoleIfNotExist(common.EmptyExecContext(), req)
	if e != nil {
		t.Fatal(e)
	}
}

func TestRemoveResFromRole(t *testing.T) {
	before(t)

	req := RemoveRoleResReq{
		RoleNo: "role_555329954676736208429",
		ResNo:  "res_555323073019904208429",
	}

	e := RemoveResFromRole(common.EmptyExecContext(), req)
	if e != nil {
		t.Fatal(e)
	}
}

func TestListRoleRes(t *testing.T) {
	before(t)

	p := common.Paging{
		Limit: 5,
		Page:  1,
	}
	req := ListRoleResReq{
		RoleNo: "role_555329954676736208429",
		Paging: p,
	}

	resp, e := ListRoleRes(common.EmptyExecContext(), req)
	if e != nil {
		t.Fatal(e)
	}

	if resp.Paging.Total < 1 {
		t.Fatal("total < 1")
	}

	t.Logf("%+v", resp)
}

func TestListRoles(t *testing.T) {
	before(t)

	p := common.Paging{
		Limit: 5,
		Page:  1,
	}
	req := ListRoleReq{
		Paging: p,
	}

	resp, e := ListRoles(common.EmptyExecContext(), req)
	if e != nil {
		t.Fatal(e)
	}

	if resp.Paging.Total < 1 {
		t.Fatal("total < 1")
	}

	t.Logf("%+v", resp)
}

func TestTestResourceAccess(t *testing.T) {
	before(t)

	ec := common.EmptyExecContext()
	LoadPathResCache(ec)
	LoadRoleResCache(ec)

	req := TestResAccessReq{
		RoleNo: "role_555329954676736208429",
		Url:    "/goauth/open/api/role/resource/add",
	}

	r, e := TestResourceAccess(ec, req)
	if e != nil {
		t.Fatal(e)
	}
	if !r.Valid {
		t.Fatal("should be valid")
	}
}

func TestGenInitialPathRoleRes(t *testing.T) {
	roleNo := "role_554107924873216177918"
	roleName := "Administrator"
	paths := []namedPath{
		{
			url:     "goauth" + server.OpenApiPath("/resource/add"),
			resName: "Add Resource",
		},
		{
			url:     "goauth" + server.OpenApiPath("/role/resource/add"),
			resName: "Add Resource To Role",
		},
		{
			url:     "goauth" + server.OpenApiPath("/role/resource/remove"),
			resName: "Remove Resource From Role",
		},
		{
			url:     "goauth" + server.OpenApiPath("/role/add"),
			resName: "Add New Role",
		},
		{
			url:     "goauth" + server.OpenApiPath("/role/list"),
			resName: "List Roles",
		},
		{
			url:     "goauth" + server.OpenApiPath("/role/resource/list"),
			resName: "List Resources of Role",
		},
		{
			url:     "goauth" + server.OpenApiPath("/path/list"),
			resName: "List Paths",
		},
		{
			url:     "goauth" + server.OpenApiPath("/path/resource/bind"),
			resName: "Bind Path to Resource",
		},
		{
			url:     "goauth" + server.OpenApiPath("/path/resource/unbind"),
			resName: "Unbind Path and Resource",
		},
		{
			url:     "goauth" + server.OpenApiPath("/path/delete"),
			resName: "Delete Path",
		},
		{
			url:     "goauth" + server.OpenApiPath("/path/add"),
			resName: "Add Path",
		},
		{
			url:     "goauth" + server.OpenApiPath("/role/info"),
			resName: "Fetch Role Info",
		},
		{
			url:     "goauth" + server.OpenApiPath("/path/update"),
			resName: "Update Path Info",
		},
	}

	initsql := fmt.Sprintf("INSERT INTO role(role_no, name) VALUES ('%s', '%s');", roleNo, roleName)
	for i, p := range paths {
		p.url = preprocessUrl(p.url)
		p.pathNo = common.GenIdP("path_")
		p.resNo = common.GenIdP("res_")
		paths[i] = p
	}

	initsql += "\n\nINSERT INTO resource(res_no, name) VALUES"
	for i, p := range paths {
		if i > 0 {
			initsql += ","
		}
		initsql += fmt.Sprintf("\n  ('%s', '%s')", p.resNo, p.resName)
	}
	initsql += ";"

	initsql += "\n\nINSERT INTO role_resource(role_no, res_no) VALUES"
	for i, p := range paths {
		if i > 0 {
			initsql += ","
		}
		initsql += fmt.Sprintf("\n  ('%s', '%s')", roleNo, p.resNo)
	}
	initsql += ";"

	initsql += "\n\nINSERT INTO path(path_no, url, ptype, res_no, pgroup) VALUES"
	for i, p := range paths {
		if i > 0 {
			initsql += ","
		}
		initsql += fmt.Sprintf("\n  ('%s', '%s', '%s', '%s', 'goauth')", p.pathNo, p.url, PT_PROTECTED, p.resNo)
	}
	initsql += ";"

	t.Log("\n\n" + initsql + "\n\n")
}

type namedPath struct {
	url     string
	resName string
	resNo   string
	pathNo  string
}

func TestGeneratedInitScript(t *testing.T) {
	before(t)

	ec := common.EmptyExecContext()
	LoadPathResCache(ec)
	LoadRoleResCache(ec)

	paths := []string{
		"goauth" + server.OpenApiPath("/resource/add"),
		"goauth" + server.OpenApiPath("/role/resource/add"),
		"goauth" + server.OpenApiPath("/role/resource/remove"),
		"goauth" + server.OpenApiPath("/role/add"),
		"goauth" + server.OpenApiPath("/role/list"),
		"goauth" + server.OpenApiPath("/role/resource/list"),
		"goauth" + server.OpenApiPath("/path/list"),
		"goauth" + server.OpenApiPath("/path/resource/bind"),
		"goauth" + server.OpenApiPath("/path/resource/unbind"),
		"goauth" + server.OpenApiPath("/path/delete"),
		"goauth" + server.OpenApiPath("/path/add"),
		"goauth" + server.OpenApiPath("/role/info"),
		"goauth" + server.OpenApiPath("/path/update"),
	}

	for _, p := range paths {
		r, e := TestResourceAccess(ec, TestResAccessReq{
			RoleNo: "role_554107924873216177918",
			Url:    p,
		})
		if e != nil {
			t.Fatal(e)
		}
		if !r.Valid {
			t.Fatalf("should be valid, url: '%s'", p)
		}
	}

}
