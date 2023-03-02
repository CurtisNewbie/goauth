package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/mysql"
	"github.com/curtisnewbie/gocommon/redis"
	"gorm.io/gorm"
)

var (
	permitted = TestResAccessResp{Valid: true}
	forbidden = TestResAccessResp{Valid: false}
)

const (
	DEFAULT_ADMIN_ROLE_NO = "role_554107924873216177918"
)

type WRole struct {
	Id         int          `json:"id"`
	RoleNo     string       `json:"roleNo"`
	Name       string       `json:"name"`
	CreateTime common.ETime `json:"createTime"`
	CreateBy   string       `json:"createBy"`
	UpdateTime common.ETime `json:"updateTime"`
	UpdateBy   string       `json:"updateBy"`
}

type CachedUrlRes struct {
	Id     int      // id
	Pgroup string   // path group
	PathNo string   // path no
	ResNo  string   // resource no
	Url    string   // url
	Ptype  PathType // path type: PROTECTED, PUBLIC
}

type ResBrief struct {
	ResNo string `json:"resNo"`
	Name  string `json:"name"`
}

type AddRoleReq struct {
	Name string `json:"name" validation:"notEmpty,maxLen:32"` // role name
}

type TestResAccessReq struct {
	RoleNo string `json:"roleNo"`
	Url    string `json:"url"`
}

type TestResAccessResp struct {
	Valid bool `json:"valid"`
}

type ListRoleReq struct {
	Paging common.Paging `json:"pagingVo"`
}

type ListRoleResp struct {
	Payload []WRole       `json:"payload"`
	Paging  common.Paging `json:"pagingVo"`
}

type RoleBrief struct {
	RoleNo string `json:"roleNo"`
	Name   string `json:"name"`
}

type ListPathReq struct {
	Pgroup string        `json:"pgroup"`
	Url    string        `json:"url"`
	Ptype  PathType      `json:"ptype"`
	Paging common.Paging `json:"pagingVo"`
}

type WPath struct {
	Id         int          `json:"id"`
	ResName    string       `json:"resName"`
	Pgroup     string       `json:"pgroup"`
	PathNo     string       `json:"pathNo"`
	ResNo      string       `json:"resNo"`
	Url        string       `json:"url"`
	Ptype      PathType     `json:"ptype"`
	CreateTime common.ETime `json:"createTime"`
	CreateBy   string       `json:"createBy"`
	UpdateTime common.ETime `json:"updateTime"`
	UpdateBy   string       `json:"updateBy"`
}

type WRes struct {
	Id         int          `json:"id"`
	ResNo      string       `json:"resNo"`
	Name       string       `json:"name"`
	CreateTime common.ETime `json:"createTime"`
	CreateBy   string       `json:"createBy"`
	UpdateTime common.ETime `json:"updateTime"`
	UpdateBy   string       `json:"updateBy"`
}

type ListPathResp struct {
	Paging  common.Paging `json:"pagingVo"`
	Payload []WPath       `json:"payload"`
}

type BindPathResReq struct {
	PathNo string `json:"pathNo" validation:"notEmpty"`
	ResNo  string `json:"resNo" validation:"notEmpty"`
}

type UnbindPathResReq struct {
	PathNo string `json:"pathNo" validation:"notEmpty"`
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
	Id         int       `json:"id"`
	ResNo      string    `json:"resNo"`
	ResName    string    `json:"resName"`
	CreateTime time.Time `json:"createTime"`
	CreateBy   string    `json:"createBy"`
}

type RoleInfoReq struct {
	RoleNo string `json:"roleNo" validation:"notEmpty"`
}

type RoleInfoResp struct {
	RoleNo string `json:"roleNo"`
	Name   string `json:"name"`
}

type UpdatePathReq struct {
	Type   PathType `json:"type" validation:"notEmpty"`
	PathNo string   `json:"pathNo" validation:"notEmpty"`
	Group  string   `json:"group" validation:"notEmpty,maxLen:20"`
}

type CreatePathReq struct {
	Type  PathType `json:"type" validation:"notEmpty"`
	Url   string   `json:"url" validation:"notEmpty,maxLen:128"`
	Group string   `json:"group" validation:"notEmpty,maxLen:20"`
}

type BatchCreatePathReq struct {
	Type  PathType `json:"type" validation:"notEmpty"`
	Urls  []string `json:"urls"`
	Group string   `json:"group" validation:"notEmpty,maxLen:20"`
}

