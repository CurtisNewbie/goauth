package domain

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/miso/core"
	"github.com/curtisnewbie/miso/mysql"
	"github.com/curtisnewbie/miso/redis"
	"gorm.io/gorm"
)

var (
	permitted = TestResAccessResp{Valid: true}
	forbidden = TestResAccessResp{Valid: false}

	roleInfoCache = redis.NewLazyObjectRCache[RoleInfoResp](10 * time.Minute)
)

type PathType string

const (
	// default roleno for admin
	DEFAULT_ADMIN_ROLE_NO = "role_554107924873216177918"

	PT_PROTECTED PathType = "PROTECTED"
	PT_PUBLIC    PathType = "PUBLIC"
)

type PathRes struct {
	Id         int    // id
	PathNo     string // path no
	ResCode    string // resource code
	CreateTime core.ETime
	CreateBy   string
	UpdateTime core.ETime
	UpdateBy   string
}

type ExtendedPathRes struct {
	Id         int      // id
	Pgroup     string   // path group
	PathNo     string   // path no
	ResCode    string   // resource code
	Desc       string   // description
	Url        string   // url
	Method     string   // http method
	Ptype      PathType // path type: PROTECTED, PUBLIC
	CreateTime core.ETime
	CreateBy   string
	UpdateTime core.ETime
	UpdateBy   string
}

type EPath struct {
	Id         int      // id
	Pgroup     string   // path group
	PathNo     string   // path no
	Desc       string   // description
	Url        string   // url
	Method     string   // method
	Ptype      PathType // path type: PROTECTED, PUBLIC
	CreateTime core.ETime
	CreateBy   string
	UpdateTime core.ETime
	UpdateBy   string
}

type ERes struct {
	Id         int    // id
	Code       string // resource code
	Name       string // resource name
	CreateTime core.ETime
	CreateBy   string
	UpdateTime core.ETime
	UpdateBy   string
}

type ERoleRes struct {
	Id         int    // id
	RoleNo     string // role no
	ResCode    string // resource code
	CreateTime core.ETime
	CreateBy   string
	UpdateTime core.ETime
	UpdateBy   string
}

type ERole struct {
	Id         int
	RoleNo     string
	Name       string
	CreateTime core.ETime
	CreateBy   string
	UpdateTime core.ETime
	UpdateBy   string
}

type WRole struct {
	Id         int        `json:"id"`
	RoleNo     string     `json:"roleNo"`
	Name       string     `json:"name"`
	CreateTime core.ETime `json:"createTime"`
	CreateBy   string     `json:"createBy"`
	UpdateTime core.ETime `json:"updateTime"`
	UpdateBy   string     `json:"updateBy"`
}

type CachedUrlRes struct {
	Id      int      // id
	Pgroup  string   // path group
	PathNo  string   // path no
	ResCode string   // resource code
	Url     string   // url
	Method  string   // http method
	Ptype   PathType // path type: PROTECTED, PUBLIC
}

type ResBrief struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type AddRoleReq struct {
	Name string `json:"name" validation:"notEmpty,maxLen:32"` // role name
}

type TestResAccessReq struct {
	RoleNo string `json:"roleNo"`
	Url    string `json:"url"`
	Method string `json:"method"`
}

type TestResAccessResp struct {
	Valid bool `json:"valid"`
}

type ListRoleReq struct {
	Paging core.Paging `json:"pagingVo"`
}

type ListRoleResp struct {
	Payload []WRole     `json:"payload"`
	Paging  core.Paging `json:"pagingVo"`
}

type RoleBrief struct {
	RoleNo string `json:"roleNo"`
	Name   string `json:"name"`
}

type ListPathReq struct {
	ResCode string      `json:"resCode"`
	Pgroup  string      `json:"pgroup"`
	Url     string      `json:"url"`
	Ptype   PathType    `json:"ptype"`
	Paging  core.Paging `json:"pagingVo"`
}

