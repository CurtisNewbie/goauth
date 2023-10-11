package domain

import (
	"fmt"
	"testing"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/miso/miso"
)

func before(t *testing.T) {
	rail := miso.EmptyRail()
	miso.LoadConfigFromFile("../app-conf-dev.yml", rail)
	if _, e := miso.InitRedisFromProp(rail); e != nil {
		t.Fatal(e)
	}
	if e := miso.InitMySQLFromProp(); e != nil {
		t.Fatal(e)
	}
}

func TestUpdatePath(t *testing.T) {
	before(t)

	req := UpdatePathReq{
		PathNo: "path_578477630062592208429",
		Type:   PtPublic,
		Group:  "goauth",
	}
	e := UpdatePath(miso.EmptyRail(), req)
	if e != nil {
		t.Fatal(e)
	}
}

func TestGetRoleInfo(t *testing.T) {
	before(t)

	req := RoleInfoReq{
		RoleNo: "role_554107924873216177918",
	}
	resp, e := GetRoleInfo(miso.EmptyRail(), req)
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("%v", resp)
}

func TestCreatePathIfNotExist(t *testing.T) {
	before(t)

	req := CreatePathReq{
		Type:  PtProtected,
		Url:   "/goauth/open/api/role/resource/add",
		Group: "goauth",
	}
	e := CreatePathIfNotExist(miso.EmptyRail(), req, common.NilUser())
	if e != nil {
		t.Fatal(e)
	}
}

func TestDeletePath(t *testing.T) {
	before(t)

	req := DeletePathReq{
		PathNo: "path_555305367076864208429",
	}

	e := DeletePath(miso.EmptyRail(), req)
	if e != nil {
		t.Fatal(e)
	}
}

func TestCreateRes(t *testing.T) {
	before(t)

	req := CreateResReq{
		Name: "GoAuth Test  ",
	}

	e := CreateResourceIfNotExist(miso.EmptyRail(), req, common.NilUser())
	if e != nil {
		t.Fatal(e)
	}
}

func TestBindPathRes(t *testing.T) {
	before(t)

	req := BindPathResReq{
		PathNo:  "path_555326806016000208429",
		ResCode: "res_555323073019904208429",
	}

	e := BindPathRes(miso.EmptyRail(), req)
	if e != nil {
		t.Fatal(e)
	}
}

func TestPreprocessUrl(t *testing.T) {
	if v := preprocessUrl(""); v != "/" {
		t.Fatal(v)
	}

	if v := preprocessUrl("/"); v != "/" {
		t.Fatal(v)
	}

	if v := preprocessUrl("///"); v != "/" {
		t.Fatal(v)
	}

	if v := preprocessUrl("/goauth/test/path"); v != "/goauth/test/path" {
		t.Fatal(v)
	}

	if v := preprocessUrl("/goauth/test/path//"); v != "/goauth/test/path" {
		t.Fatal(v)
	}

	if v := preprocessUrl("goauth/test/path//"); v != "/goauth/test/path" {
		t.Fatal(v)
	}

	if v := preprocessUrl("goauth/test/path?abc=123"); v != "/goauth/test/path" {
		t.Fatal(v)
	}
}

func TestUnbindPathRes(t *testing.T) {
	before(t)

	req := UnbindPathResReq{
		PathNo: "path_555326806016000208429",
	}

	e := UnbindPathRes(miso.EmptyRail(), req)
	if e != nil {
		t.Fatal(e)
	}
}

func TestAddRole(t *testing.T) {
	before(t)

	req := AddRoleReq{
		Name: "Guest",
	}

	e := AddRole(miso.EmptyRail(), req, common.NilUser())
	if e != nil {
		t.Fatal(e)
	}
}

func TestAddResToRole(t *testing.T) {
	before(t)

	req := AddRoleResReq{
		RoleNo:  "role_555329954676736208429",
		ResCode: "res_555323073019904208429",
	}

	e := AddResToRoleIfNotExist(miso.EmptyRail(), req, common.NilUser())
	if e != nil {
		t.Fatal(e)
	}
}

func TestGenPathNo(t *testing.T) {
	pathNo := genPathNo("test", "/core/path/is/that/okay/if/i/amy/very", "GET")
	if pathNo == "" {
		t.Error("pathNo is empty")
		return
	}
	t.Log(pathNo)
}

func TestRemoveResFromRole(t *testing.T) {
	before(t)

	req := RemoveRoleResReq{
		RoleNo:  "role_555329954676736208429",
		ResCode: "res_555323073019904208429",
	}

	e := RemoveResFromRole(miso.EmptyRail(), req)
	if e != nil {
		t.Fatal(e)
	}
}