type DeletePathReq struct {
	PathNo string `json:"pathNo" validation:"notEmpty"`
}

type ListResReq struct {
	Paging common.Paging `json:"pagingVo"`
}

type ListResResp struct {
	Paging  common.Paging `json:"pagingVo"`
	Payload []WRes        `json:"payload"`
}

type CreateResReq struct {
	Name string `json:"name" validation:"notEmpty,maxLen:32"`
}

var (
	urlResCache  = redis.NewLazyRCache(30 * time.Minute) // cache for url's resource, url -> CachedUrlRes
	roleResCache = redis.NewLazyRCache(1 * time.Hour)    // cache for role's resource, role + res -> flag ("1")
)

func ListResourceCandidatesForRole(ec common.ExecContext, roleNo string) ([]ResBrief, error) {
	if roleNo == "" {
		return []ResBrief{}, nil
	}

	var res []ResBrief
	tx := mysql.GetMySql().
		Raw(`select r.res_no, r.name from resource r 
			where not exists (select * from role_resource where role_no = ? and res_no = r.res_no)`, roleNo).
		Scan(&res)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if res == nil {
		res = []ResBrief{}
	}
	return res, nil
}

func ListAllResBriefs(ec common.ExecContext) ([]ResBrief, error) {
	var res []ResBrief
	tx := mysql.GetMySql().Raw("select res_no, name, name from resource").Scan(&res)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if res == nil {
		res = []ResBrief{}
	}
	return res, nil
}

func ListResources(ec common.ExecContext, req ListResReq) (ListResResp, error) {
	var resources []WRes
	tx := mysql.GetMySql().
		Raw("select * from resource order by id desc limit ?, ?", req.Paging.GetOffset(), req.Paging.GetLimit()).
		Scan(&resources)
	if tx.Error != nil {
		return ListResResp{}, tx.Error
	}
	if resources == nil {
		resources = []WRes{}
	}

	var count int
	tx = mysql.GetMySql().Raw("select count(*) from resource").Scan(&count)
	if tx.Error != nil {
		return ListResResp{}, tx.Error
	}

	return ListResResp{Paging: common.RespPage(req.Paging, count), Payload: resources}, nil
}

func UpdatePath(ec common.ExecContext, req UpdatePathReq) error {
	_, e := redis.RLockRun(ec, "goauth:path:"+req.PathNo, func() (any, error) {
		tx := mysql.GetMySql().Exec(`update path set pgroup = ?, ptype = ? where path_no = ?`, req.Group, req.Type,
			req.PathNo)
		return nil, tx.Error
	})

	if e == nil {
		go func(ec common.ExecContext, pathNo string) {
			ec.Log.Infof("Refreshing path cache, pathNo: %s", pathNo)
			ep, e := findPath(pathNo)
			if e != nil {
				ec.Log.Errorf("Failed to reload path cache, pathNo: %s, %v", pathNo, e)
				return
			}

			ep.Url = preprocessUrl(ep.Url)
			cachedStr, e := prepCachedUrlResStr(ec, ep)
			if e != nil {
				ec.Log.Errorf("Failed to prepare cached url resource, pathNo: %s, %v", pathNo, e)
				return
			}

			if e := urlResCache.Put(ec, ep.Url, cachedStr); e != nil {
				ec.Log.Errorf("Failed to save cached url resource, pathNo: %s, %v", pathNo, e)
				return
			}
		}(ec, req.PathNo)
	}
	return e
}

func GetRoleInfo(ec common.ExecContext, req RoleInfoReq) (RoleInfoResp, error) {
	var resp RoleInfoResp
	tx := mysql.GetMySql().Raw("select role_no, name from role where role_no = ?", req.RoleNo).Scan(&resp)
	if tx.Error != nil {
		return resp, tx.Error
	}

	if tx.RowsAffected < 1 {
		return resp, common.NewWebErrCode(EC_ROLE_NOT_FOUND, "Role not found")
	}
	return resp, tx.Error
}

