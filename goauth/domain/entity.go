package domain

import "github.com/curtisnewbie/gocommon/common"

type PathType string

const (
	PT_PROTECTED = "PROTECTED"
	PT_PUBLIC   = "PUBLIC"
)

type EPath struct {
	Id         int      // id
	Pgroup     string   // path group
	PathNo     string   // path no
	ResNo      string   // resource no
	Url        string   // url
	Ptype      PathType // path type: PROTECTED, PUBLIC
	CreateTime common.ETime
	CreateBy   string
	UpdateTime common.ETime
	UpdateBy   string
	IsDel      common.IS_DEL
}

type ERes struct {
	Id         int    // id
	ResNo      string // resource no
	Url        string // url
	CreateTime common.ETime
	CreateBy   string
	UpdateTime common.ETime
	UpdateBy   string
	IsDel      common.IS_DEL
}

type ERoleRes struct {
	Id         int    // id
	RoleNo     string // role no
	ResNo      string // resource no
	Url        string // url
	CreateTime common.ETime
	CreateBy   string
	UpdateTime common.ETime
	UpdateBy   string
	IsDel      common.IS_DEL
}

type ERole struct {
	Id         int    // id
	RoleNo     string // role no
	Name       string // role name
	CreateTime common.ETime
	CreateBy   string
	UpdateTime common.ETime
	UpdateBy   string
	IsDel      common.IS_DEL
}
