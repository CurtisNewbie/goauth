package domain

import "github.com/curtisnewbie/gocommon/common"

type PathType string

const (
	PT_PROTECTED PathType = "PROTECTED"
	PT_PUBLIC    PathType = "PUBLIC"
)

type EPath struct {
	Id         int      // id
	Pgroup     string   // path group
	PathNo     string   // path no
	Desc       string   // description
	ResNo      string   // resource no
	Url        string   // url
	Ptype      PathType // path type: PROTECTED, PUBLIC
	CreateTime common.ETime
	CreateBy   string
	UpdateTime common.ETime
	UpdateBy   string
}

type ERes struct {
	Id         int    // id
	ResNo      string // resource no
	Code       string // resource code
	Name       string // resource name
	CreateTime common.ETime
	CreateBy   string
	UpdateTime common.ETime
	UpdateBy   string
}

type ERoleRes struct {
	Id         int    // id
	RoleNo     string // role no
	ResNo      string // resource no
	CreateTime common.ETime
	CreateBy   string
	UpdateTime common.ETime
	UpdateBy   string
}

type ERole struct {
	Id         int
	RoleNo     string
	Name       string
	CreateTime common.ETime
	CreateBy   string
	UpdateTime common.ETime
	UpdateBy   string
}