func CreateResourceIfNotExist(ec common.ExecContext, req CreateResReq) error {
	req.Name = strings.TrimSpace(req.Name)
	_, e := redis.RLockRun(ec, "goauth:resource:create"+req.Name, func() (any, error) {
		var id int
		tx := mysql.GetMySql().Raw(`select id from resource where name = ? limit 1`, req.Name).Scan(&id)
		if tx.Error != nil {
			return nil, tx.Error
		}

		if id > 0 {
			ec.Log.Infof("Resource '%s' already exist", req.Name)
			return nil, nil
		}

		res := ERes{
			ResNo:    common.GenIdP("res_"),
			Name:     req.Name,
			CreateBy: ec.User.Username,
			UpdateBy: ec.User.Username,
		}

		tx = mysql.GetMySql().
			Table("resource").
			Omit("Id", "CreateTime", "UpdateTime").
			Create(&res)
		return nil, tx.Error
	})
	return e
}

func BatchCreatePathIfNotExist(ec common.ExecContext, req BatchCreatePathReq) error {
	if req.Urls == nil {
		return nil
	}

	for _, u := range req.Urls {
		if utf8.RuneCountInString(u) > 128 {
			return common.NewWebErr("URL exceeded maximum length 128")
		}

		if e := CreatePathIfNotExist(ec, CreatePathReq{
			Type:  req.Type,
			Group: req.Group,
			Url:   u,
		}); e != nil {
			return e
		}
	}
	return nil
}

func CreatePathIfNotExist(ec common.ExecContext, req CreatePathReq) error {
	_, e := redis.RLockRun(ec, "goauth:path:url"+req.Url, func() (any, error) {

		var id int
		tx := mysql.GetMySql().Raw(`select id from path where url = ? limit 1`, req.Url).Scan(&id)
		if tx.Error != nil {
			return nil, tx.Error
		}

		if id > 0 {
			// ec.Log.Infof("Path '%s' already exist", req.Url)
			return nil, nil
		}

		ep := EPath{
			Url:      req.Url,
			Ptype:    req.Type,
			Pgroup:   req.Group,
			PathNo:   common.GenIdP("path_"),
			CreateBy: ec.User.Username,
			UpdateBy: ec.User.Username,
		}
		tx = mysql.GetMySql().
			Table("path").
			Omit("Id", "CreateTime", "UpdateTime").
			Create(&ep)
		return nil, tx.Error
	})
	return e
}

func DeletePath(ec common.ExecContext, req DeletePathReq) error {
	_, e := redis.RLockRun(ec, "goauth:path:"+req.PathNo, func() (any, error) {
		tx := mysql.GetMySql().Exec(`delete from path where path_no = ?`, req.PathNo)
		return nil, tx.Error
	})
	return e
}

func UnbindPathRes(ec common.ExecContext, req UnbindPathResReq) error {
	_, e := redis.RLockRun(ec, "goauth:path:"+req.PathNo, func() (any, error) {
		tx := mysql.GetMySql().Exec(`update path set res_no = '' where path_no = ?`, req.PathNo)
		return nil, tx.Error
	})
	return e
}

func BindPathRes(ec common.ExecContext, req BindPathResReq) error {
	_, e := redis.RLockRun(ec, "goauth:path:"+req.PathNo, func() (any, error) {
		tx := mysql.GetMySql().Exec(`update path set res_no = ? where path_no = ?`, req.ResNo, req.PathNo)
		return nil, tx.Error
	})
	return e
}

func ListPaths(ec common.ExecContext, req ListPathReq) (ListPathResp, error) {
	var paths []WPath
	tx := mysql.GetMySql().
		Table("path p").
		Select("p.*, r.name res_name").
		Joins("left join resource r on p.res_no = r.res_no").
		Order("id desc")

	if req.Pgroup != "" {
		tx = tx.Where("p.pgroup = ?", req.Pgroup)
	}
	if req.Url != "" {
		tx = tx.Where("p.url like ?", req.Url+"%")
	}
	if req.Ptype != "" {
		tx = tx.Where("p.ptype = ?", req.Ptype)
	}

	tx = tx.Offset(req.Paging.GetOffset()).
		Limit(req.Paging.GetLimit()).
		Scan(&paths)
	if tx.Error != nil {
		return ListPathResp{}, tx.Error
	}

	var count int
	tx = mysql.GetMySql().
		Table("path p").
		Select("count(*)").
		Joins("left join resource r on p.res_no = r.res_no")

	if req.Pgroup != "" {
		tx = tx.Where("p.pgroup = ?", req.Pgroup)
	}
	if req.Url != "" {
		tx = tx.Where("p.url like ?", req.Url+"%")
	}
	if req.Ptype != "" {
		tx = tx.Where("p.ptype = ?", req.Ptype)
	}

	tx = tx.Scan(&count)
	if tx.Error != nil {
		return ListPathResp{}, tx.Error
	}

	return ListPathResp{Payload: paths, Paging: common.Paging{Limit: req.Paging.Limit, Page: req.Paging.Page, Total: count}}, nil
}

