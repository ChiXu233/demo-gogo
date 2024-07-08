package apimodel

import (
	"demo-gogo/database/model"
	"demo-gogo/httpserver/errcode"
	"fmt"
	"gorm.io/gorm/utils"
)

const (
	DefaultPageNo   = 1
	DefaultPageSize = 0
	DefaultOrderBy  = "created_at"
	OrderDesc       = "desc"
	OrderAsc        = "asc"

	ValidOptCreateOrUpdate = "save"
	ValidOptList           = "query"
	ValidOptDel            = "del"
)

var (
	DefaultPaginationRequest = PaginationRequest{
		PageNo:   DefaultPageNo,
		PageSize: DefaultPageSize,
		OrderBy:  DefaultOrderBy,
		Order:    OrderDesc,
	}
)

type PaginationRequest struct {
	PageNo   int    `json:"page_no" form:"page_no"`
	PageSize int    `json:"page_size" form:"page_size"`
	OrderBy  string `json:"order_by" form:"order_by"`
	Order    string `json:"order" form:"order"`
}

func (req MapRequest) Valid(opt string) error {
	if opt == ValidOptCreateOrUpdate {
		if req.ID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.Name == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "name")
		}
	} else if opt == ValidOptDel {
		if req.ID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
	} else {
		orderByFields := []string{model.FieldID, model.FieldName, model.FieldCreatedTime, model.FieldUpdatedTime}
		return req.PaginationRequest.Valid(orderByFields)
	}
	return nil
}

func (req PaginationRequest) Valid(orderByList []string) error {
	// pageSize为0代表不分页
	if req.PageSize < 0 {
		return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "[page_size]")
	}
	if req.PageNo <= 0 {
		return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "[page_no]")
	}
	if !utils.Contains(orderByList, req.OrderBy) {
		return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "[order_by]")
	}
	if req.Order != OrderDesc && req.Order != OrderAsc {
		return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "[order]")
	}
	return nil
}

type PaginationResponse struct {
	TotalSize int `json:"total_size"`
}
