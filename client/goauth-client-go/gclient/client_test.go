package gclient

import (
	"context"
	"testing"
)

func TestTestResourceAccess(t *testing.T) {
	req := TestResAccessReq{
		Url:    "/open/api/resource/add",
		RoleNo: "role_554107924873216177918",
	}

	r, e := TestResourceAccess(context.Background(), req)
	if e != nil {
		t.Fatal(e)
	}
	if r == nil {
		t.Fatal("r is nil")
	}
	if !r.Valid {
		t.Fatal("Should be valid")
	}
	t.Logf("r: %+v", r)
}

func TestAddPath(t *testing.T) {
	req := CreatePathReq{
		Url:   "/test/url/gclient",
		Type:  PT_PROTECTED,
		Group: "goauth-client-java",
	}

	e := AddPath(context.Background(), req)
	if e != nil {
		t.Fatal(e)
	}
}

func TestGetRoleInfo(t *testing.T) {
	req := RoleInfoReq{
		RoleNo: "role_554107924873216177918",
	}

	r, e := GetRoleInfo(context.Background(), req)
	if e != nil {
		t.Fatal(e)
	}
	if r == nil {
		t.Fatal("r is nil")
	}
	if r.RoleNo == "" {
		t.Fatal("roleNo should not be empty")
	}
	if r.Name == "" {
		t.Fatal("name should not be empty")
	}
	t.Logf("r: %+v", r)
}