func AddRole(ec common.ExecContext, req AddRoleReq) error {
	_, e := redis.RLockRun(ec, "goauth:role:add"+req.Name, func() (any, error) {
		r := ERole{
			RoleNo:   common.GenIdP("role_"),
			Name:     req.Name,
			CreateBy: ec.User.Username,
			UpdateBy: ec.User.Username,
		}
		return nil, mysql.GetMySql().
			Table("role").
			Omit("Id", "CreateTime", "UpdateTime").
			Create(&r).Error
	})
	return e
}

func RemoveResFromRole(ec common.ExecContext, req RemoveRoleResReq) error {
	_, e := redis.RLockRun(ec, "goauth:role:"+req.RoleNo, func() (any, error) {
		tx := mysql.GetMySql().Exec(`delete from role_resource where role_no = ? and res_no = ?`, req.RoleNo, req.ResNo)
		return nil, tx.Error
	})
	return e
}

func AddResToRoleIfNotExist(ec common.ExecContext, req AddRoleResReq) error {
	_, e := redis.RLockRun(ec, "goauth:role:"+req.RoleNo, func() (any, error) {
		var id int
		tx := mysql.GetMySql().Raw(`select id from role_resource where role_no = ? and res_no = ?`, req.RoleNo, req.ResNo).Scan(&id)
		if tx.Error != nil {
			return nil, tx.Error
		}
		if id > 0 {
			return nil, nil
		}

		rr := ERoleRes{
			RoleNo:   req.RoleNo,
			ResNo:    req.ResNo,
			CreateBy: ec.User.Username,
			UpdateBy: ec.User.Username,
		}

		tx = mysql.GetMySql().
			Table("role_resource").
			Omit("Id", "CreateTime", "UpdateTime").
			Create(&rr)
		return nil, tx.Error
	})
	return e
}

