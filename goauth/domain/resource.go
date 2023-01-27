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

var (
	urlResCache  = redis.NewLazyRCache(30 * time.Minute) // cache for url's resource, url -> CachedUrlRes
	roleResCache = redis.NewLazyRCache(1 * time.Hour)    // cache for role's resource, role + res -> flag ("1")
)

// Test access to resource
func TestResourceAccess(ec common.ExecContext, url string, roleNo string) error {

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
