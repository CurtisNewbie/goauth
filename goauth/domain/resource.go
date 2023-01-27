package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/mysql"
	"github.com/curtisnewbie/gocommon/redis"
	"gorm.io/gorm"
)

type CachedUrlRes struct {
	Id     int      // id
	Pgroup string   // path group
	PathNo string   // path no
	ResNo  string   // resource no
	Url    string   // url
	Ptype  PathType // path type: PROTECTED, PUBLIC
}

type AddRoleReq struct {
	Name string // role name
}

type TestResAccessReq struct {
	RoleNo string `json:"roleNo"`
	Url    string `json:"url"`
}

type ListRoleReq struct {
	Paging common.Paging `json:"pagingVo"`
}

type ListRoleResp struct {
	Payload []ERole       `json:"payload"`
	Paging  common.Paging `json:"pagingVo"`
}

type ListPathReq struct {
	Paging common.Paging `json:"pagingVo"`
}

type ListPathResp struct {
	Paging  common.Paging `json:"pagingVo"`
	Payload []EPath       `json:"payload"`
}

type BindPathResReq struct {
	PathNo string `json:"pathNo"`
	ResNo  string `json:"resNo"`
}

type UnbindPathResReq struct {
	PathNo string `json:"pathNo"`
}

type ListRoleResReq struct {
	Paging common.Paging `json:"pagingVo"`
	RoleNo string        `json:"roleNo" validation:"notEmpty"`
}

type RemoveRoleResReq struct {
	RoleNo string `json:"roleNo" validation:"notEmpty"`
	ResNo  string `json:"resNo" validation:"notEmpty"`
}

type AddRoleResReq struct {
	RoleNo string `json:"roleNo" validation:"notEmpty"`
	ResNo  string `json:"resNo" validation:"notEmpty"`
}

type ListRoleResResp struct {
	Paging  common.Paging   `json:"pagingVo"`
	Payload []ListedRoleRes `json:"payload"`
}

type ListedRoleRes struct {
	Id         int    // id
	RoleNo     string // role no
	ResNo      string // resource no
	ResName    string // resource name
	Url        string // url
	CreateTime common.ETime
	CreateBy   string
	UpdateTime common.ETime
	UpdateBy   string
}

var (
	urlResCache  = redis.NewLazyRCache(30 * time.Minute) // cache for url's resource, url -> CachedUrlRes
	roleResCache = redis.NewLazyRCache(1 * time.Hour)    // cache for role's resource, role + res -> flag ("1")
)

func UnbindPathRes(ec common.ExecContext, req UnbindPathResReq) error {
	_, e := redis.RLockRun(ec, "goauth:path:"+req.PathNo, func() (any, error) {
		tx := mysql.GetMySql().Raw(`update path set res_no = '' where path_no = ?`, req.PathNo)
		return nil, tx.Error
	})
	return e
}

func BindPathRes(ec common.ExecContext, req BindPathResReq) error {
	_, e := redis.RLockRun(ec, "goauth:path:"+req.PathNo, func() (any, error) {
		tx := mysql.GetMySql().Raw(`update path set res_no = ? where path_no = ?`, req.ResNo, req.PathNo)
		return nil, tx.Error
	})
	return e
}

func ListPaths(ec common.ExecContext, req ListPathReq) (ListPathResp, error) {
	var paths []EPath
	tx := mysql.GetMySql().
		Raw("select * from path where limit ?, ?", common.CalcOffset(&req.Paging), req.Paging.Limit).
		Scan(&paths)
	if tx.Error != nil {
		return ListPathResp{}, tx.Error
	}

	var count int
	tx = mysql.GetMySql().
		Raw("select count(*) from path where limit ?, ?", common.CalcOffset(&req.Paging), req.Paging.Limit).
		Scan(&count)
	if tx.Error != nil {
		return ListPathResp{}, tx.Error
	}

	return ListPathResp{Payload: paths, Paging: common.Paging{Limit: req.Paging.Limit, Page: req.Paging.Page, Total: count}}, nil
}

