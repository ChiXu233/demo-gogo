package handler

import (
	"demo-gogo/api/apimodel"
	"demo-gogo/httpserver/app"
	"demo-gogo/httpserver/errcode"
	"github.com/gin-gonic/gin"
)

func (handler *RestHandler) CreateOrUpdateMap(c *gin.Context) {
	var req apimodel.MapRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	err = req.Valid(apimodel.ValidOptCreateOrUpdate)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgCreateOrUpdate, err)
		return
	}
	err = handler.Operator.CreateOrUpdateMap(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgCreateOrUpdate, err)
		return
	}
	app.Success(c, nil)
}

func (handler *RestHandler) ListMap(c *gin.Context) {
	req := apimodel.MapRequest{
		PaginationRequest: apimodel.DefaultPaginationRequest,
	}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	err = req.Valid(apimodel.ValidOptList)
	if err != nil {
		app.SendParameterErrorResponse(c, err.Error())
		return
	}
	resp, err := handler.Operator.ListMap(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgListData, err)
		return
	}
	app.Success(c, resp)
}

func (handler *RestHandler) DeleteMap(c *gin.Context) {
	var req apimodel.MapRequest
	err := c.ShouldBindUri(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	err = req.Valid(apimodel.ValidOptDel)
	if err != nil {
		app.SendParameterErrorResponse(c, err.Error())
		return
	}
	err = handler.Operator.DeleteMap(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgDeleteData, err)
		return
	}
	app.Success(c, nil)
}