func ListRoleRes(ec common.ExecContext, req ListRoleResReq) (ListRoleResResp, error) {
	var res []ListedRoleRes
	tx := mysql.GetMySql().
		Raw(`select rr.id, rr.res_no, rr.create_time, rr.create_by, r.name 'res_name' from role_resource rr 
			left join resource r on rr.res_no = r.res_no
			where rr.role_no = ? order by rr.id desc limit ?, ?`, req.RoleNo, req.Paging.GetOffset(), req.Paging.GetLimit()).
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
			where rr.role_no = ?`, req.RoleNo).
		Scan(&count)

	if tx.Error != nil {
		return ListRoleResResp{}, tx.Error
	}

	return ListRoleResResp{Payload: res, Paging: common.Paging{Limit: req.Paging.Limit, Page: req.Paging.Page, Total: count}}, nil
}

func ListAllRoleBriefs(ec common.ExecContext) ([]RoleBrief, error) {
	var roles []RoleBrief
	tx := mysql.GetMySql().Raw("select role_no, name from role").Scan(&roles)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if roles == nil {
		roles = []RoleBrief{}
	}
	return roles, nil
}

func ListRoles(ec common.ExecContext, req ListRoleReq) (ListRoleResp, error) {
	var roles []WRole
	tx := mysql.GetMySql().
		Raw("select * from role order by id desc limit ?, ?", req.Paging.GetOffset(), req.Paging.GetLimit()).
		Scan(&roles)
	if tx.Error != nil {
		return ListRoleResp{}, tx.Error
	}
	if roles == nil {
		roles = []WRole{}
	}

	var count int
	tx = mysql.GetMySql().Raw("select count(*) from role").Scan(&count)
	if tx.Error != nil {
		return ListRoleResp{}, tx.Error
	}

	return ListRoleResp{Payload: roles, Paging: common.Paging{Limit: req.Paging.Limit, Page: req.Paging.Page, Total: count}}, nil
}

// Test access to resource
func TestResourceAccess(ec common.ExecContext, req TestResAccessReq) (TestResAccessResp, error) {
	url := req.Url
	roleNo := req.RoleNo

	if roleNo == DEFAULT_ADMIN_ROLE_NO {
		return permitted, nil
	}

	// some sanitization & standardization for the url
	url = preprocessUrl(url)

	// find resource required for the url
	cur, e := lookupUrlRes(ec, url)
	if e != nil {
		return forbidden, e
	}

	// public path type, doesn't require access to resource
	if cur.Ptype == PT_PUBLIC {
		return permitted, nil
	}
	ec.Log.Infof("'%s' is protected, validating resource access", url)

	// doesn't even have role
	roleNo = strings.TrimSpace(roleNo)
	if roleNo == "" {
		ec.Log.Infof("Rejected '%s', user doesn't have roleNo", url)
		return forbidden, nil
	}

	// the requiredRes resources no
	requiredRes := cur.ResNo
	if requiredRes == "" {
		ec.Log.Infof("Rejected '%s', path doesn't have any resource bound yet", url)
		return forbidden, nil
	}

	ok, e := checkRoleRes(ec, roleNo, requiredRes)
	if e != nil {
		return forbidden, e
	}

	// the role doesn't have access to the required resource
	if !ok {
		ec.Log.Infof("Rejected '%s', roleNo: '%s', role doesn't have access to required resource '%s'", url, roleNo, requiredRes)
		return forbidden, nil
	}

	return permitted, nil
}

func checkRoleRes(ec common.ExecContext, roleNo string, resNo string) (bool, error) {
	r, e := roleResCache.Get(ec, fmt.Sprintf("role:%s:res:%s", roleNo, resNo))
	if e != nil {
		return false, e
	}

	return r != "", nil
}

func LoadRoleResCache(ec common.ExecContext) error {

	_, e := redis.RLockRun(ec, "goauth:role:res:cache", func() (any, error) {

		// ec.Log.Info("Loading role resource cache")

		lr, e := listRoleNos(ec)
		if e != nil {
			return nil, e
		}

		for _, roleNo := range lr {
			roleResList, e := listRoleRes(ec, roleNo)
			if e != nil {
				return nil, e
			}

			for _, rr := range roleResList {
				roleResCache.Put(ec, fmt.Sprintf("role:%s:res:%s", rr.RoleNo, rr.ResNo), "1")
				// ec.Log.Infof("Loaded RoleRes: '%s' -> '%s'", rr.RoleNo, rr.ResNo)
			}
		}
		return nil, nil
	})
	return e
}

func listRoleNos(ec common.ExecContext) ([]string, error) {
	var ern []string
	t := mysql.GetMySql().Raw("select role_no from role").Scan(&ern)
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
	if js == "" {
		return CachedUrlRes{}, common.NewWebErr(fmt.Sprintf("Unable to find path '%s'", url))
	}

	var cur CachedUrlRes
	if e = json.Unmarshal([]byte(js), &cur); e != nil {
		return CachedUrlRes{}, e
	}

	return cur, nil
}

func LoadPathResCache(ec common.ExecContext) error {

	_, e := redis.RLockRun(ec, "goauth:path:res:cache", func() (any, error) {

		// ec.Log.Info("Loading path resource cache")

		var paths []EPath
		tx := mysql.GetMySql().Raw("select * from path").Scan(&paths)
		if tx.Error != nil {
			return nil, tx.Error
		}
		if paths == nil {
			return nil, nil
		}

		for _, ep := range paths {
			ep.Url = preprocessUrl(ep.Url)
			cachedStr, e := prepCachedUrlResStr(ec, ep)
			if e != nil {
				return nil, e
			}
			if e := urlResCache.Put(ec, ep.Url, cachedStr); e != nil {
				return nil, e
			}
			// ec.Log.Infof("Loaded PathRes: '%s', '%s', '%s'", ep.Url, ep.Ptype, ep.ResNo)
		}
		return nil, nil
	})

	return e
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
func preprocessUrl(url string) string {
	ru := []rune(strings.TrimSpace(url))
	l := len(ru)
	if l < 1 {
		return "/"
	}

	j := strings.LastIndex(url, "?")
	if j > -1 {
		ru = ru[0:j]
		l = len(ru)
	}

	// never ends with '/'
	if ru[l-1] == '/' {
		ru = ru[0 : l-1]
		url = string(ru)
	}

	// always start with '/'
	if ru[0] != '/' {
		url = "/" + url
	}
	return url
}

func findPath(pathNo string) (EPath, error) {
	var ep EPath
	tx := mysql.GetMySql().Raw("select * from path where path_no = ?", pathNo).Scan(&ep)
	if tx.Error != nil {
		return ep, tx.Error
	}

	if tx.RowsAffected < 1 {
		return ep, common.NewWebErr("Path not found")
	}

	return ep, nil
}