type WPath struct {
	Id         int        `json:"id"`
	Pgroup     string     `json:"pgroup"`
	PathNo     string     `json:"pathNo"`
	Method     string     `json:"method"`
	Desc       string     `json:"desc"`
	Url        string     `json:"url"`
	Ptype      PathType   `json:"ptype"`
	CreateTime core.ETime `json:"createTime"`
	CreateBy   string     `json:"createBy"`
	UpdateTime core.ETime `json:"updateTime"`
	UpdateBy   string     `json:"updateBy"`
}

type WRes struct {
	Id         int        `json:"id"`
	Code       string     `json:"code"`
	Name       string     `json:"name"`
	CreateTime core.ETime `json:"createTime"`
	CreateBy   string     `json:"createBy"`
	UpdateTime core.ETime `json:"updateTime"`
	UpdateBy   string     `json:"updateBy"`
}

type ListPathResp struct {
	Paging  core.Paging `json:"pagingVo"`
	Payload []WPath     `json:"payload"`
}

type BindPathResReq struct {
	PathNo  string `json:"pathNo" validation:"notEmpty"`
	ResCode string `json:"resCode" validation:"notEmpty"`
}

type UnbindPathResReq struct {
	PathNo  string `json:"pathNo" validation:"notEmpty"`
	ResCode string `json:"resCode" validation:"notEmpty"`
}

type ListRoleResReq struct {
	Paging core.Paging `json:"pagingVo"`
	RoleNo string      `json:"roleNo" validation:"notEmpty"`
}

type RemoveRoleResReq struct {
	RoleNo  string `json:"roleNo" validation:"notEmpty"`
	ResCode string `json:"resCode" validation:"notEmpty"`
}

type AddRoleResReq struct {
	RoleNo  string `json:"roleNo" validation:"notEmpty"`
	ResCode string `json:"resCode" validation:"notEmpty"`
}

type ListRoleResResp struct {
	Paging  core.Paging     `json:"pagingVo"`
	Payload []ListedRoleRes `json:"payload"`
}

type ListedRoleRes struct {
	Id         int       `json:"id"`
	ResCode    string    `json:"resCode"`
	ResName    string    `json:"resName"`
	CreateTime time.Time `json:"createTime"`
	CreateBy   string    `json:"createBy"`
}

type RoleInfoReq struct {
	RoleNo string `json:"roleNo" validation:"notEmpty"`
}

type GenResScriptReq struct {
	ResCodes []string `json:"resCodes" validation:"notEmpty"`
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
	Type    PathType `json:"type" validation:"notEmpty"`
	Url     string   `json:"url" validation:"notEmpty,maxLen:128"`
	Group   string   `json:"group" validation:"notEmpty,maxLen:20"`
	Method  string   `json:"method" validation:"notEmpty,maxLen:10"`
	Desc    string   `json:"desc" validation:"maxLen:255"`
	ResCode string   `json:"resCode"`
}

type DeletePathReq struct {
	PathNo string `json:"pathNo" validation:"notEmpty"`
}

type ListResReq struct {
	Paging core.Paging `json:"pagingVo"`
}

type ListResResp struct {
	Paging  core.Paging `json:"pagingVo"`
	Payload []WRes      `json:"payload"`
}

type CreateResReq struct {
	Name string `json:"name" validation:"notEmpty,maxLen:32"`
	Code string `json:"code" validation:"notEmpty,maxLen:32"`
}

type DeleteResourceReq struct {
	ResCode string `json:"resCode" validation:"notEmpty"`
}

var (
	urlResCache  = redis.NewLazyRCache(30 * time.Minute) // cache for url's resource, url -> CachedUrlRes
	roleResCache = redis.NewLazyRCache(1 * time.Hour)    // cache for role's resource, role + res -> flag ("1")
)