func AddRole(ec common.ExecContext, req AddRoleReq) error {
	r := ERole{
		Name: req.Name,
	}
	return mysql.GetMySql().Table("role").Save(r).Error
}

func RemoveResFromRole(ec common.ExecContext, req RemoveRoleResReq) error {
	_, e := redis.RLockRun(ec, "goauth:role:"+req.RoleNo, func() (any, error) {
		tx := mysql.GetMySql().Raw(`delete from role_resource where role_no = ? and res_no = ?`, req.RoleNo, req.ResNo)
		return nil, tx.Error
	})
	return e
}

func AddResToRole(ec common.ExecContext, req AddRoleResReq) error {
	_, e := redis.RLockRun(ec, "goauth:role:"+req.RoleNo, func() (any, error) {
		tx := mysql.GetMySql().Raw(`select id from role_resource where role_no = ? and res_no = ?`, req.RoleNo, req.ResNo)
		if tx.Error != nil {
			return nil, tx.Error
		}

		if tx.RowsAffected > 0 {
			return nil, common.NewWebErr("Resource already exists in this role")
		}

		rr := ERoleRes{
			RoleNo:   req.RoleNo,
			ResNo:    req.ResNo,
			CreateBy: ec.User.Username,
		}

		tx = mysql.GetMySql().Table("role_resource").Save(rr)
		return nil, tx.Error
	})
	return e
}

func ListRoleRes(ec common.ExecContext, req ListRoleResReq) (ListRoleResResp, error) {
	var res []ListedRoleRes
	offset := common.CalcOffset(&req.Paging)
	tx := mysql.GetMySql().
		Raw(`select rr.*, r.name 'res_name' from role_resource rr 
			left join resource r on rr.res_no = r.res_no
			where rr.role_no = ? limit ?, ?`, req.RoleNo, offset, req.Paging.Limit).
		Scan(&res)

	if tx.Error != nil {
		return ListRoleResResp{}, tx.Error
	}

	if res == nil {
		res = []ListedRoleRes{}
	}

	var count int
	tx = mysql.GetMySql().
		Raw(`select count(*) from role_resource rr 
			left join resource r on rr.res_no = r.res_no
			where rr.role_no = ? limit ?, ?`, req.RoleNo, offset, req.Paging.Limit).
		Scan(&res)

	if tx.Error != nil {
		return ListRoleResResp{}, tx.Error
	}

	return ListRoleResResp{Payload: res, Paging: common.Paging{Limit: req.Paging.Limit, Page: req.Paging.Page, Total: count}}, nil
}

func ListRoles(ec common.ExecContext, req ListRoleReq) (ListRoleResp, error) {
	var roles []ERole
	offset := common.CalcOffset(&req.Paging)
	tx := mysql.GetMySql().Raw("select * from role limit ?, ?", offset, req.Paging.Limit).Scan(&roles)
	if tx.Error != nil {
		return ListRoleResp{}, tx.Error
	}
	if roles == nil {
		roles = []ERole{}
	}

	var count int
	tx = mysql.GetMySql().Raw("select count(*) from role limit ?, ?", offset, req.Paging.Limit).Scan(&count)
	if tx.Error != nil {
		return ListRoleResp{}, tx.Error
	}

	return ListRoleResp{Payload: roles, Paging: common.Paging{Limit: req.Paging.Limit, Page: req.Paging.Page, Total: count}}, nil
}

// Test access to resource
func TestResourceAccess(ec common.ExecContext, req TestResAccessReq) error {
	url := req.Url
	roleNo := req.RoleNo

	// some sanitization & standardization for the url
	url = preprocessUrl(url)

	// find resource required for the url
	cur, e := lookupUrlRes(ec, url)
	if e != nil {
		return e
	}

	// public path type, doesn't require access to resource
	if cur.Ptype == PT_PUBLIC {
		return nil
	}

	// doesn't even have role
	roleNo = strings.TrimSpace(roleNo)
	if roleNo == "" {
		ec.Log.Infof("Rejected '%s', roleNo: '%s', role is empty", url, roleNo)
		return common.NewWebErr("not permitted access")
	}

	// the requiredRes resources no
	requiredRes := cur.ResNo
	ok, e := CheckRoleRes(ec, roleNo, requiredRes)
	if e != nil {
		return e
	}

	// the role doesn't have access to the required resource
	if !ok {
		ec.Log.Infof("Rejected '%s', roleNo: '%s', require access to resource '%s'", url, roleNo, requiredRes)
		return common.NewWebErr("not permitted access")
	}

	return nil
}

