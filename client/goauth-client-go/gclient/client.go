package gclient

import (
	"context"
	"errors"
	"strings"

	"github.com/curtisnewbie/gocommon/client"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/sirupsen/logrus"
)

type PathType string

type PathDoc struct {
	Desc   string
	Type   PathType
	Code   string
}

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

// Register GET request handler on server and report path to goauth
func Get(url string, handler server.TRouteHandler, doc PathDoc) {
	server.Get(url, handler)
	reportPathOnServerBootstrapted(url, "GET", doc)
}

// Register POST request handler on server and report path to goauth
func Post(url string, handler server.TRouteHandler, doc PathDoc) {
	server.Post(url, handler)
	reportPathOnServerBootstrapted(url, "POST", doc)
}

// Register Json POST request handler and report path to goauth
func PostJ[T any](url string, handler server.JTRouteHandler[T], doc PathDoc) {
	server.PostJ(url, handler)
	reportPathOnServerBootstrapted(url, "POST", doc)
}

// Register PUT request handler and report path to goauth
func Put(url string, handler server.TRouteHandler, doc PathDoc) {
	server.Put(url, handler)
	reportPathOnServerBootstrapted(url, "PUT", doc)
}

// Register DELETE request handler and report path to goauth
func Delete(url string, handler server.TRouteHandler, doc PathDoc) {
	server.Delete(url, handler)
	reportPathOnServerBootstrapted(url, "DELETE", doc)
}

func reportPathOnServerBootstrapted(url string, method string, doc PathDoc) {
	app := common.GetPropStr(common.PROP_APP_NAME)

	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}

	server.OnServerBootstrapped(func() {
		r := CreatePathReq{
			Method:  method,
			Group:   app,
			Url:     app + url,
			Type:    doc.Type,
			Desc:    doc.Code,
			ResCode: doc.Code,
		}
		if e := AddPath(context.Background(), r); e != nil {
			logrus.Fatal(e)
		}
	})
}