func DeleteResource(ec core.Rail, req DeleteResourceReq) error {

	_, e := lockResourceGlobal(ec, func() (any, error) {
		return nil, mysql.GetConn().Transaction(func(tx *gorm.DB) error {
			if t := tx.Exec(`delete from resource where code = ?`, req.ResCode); t != nil {
				return t.Error
			}
			if t := tx.Exec(`delete from role_resource where res_code = ?`, req.ResCode); t != nil {
				return t.Error
			}
			return tx.Exec(`delete from path_resource where res_code = ?`, req.ResCode).Error
		})
	})

	if e != nil {
		// asynchronously reload the cache of paths and resources
		go func() {
			if e := LoadPathResCache(ec); e != nil {
				ec.Errorf("Failed to load path resource cache, %v", e)
			}
		}()
		// asynchronously reload the cache of role and resources
		go func() {
			if e := LoadRoleResCache(ec); e != nil {
				ec.Errorf("Failed to load role resource cache, %v", e)
			}
		}()
	}

	return e
}

func ListResourceCandidatesForRole(ec core.Rail, roleNo string) ([]ResBrief, error) {
	if roleNo == "" {
		return []ResBrief{}, nil
	}

	var res []ResBrief
	tx := mysql.GetConn().
		Select("r.name, r.code").
		Table("resource r").
		Where("NOT EXISTS (SELECT * FROM role_resource WHERE role_no = ? and res_code = r.code)", roleNo).
		Scan(&res)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if res == nil {
		res = []ResBrief{}
	}
	return res, nil
}

func ListAllResBriefsOfRole(ec core.Rail, roleNo string) ([]ResBrief, error) {
	var res []ResBrief

	if roleNo == DEFAULT_ADMIN_ROLE_NO {
		return ListAllResBriefs(ec)
	}

	tx := mysql.GetConn().
		Select(`r.name, r.code`).
		Table(`role_resource rr`).
		Joins(`LEFT JOIN resource r ON r.code = rr.res_code`).
		Where(`rr.role_no = ?`, roleNo).
		Scan(&res)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if res == nil {
		res = []ResBrief{}
	}
	return res, nil
}

func ListAllResBriefs(ec core.Rail) ([]ResBrief, error) {
	var res []ResBrief
	tx := mysql.GetConn().Raw("select name, code from resource").Scan(&res)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if res == nil {
		res = []ResBrief{}
	}
	return res, nil
}

func ListResources(ec core.Rail, req ListResReq) (ListResResp, error) {
	var resources []WRes
	tx := mysql.GetConn().
		Raw("select * from resource order by id desc limit ?, ?", req.Paging.GetOffset(), req.Paging.GetLimit()).
		Scan(&resources)
	if tx.Error != nil {
		return ListResResp{}, tx.Error
	}
	if resources == nil {
		resources = []WRes{}
	}

	var count int
	tx = mysql.GetConn().Raw("select count(*) from resource").Scan(&count)
	if tx.Error != nil {
		return ListResResp{}, tx.Error
	}

	return ListResResp{Paging: core.RespPage(req.Paging, count), Payload: resources}, nil
}

func UpdatePath(ec core.Rail, req UpdatePathReq) error {
	_, e := lockPath(ec, req.PathNo, func() (any, error) {
		tx := mysql.GetConn().Exec(`update path set pgroup = ?, ptype = ? where path_no = ?`,
			req.Group, req.Type, req.PathNo)
		return nil, tx.Error
	})

	if e == nil {
		loadOnePathResCacheAsync(ec, req.PathNo)
	}
	return e
}

func loadOnePathResCacheAsync(ec core.Rail, pathNo string) {
	go func(ec core.Rail, pathNo string) {
		// ec.Infof("Refreshing path cache, pathNo: %s", pathNo)
		ep, e := findPathRes(pathNo)
		if e != nil {
			ec.Errorf("Failed to reload path cache, pathNo: %s, %v", pathNo, e)
			return
		}

		ep.Url = preprocessUrl(ep.Url)
		cachedStr, e := prepCachedUrlResStr(ec, ep)
		if e != nil {
			ec.Errorf("Failed to prepare cached url resource, pathNo: %s, %v", pathNo, e)
			return
		}

		if e := urlResCache.Put(ec, ep.Method+":"+ep.Url, cachedStr); e != nil {
			ec.Errorf("Failed to save cached url resource, pathNo: %s, %v", pathNo, e)
			return
		}
	}(ec, pathNo)
}