func CheckRoleRes(ec common.ExecContext, roleNo string, resNo string) (bool, error) {
	r, e := roleResCache.Get(ec, fmt.Sprintf("role:%s:res:%s", roleNo, resNo))
	if e != nil {
		return false, e
	}

	return r != "", nil
}

func LoadRoleResCache(ec common.ExecContext) error {
	ec.Log.Info("Loading role resource cache")
	lr, e := listRoles(ec)
	if e != nil {
		return e
	}

	for _, roleNo := range lr {
		roleResList, e := listRoleRes(ec, roleNo)
		if e != nil {
			return e
		}

		for _, rr := range roleResList {
			roleResCache.Put(ec, fmt.Sprintf("role:%s:res:%s", rr.RoleNo, rr.ResNo), "1")
		}
	}
	return nil
}

func listRoles(ec common.ExecContext) ([]string, error) {
	var ern []string
	t := mysql.GetMySql().Raw("select * from role").Scan(&ern)
	if t.Error != nil {
		return nil, t.Error
	}

	if ern == nil {
		ern = []string{}
	}
	return ern, nil
}

func listRoleRes(ec common.ExecContext, roleNo string) ([]ERoleRes, error) {
	var rr []ERoleRes
	t := mysql.GetMySql().Raw("select * from role_resource where role_no = ?", roleNo).Scan(&rr)
	if t.Error != nil {
		if errors.Is(t.Error, gorm.ErrRecordNotFound) {
			return []ERoleRes{}, nil
		}
		return nil, t.Error
	}

	return rr, nil
}

func lookupUrlRes(ec common.ExecContext, url string) (CachedUrlRes, error) {
	js, e := urlResCache.Get(ec, url)
	if e != nil {
		return CachedUrlRes{}, e
	}

	var cur CachedUrlRes
	if e = json.Unmarshal([]byte(js), &cur); e != nil {
		return CachedUrlRes{}, e
	}

	return cur, nil
}

func LoadPathResCache(ec common.ExecContext) error {
	ec.Log.Info("Loading path resource cache")
	var paths []EPath
	tx := mysql.GetMySql().Raw("select * from path").Scan(&paths)
	if tx.Error != nil {
		return tx.Error
	}
	if paths == nil {
		return nil
	}

	for _, ep := range paths {
		cachedStr, e := prepCachedUrlResStr(ec, ep)
		if e != nil {
			return e
		}
		if e := urlResCache.Put(ec, ep.Url, cachedStr); e != nil {
			return e
		}
	}
	return nil
}

func prepCachedUrlResStr(ec common.ExecContext, epath EPath) (string, error) {
	url := epath.Url
	cur := CachedUrlRes{
		Id:     epath.Id,
		Pgroup: epath.Pgroup,
		PathNo: epath.PathNo,
		ResNo:  epath.ResNo,
		Url:    epath.Url,
		Ptype:  epath.Ptype,
	}

	j, e := json.Marshal(cur)
	if e != nil {
		ec.Log.Errorf("Failed to marshal EPath for '%s', %v", url, e)
		return "", e
	}
	return string(j), nil
}

// preprocess url, the processed url will always starts with '/' and never ends with '/'
func preprocessUrl(url string) (processed_url string) {
	ru := []rune(strings.TrimSpace(url))
	l := len(ru)
	if l < 1 {
		processed_url = "/"
		return
	}

	j := strings.LastIndex(url, "?")
	if j > -1 {
		ru = ru[0:j]
		l = len(ru)
	}

	// never ends with '/'
	if ru[l-1] == '/' {
		ru = ru[0 : l-1]
		processed_url = string(ru)
	}

	// always start with '/'
	if ru[0] != '/' {
		processed_url = "/" + processed_url
	}

	return
}
