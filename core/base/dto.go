package base

import "github.com/pkg/errors"

type ReqBase struct {
	ReqId string `json:"reqId" form:"reqId"` // 请求id 链路跟踪
}

type ReqId struct {
	Id uint64 `json:"id" form:"id"` // 主键ID
}

type ReqStrId struct {
	Id string `json:"id" form:"id"` // 主键ID
}

type ReqIds struct {
	Ids []uint64 `json:"ids" form:"ids[]"` //多id
}

func (x *ReqIds) Valid() error {
	if len(x.Ids) == 0 {
		return errors.New("ids cannot be empty")
	}
	return nil
}

type ReqPage struct {
	Page     int `json:"page" form:"page" query:"-"`         // 页码
	PageSize int `json:"pageSize" form:"pageSize" query:"-"` // 每页大小
}

func (e *ReqPage) GetPage() int {
	if e.Page < 1 {
		return 1
	}
	return e.Page
}

func (e *ReqPage) GetSize() int {
	if e.PageSize < 1 {
		return 10
	}
	return e.PageSize
}

func (e *ReqPage) GetOffset() int {
	return (e.GetPage() - 1) * e.GetSize()
}