func GetRoleInfo(ec core.Rail, req RoleInfoReq) (RoleInfoResp, error) {
	resp, _, err := roleInfoCache.GetElse(ec, req.RoleNo, func() (RoleInfoResp, bool, error) {
		var resp RoleInfoResp
		tx := mysql.GetConn().Raw("select role_no, name from role where role_no = ?", req.RoleNo).Scan(&resp)
		if tx.Error != nil {
			return resp, false, tx.Error
		}

		if tx.RowsAffected < 1 {
			return resp, false, core.NewWebErrCode(EC_ROLE_NOT_FOUND, "Role not found")
		}
		return resp, true, nil
	})
	return resp, err
}

func CreateResourceIfNotExist(ec core.Rail, req CreateResReq, user common.User) error {
	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(req.Code)
	_, e := lockResourceGlobal(ec, func() (any, error) {
		var id int
		tx := mysql.GetConn().Raw(`select id from resource where code = ? limit 1`, req.Code).Scan(&id)
		if tx.Error != nil {
			return nil, tx.Error
		}

		if id > 0 {
			ec.Debugf("Resource '%s' (%s) already exist", req.Code, req.Name)
			return nil, nil
		}

		res := ERes{
			Name:     req.Name,
			Code:     req.Code,
			CreateBy: user.Username,
			UpdateBy: user.Username,
		}

		tx = mysql.GetConn().
			Table("resource").
			Omit("Id", "CreateTime", "UpdateTime").
			Create(&res)
		return nil, tx.Error
	})
	return e
}

func genPathNo(group string, url string, method string) string {
	cksum := md5.Sum([]byte(group + method + url))
	return "path_" + base64.StdEncoding.EncodeToString(cksum[:])
}

func CreatePathIfNotExist(ec core.Rail, req CreatePathReq, user common.User) error {
	req.Url = preprocessUrl(req.Url)
	req.Group = strings.TrimSpace(req.Group)
	req.Method = strings.ToUpper(strings.TrimSpace(req.Method))
	pathNo := genPathNo(req.Group, req.Url, req.Method)

	res, e := lockPath(ec, pathNo, func() (any, error) {
		var id int
		tx := mysql.GetConn().Raw(`select id from path where path_no = ? limit 1`, pathNo).Scan(&id)
		if tx.Error != nil {
			return false, tx.Error
		}
		if id > 0 { // exists already
			ec.Debugf("Path '%s %s' (%s) already exists", req.Method, req.Url, pathNo)
			return false, nil
		}

		ep := EPath{
			Url:      req.Url,
			Desc:     req.Desc,
			Ptype:    req.Type,
			Pgroup:   req.Group,
			Method:   req.Method,
			PathNo:   pathNo,
			CreateBy: user.Username,
			UpdateBy: user.Username,
		}
		tx = mysql.GetConn().
			Table("path").
			Omit("Id", "CreateTime", "UpdateTime").
			Create(&ep)
		if tx.Error != nil {
			return false, tx.Error
		}

		ec.Infof("Created path (%s) '{%s}'", pathNo, req.Url)
		return true, nil
	})
	if e != nil {
		return e
	}

	created := res.(bool)
	if created { // reload cache for the path
		loadOnePathResCacheAsync(ec, pathNo)
	}

	if req.ResCode != "" { // rebind path and resource
		return BindPathRes(ec, BindPathResReq{PathNo: pathNo, ResCode: req.ResCode})
	}

	return nil
}

