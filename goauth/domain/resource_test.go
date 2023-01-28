package domain

import (
	"testing"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/mysql"
	"github.com/curtisnewbie/gocommon/redis"
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