func TestListRoleRes(t *testing.T) {
	before(t)

	p := miso.Paging{
		Limit: 5,
		Page:  1,
	}
	req := ListRoleResReq{
		RoleNo: "role_555329954676736208429",
		Paging: p,
	}

	resp, e := ListRoleRes(miso.EmptyRail(), req)
	if e != nil {
		t.Fatal(e)
	}

	if resp.Paging.Total < 1 {
		t.Fatal("total < 1")
	}

	t.Logf("%+v", resp)
}

func TestListAllRoleBriefs(t *testing.T) {
	before(t)

	resp, e := ListAllRoleBriefs(miso.EmptyRail())
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("%+v", resp)
}

func TestListRoles(t *testing.T) {
	before(t)

	p := miso.Paging{
		Limit: 5,
		Page:  1,
	}
	req := ListRoleReq{
		Paging: p,
	}

	resp, e := ListRoles(miso.EmptyRail(), req)
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

	ec := miso.EmptyRail()
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
			resNo:   "res_578477630062593208429",
			pathNo:  "path_578477630062592208429",
			url:     "goauth" + miso.OpenApiPath("/resource/add"),
			resName: "Add Resource",
		},
		{
			resNo:   "res_578477630062595208429",
			pathNo:  "path_578477630062594208429",
			url:     "goauth" + miso.OpenApiPath("/role/resource/add"),
			resName: "Add Resource To Role",
		},
		{
			resNo:   "res_578477630062597208429",
			pathNo:  "path_578477630062596208429",
			url:     "goauth" + miso.OpenApiPath("/role/resource/remove"),
			resName: "Remove Resource From Role",
		},
		{
			resNo:   "res_578477630062599208429",
			pathNo:  "path_578477630062598208429",
			url:     "goauth" + miso.OpenApiPath("/role/add"),
			resName: "Add New Role",
		},
		{
			resNo:   "res_578477630062601208429",
			pathNo:  "path_578477630062600208429",
			url:     "goauth" + miso.OpenApiPath("/role/list"),
			resName: "List Roles",
		},
		{
			resNo:   "res_578477630062603208429",
			pathNo:  "path_578477630062602208429",
			url:     "goauth" + miso.OpenApiPath("/role/resource/list"),
			resName: "List Resources of Role",
		},
		{
			resNo:   "res_578477630062605208429",
			pathNo:  "path_578477630062604208429",
			url:     "goauth" + miso.OpenApiPath("/path/list"),
			resName: "List Paths",
		},
		{
			resNo:   "res_578477630062607208429",
			pathNo:  "path_578477630062606208429",
			url:     "goauth" + miso.OpenApiPath("/path/resource/bind"),
			resName: "Bind Path to Resource",
		},
		{
			resNo:   "res_578477630062609208429",
			pathNo:  "path_578477630062608208429",
			url:     "goauth" + miso.OpenApiPath("/path/resource/unbind"),
			resName: "Unbind Path and Resource",
		},
		{
			resNo:   "res_578477630062611208429",
			pathNo:  "path_578477630062610208429",
			url:     "goauth" + miso.OpenApiPath("/path/delete"),
			resName: "Delete Path",
		},
		{
			resNo:   "res_578477630062613208429",
			pathNo:  "path_578477630062612208429",
			url:     "goauth" + miso.OpenApiPath("/path/add"),
			resName: "Add Path",
		},
		{
			resNo:   "res_578477630062615208429",
			pathNo:  "path_578477630062614208429",
			url:     "goauth" + miso.OpenApiPath("/role/info"),
			resName: "Fetch Role Info",
		},
		{
			resNo:   "res_578477630062617208429",
			pathNo:  "path_578477630062616208429",
			url:     "goauth" + miso.OpenApiPath("/path/update"),
			resName: "Update Path Info",
		},
		{
			resNo:   "res_585463207870465208429",
			pathNo:  "path_585463207870464208429",
			url:     "goauth" + miso.OpenApiPath("/role/all"),
			resName: "List All Role Briefs",
		},
		{
			pathNo:  "path_591212357369856208429",
			resNo:   "res_591212357369857208429",
			url:     "goauth" + miso.OpenApiPath("/resource/list"),
			resName: "List Resources",
		},
	}

	initsql := fmt.Sprintf("INSERT INTO role(role_no, name) VALUES ('%s', '%s');", roleNo, roleName)
	for i, p := range paths {
		p.url = preprocessUrl(p.url)
		if p.pathNo == "" {
			p.pathNo = miso.GenIdP("path_")
		}
		if p.resNo == "" {
			p.resNo = miso.GenIdP("res_")
		}
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

	initsql += "\n\nINSERT INTO path(path_no, url, ptype, res_no, pgroup) VALUES"
	for i, p := range paths {
		if i > 0 {
			initsql += ","
		}
		initsql += fmt.Sprintf("\n  ('%s', '%s', '%s', '%s', 'goauth')", p.pathNo, p.url, PtProtected, p.resNo)
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