func DeletePath(ec core.Rail, req DeletePathReq) error {
	req.PathNo = strings.TrimSpace(req.PathNo)
	_, e := lockPath(ec, req.PathNo, func() (any, error) {
		er := mysql.GetConn().Transaction(func(tx *gorm.DB) error {
			tx = tx.Exec(`delete from path where path_no = ?`, req.PathNo)
			if tx.Error != nil {
				return tx.Error
			}

			return tx.Exec(`delete from path_resource where path_no = ?`, req.PathNo).Error
		})

		return nil, er
	})
	return e
}

func UnbindPathRes(ec core.Rail, req UnbindPathResReq) error {
	req.PathNo = strings.TrimSpace(req.PathNo)
	_, e := lockPath(ec, req.PathNo, func() (any, error) {
		tx := mysql.GetConn().Exec(`delete from path_resource where path_no = ?`, req.PathNo)
		return nil, tx.Error
	})

	if e != nil {
		// asynchronously reload the cache of paths and resources
		go func() {
			if e := LoadPathResCache(ec); e != nil {
				ec.Errorf("Failed to load path resource cache, %v", e)
			}
		}()
	}
	return e
}

func BindPathRes(rail core.Rail, req BindPathResReq) error {
	req.PathNo = strings.TrimSpace(req.PathNo)
	e := lockPathExec(rail, req.PathNo, func() error { // lock for path
		return lockResourceGlobalExec(rail, func() error {

			// check if resource exist
			var resId int
			tx := mysql.GetConn().
				Raw(`SELECT id FROM resource WHERE code = ?`, req.ResCode).
				Scan(&resId)
			if tx.Error != nil {
				return tx.Error
			}
			if resId < 1 {
				rail.Errorf("Resource %v not found", req.ResCode)
				return core.NewWebErr("Resource not found")
			}

			// check if the path is already bound to current resource
			var prid int
			tx = mysql.GetConn().
				Raw(`SELECT id FROM path_resource WHERE path_no = ? AND res_code = ? LIMIT 1`, req.PathNo, req.ResCode).
				Scan(&prid)

			if tx.Error != nil {
				rail.Errorf("Failed to bind path %v to resource %v, %v", req.PathNo, req.ResCode, tx.Error)
				return tx.Error
			}
			if prid > 0 {
				rail.Debugf("Path %v already bound to resource %v", req.PathNo, req.ResCode)
				return tx.Error
			}

			// bind resource to path
			return mysql.GetConn().
				Exec(`INSERT INTO path_resource (path_no, res_code) VALUES (?, ?)`, req.PathNo, req.ResCode).
				Error
		})
	})

	if e == nil {
		// asynchronously reload the cache of paths and resources
		loadOnePathResCacheAsync(rail, req.PathNo)
	}
	return e
}

func ListPaths(ec core.Rail, req ListPathReq) (ListPathResp, error) {

	applyCond := func(t *gorm.DB) *gorm.DB {
		if req.Pgroup != "" {
			t = t.Where("p.pgroup = ?", req.Pgroup)
		}
		if req.ResCode != "" {
			t = t.Joins("LEFT JOIN path_resource pr ON p.path_no = pr.path_no").
				Where("pr.res_code = ?", req.ResCode)
		}
		if req.Url != "" {
			t = t.Where("p.url LIKE ?", "%"+req.Url+"%")
		}
		if req.Ptype != "" {
			t = t.Where("p.ptype = ?", req.Ptype)
		}
		return t
	}

	var paths []WPath
	tx := mysql.GetConn().
		Table("path p").
		Select("p.*").
		Order("id DESC")

	tx = applyCond(tx).
		Offset(req.Paging.GetOffset()).
		Limit(req.Paging.GetLimit()).
		Scan(&paths)
	if tx.Error != nil {
		return ListPathResp{}, tx.Error
	}

	var count int
	tx = mysql.GetConn().
		Table("path p").
		Select("COUNT(*)")

	tx = applyCond(tx).
		Scan(&count)

	if tx.Error != nil {
		return ListPathResp{}, tx.Error
	}

	return ListPathResp{Payload: paths, Paging: core.Paging{Limit: req.Paging.Limit, Page: req.Paging.Page, Total: count}}, nil
}

