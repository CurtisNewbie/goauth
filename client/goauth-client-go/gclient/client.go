package gclient

import (
	"context"
	"errors"

	"github.com/curtisnewbie/gocommon/client"
	"github.com/curtisnewbie/gocommon/common"
)

type PathType string

const (
	PT_PROTECTED PathType = "PROTECTED"
	PT_PUBLIC    PathType = "PUBLIC"
)

type RoleInfoReq struct {
	RoleNo string `json:"roleNo" `
}

type RoleInfoResp struct {
	RoleNo string `json:"roleNo"`
	Name   string `json:"name"`
}

type CreatePathReq struct {
	Type    PathType `json:"type"`
	Url     string   `json:"url"`
	Group   string   `json:"group"`
	Desc    string   `json:"desc"`
	ResCode string   `json:"resCode"`
	Method  string   `json:"method"`
}

type TestResAccessReq struct {
	RoleNo string `json:"roleNo"`
	Url    string `json:"url"`
}

type TestResAccessResp struct {
	Valid bool `json:"valid"`
}

type AddResourceReq struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func TestResourceAccess(ctx context.Context, req TestResAccessReq) (*TestResAccessResp, error) {
	tr := client.NewDynTClient(ctx, "/remote/path/resource/access-test", "goauth").
		EnableTracing().
		PostJson(req)

	if tr.Err != nil {
		return nil, tr.Err
	}
	defer tr.Close()

	var r common.GnResp[*TestResAccessResp]
	if e := tr.ReadJson(&r); e != nil {
		return nil, e
	}

	if r.Error {
		return nil, common.NewWebErr(r.Msg)
	}

	if r.Data == nil {
		return nil, errors.New("data is nil, unable to retrieve TestResAccessResp")
	}

	return r.Data, nil
}

func AddResource(ctx context.Context, req AddResourceReq) error {
	tr := client.NewDynTClient(ctx, "/remote/resource/add", "goauth").
		EnableTracing().
		PostJson(req)

	if tr.Err != nil {
		return tr.Err
	}
	defer tr.Close()

	var r common.Resp
	if e := tr.ReadJson(&r); e != nil {
		return e
	}
	if r.Error {
		return common.NewWebErr(r.Msg)
	}

	return nil
}

func AddPath(ctx context.Context, req CreatePathReq) error {
	tr := client.NewDynTClient(ctx, "/remote/path/add", "goauth").
		EnableTracing().
		PostJson(req)

	if tr.Err != nil {
		return tr.Err
	}
	defer tr.Close()

	var r common.Resp
	if e := tr.ReadJson(&r); e != nil {
		return e
	}
	if r.Error {
		return common.NewWebErr(r.Msg)
	}

	return nil
}

func GetRoleInfo(ctx context.Context, req RoleInfoReq) (*RoleInfoResp, error) {
	tr := client.NewDynTClient(ctx, "/remote/role/info", "goauth").
		EnableTracing().
		PostJson(req)

	if tr.Err != nil {
		return nil, tr.Err
	}
	defer tr.Close()

	var r common.GnResp[*RoleInfoResp]
	if e := tr.ReadJson(&r); e != nil {
		return nil, e
	}
	if r.Data == nil {
		return nil, errors.New("data is nil, unable to retrieve RoleInfoResp")
	}

	return r.Data, nil
}