func AddRole(ec core.Rail, req AddRoleReq, user common.User) error {
	_, e := redis.RLockRun(ec, "goauth:role:add"+req.Name, func() (any, error) {
		r := ERole{
			RoleNo:   core.GenIdP("role_"),
			Name:     req.Name,
			CreateBy: user.Username,
			UpdateBy: user.Username,
		}
		return nil, mysql.GetConn().
			Table("role").
			Omit("Id", "CreateTime", "UpdateTime").
			Create(&r).Error
	})
	return e
}

func RemoveResFromRole(ec core.Rail, req RemoveRoleResReq) error {
	_, e := redis.RLockRun(ec, "goauth:role:"+req.RoleNo, func() (any, error) {
		tx := mysql.GetConn().Exec(`delete from role_resource where role_no = ? and res_code = ?`, req.RoleNo, req.ResCode)
		return nil, tx.Error
	})

	if e != nil {
		e = roleResCache.Put(ec, fmt.Sprintf("role:%s:res:%s", req.RoleNo, req.ResCode), "")
	}

	return e
}

func AddResToRoleIfNotExist(ec core.Rail, req AddRoleResReq, user common.User) error {

	res, e := redis.RLockRun(ec, "goauth:role:"+req.RoleNo, func() (any, error) { // lock for role
		return lockResourceGlobal(ec, func() (any, error) {
			// check if resource exist
			var resId int
			tx := mysql.GetConn().Raw(`select id from resource where code = ?`, req.ResCode).Scan(&resId)
			if tx.Error != nil {
				return false, tx.Error
			}
			if resId < 1 {
				return false, core.NewWebErr("Resource not found")
			}

			// check if role-resource relation exists
			var id int
			tx = mysql.GetConn().Raw(`select id from role_resource where role_no = ? and res_code = ?`, req.RoleNo, req.ResCode).Scan(&id)
			if tx.Error != nil {
				return false, tx.Error
			}
			if id > 0 { // relation exists already
				return false, nil
			}

			// create role-resource relation
			rr := ERoleRes{
				RoleNo:   req.RoleNo,
				ResCode:  req.ResCode,
				CreateBy: user.Username,
				UpdateBy: user.Username,
			}

			return true, mysql.GetConn().
				Table("role_resource").
				Omit("Id", "CreateTime", "UpdateTime").
				Create(&rr).Error
		})
	})

	if e != nil {
		return e
	}

	if isAdded := res.(bool); isAdded {
		e = _loadResOfRole(ec, req.RoleNo)
	}

	return e
}

func ListRoleRes(ec core.Rail, req ListRoleResReq) (ListRoleResResp, error) {
	var res []ListedRoleRes
	tx := mysql.GetConn().
		Raw(`select rr.id, rr.res_code, rr.create_time, rr.create_by, r.name 'res_name' from role_resource rr
			left join resource r on rr.res_code = r.code
			where rr.role_no = ? order by rr.id desc limit ?, ?`, req.RoleNo, req.Paging.GetOffset(), req.Paging.GetLimit()).
		Scan(&res)

	if tx.Error != nil {
		return ListRoleResResp{}, tx.Error
	}

	if res == nil {
		res = []ListedRoleRes{}
	}

	var count int
	tx = mysql.GetConn().
		Raw(`select count(*) from role_resource rr
			left join resource r on rr.res_code = r.code
			where rr.role_no = ?`, req.RoleNo).
		Scan(&count)

	if tx.Error != nil {
		return ListRoleResResp{}, tx.Error
	}

	return ListRoleResResp{Payload: res, Paging: core.Paging{Limit: req.Paging.Limit, Page: req.Paging.Page, Total: count}}, nil
}

func ListAllRoleBriefs(ec core.Rail) ([]RoleBrief, error) {
	var roles []RoleBrief
	tx := mysql.GetConn().Raw("select role_no, name from role").Scan(&roles)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if roles == nil {
		roles = []RoleBrief{}
	}
	return roles, nil
}

func ListRoles(ec core.Rail, req ListRoleReq) (ListRoleResp, error) {
	var roles []WRole
	tx := mysql.GetConn().
		Raw("select * from role order by id desc limit ?, ?", req.Paging.GetOffset(), req.Paging.GetLimit()).
		Scan(&roles)
	if tx.Error != nil {
		return ListRoleResp{}, tx.Error
	}
	if roles == nil {
		roles = []WRole{}
	}

	var count int
	tx = mysql.GetConn().Raw("select count(*) from role").Scan(&count)
	if tx.Error != nil {
		return ListRoleResp{}, tx.Error
	}

	return ListRoleResp{Payload: roles, Paging: core.Paging{Limit: req.Paging.Limit, Page: req.Paging.Page, Total: count}}, nil
}

// Test access to resource
func TestResourceAccess(ec core.Rail, req TestResAccessReq) (TestResAccessResp, error) {
	url := req.Url
	roleNo := req.RoleNo

	if roleNo == DEFAULT_ADMIN_ROLE_NO {
		return permitted, nil
	}

	// some sanitization & standardization for the url
	url = preprocessUrl(url)
	method := strings.ToUpper(strings.TrimSpace(req.Method))

	// find resource required for the url
	cur, e := lookupUrlRes(ec, url, method)
	if e != nil {
		ec.Infof("Rejected '%s' (%s), path not found", url, method)
		return forbidden, nil
	}

	// public path type, doesn't require access to resource
	if cur.Ptype == PT_PUBLIC {
		return permitted, nil
	}

	// doesn't even have role
	roleNo = strings.TrimSpace(roleNo)
	if roleNo == "" {
		ec.Infof("Rejected '%s', user doesn't have roleNo", url)
		return forbidden, nil
	}

	// the requiredRes resources no
	requiredRes := cur.ResCode
	if requiredRes == "" {
		ec.Infof("Rejected '%s', path doesn't have any resource bound yet", url)
		return forbidden, nil
	}

	ok, e := checkRoleRes(ec, roleNo, requiredRes)
	if e != nil {
		return forbidden, e
	}

	// the role doesn't have access to the required resource
	if !ok {
		ec.Infof("Rejected '%s', roleNo: '%s', role doesn't have access to required resource '%s'", url, roleNo, requiredRes)
		return forbidden, nil
	}

	return permitted, nil
}

func checkRoleRes(ec core.Rail, roleNo string, resCode string) (bool, error) {
	r, e := roleResCache.Get(ec, fmt.Sprintf("role:%s:res:%s", roleNo, resCode))
	if e != nil {
		return false, e
	}

	return r != "", nil
}

// Load cache for role -> resources
func LoadRoleResCache(ec core.Rail) error {

	_, e := lockRoleResCache(ec, func() (any, error) {

		lr, e := listRoleNos(ec)
		if e != nil {
			return nil, e
		}

		for _, roleNo := range lr {
			e = _loadResOfRole(ec, roleNo)
			if e != nil {
				return nil, e
			}
		}
		return nil, nil
	})
	return e
}

func _loadResOfRole(ec core.Rail, roleNo string) error {
	roleResList, e := listRoleRes(ec, roleNo)
	if e != nil {
		return e
	}

	for _, rr := range roleResList {
		roleResCache.Put(ec, fmt.Sprintf("role:%s:res:%s", rr.RoleNo, rr.ResCode), "1")
	}
	return nil
}

func listRoleNos(ec core.Rail) ([]string, error) {
	var ern []string
	t := mysql.GetConn().Raw("select role_no from role").Scan(&ern)
	if t.Error != nil {
		return nil, t.Error
	}

	if ern == nil {
		ern = []string{}
	}
	return ern, nil
}

func listRoleRes(ec core.Rail, roleNo string) ([]ERoleRes, error) {
	var rr []ERoleRes
	t := mysql.GetConn().Raw("select * from role_resource where role_no = ?", roleNo).Scan(&rr)
	if t.Error != nil {
		if errors.Is(t.Error, gorm.ErrRecordNotFound) {
			return []ERoleRes{}, nil
		}
		return nil, t.Error
	}

	return rr, nil
}

func lookupUrlRes(ec core.Rail, url string, method string) (CachedUrlRes, error) {
	js, e := urlResCache.Get(ec, method+":"+url)
	if e != nil {
		return CachedUrlRes{}, e
	}
	if js == "" {
		return CachedUrlRes{}, core.NewWebErr(fmt.Sprintf("Unable to find path '%s'", url))
	}

	var cur CachedUrlRes
	if e = json.Unmarshal([]byte(js), &cur); e != nil {
		return CachedUrlRes{}, e
	}

	return cur, nil
}

// Load cache for path -> resource
func LoadPathResCache(ec core.Rail) error {

	_, e := redis.RLockRun(ec, "goauth:path:res:cache", func() (any, error) {

		// ec.Info("Loading path resource cache")

		var paths []ExtendedPathRes
		tx := mysql.GetConn().
			Raw("select p.*, pr.res_code from path p left join path_resource pr on p.path_no = pr.path_no").
			Scan(&paths)
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
			if e := urlResCache.Put(ec, ep.Method+":"+ep.Url, cachedStr); e != nil {
				return nil, e
			}
			// ec.Infof("Loaded PathRes: '%s', '%s', '%s'", ep.Url, ep.Ptype, ep.ResNo)
		}
		return nil, nil
	})

	return e
}

func prepCachedUrlResStr(ec core.Rail, epath ExtendedPathRes) (string, error) {
	url := epath.Url
	cur := CachedUrlRes{
		Id:      epath.Id,
		Pgroup:  epath.Pgroup,
		PathNo:  epath.PathNo,
		ResCode: epath.ResCode,
		Url:     epath.Url,
		Method:  epath.Method,
		Ptype:   epath.Ptype,
	}

	j, e := json.Marshal(cur)
	if e != nil {
		ec.Errorf("Failed to marshal EPath for '%s', %v", url, e)
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
	if ru[l-1] == '/' && l > 1 {
		lj := l - 1
		for lj > 1 && ru[lj-1] == '/' {
			lj -= 1
		}

		ru = ru[0:lj]
	}

	// always start with '/'
	if ru[0] != '/' {
		return "/" + string(ru)
	}
	return string(ru)
}

func findPathRes(pathNo string) (ExtendedPathRes, error) {
	var ep ExtendedPathRes
	tx := mysql.GetConn().
		Raw("select p.*, pr.res_code from path p left join path_resource pr on p.path_no = pr.path_no where p.path_no = ? limit 1", pathNo).
		Scan(&ep)
	if tx.Error != nil {
		return ep, tx.Error
	}

	if tx.RowsAffected < 1 {
		return ep, core.NewWebErr("Path not found")
	}

	return ep, nil
}

// global lock for resources
func lockResourceGlobal(ec core.Rail, runnable redis.LRunnable[any]) (any, error) {
	return redis.RLockRun(ec, "goauth:resource:global", runnable)
}

// global lock for resources
func lockResourceGlobalExec(ec core.Rail, runnable redis.Runnable) error {
	return redis.RLockExec(ec, "goauth:resource:global", runnable)
}

// lock for path
func lockPath(ec core.Rail, pathNo string, runnable redis.LRunnable[any]) (any, error) {
	return redis.RLockRun(ec, "goauth:path:"+pathNo, runnable)
}

// lock for path
func lockPathExec(ec core.Rail, pathNo string, runnable redis.Runnable) error {
	return redis.RLockExec(ec, "goauth:path:"+pathNo, runnable)
}

// lock for role-resource cache
func lockRoleResCache(ec core.Rail, runnable redis.LRunnable[any]) (any, error) {
	return redis.RLockRun(ec, "goauth:role:res:cache", runnable)
}
