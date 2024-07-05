package handler

import (
	"demo-gogo/api/apimodel"
	"demo-gogo/httpserver/app"
	"demo-gogo/httpserver/errcode"
	"github.com/gin-gonic/gin"

	"strconv"
)

func (handler *RestHandler) CreateOrUpdateBaseMap(c *gin.Context) {
	var req apimodel.BaseMapRequest
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
	err = handler.Operator.CreateOrUpdateBaseMap(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgCreateOrUpdate, err)
		return
	}
	app.Success(c, nil)
}

func (handler *RestHandler) ListBaseMap(c *gin.Context) {
	req := apimodel.BaseMapRequest{
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
	resp, err := handler.Operator.ListBaseMap(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgListData, err)
		return
	}
	app.Success(c, resp)
}

func (handler *RestHandler) DeleteBaseMap(c *gin.Context) {
	var req apimodel.BaseMapRequest
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
	err = handler.Operator.DeleteBaseMap(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgDeleteData, err)
		return
	}
	app.Success(c, nil)
}

func (handler *RestHandler) CreateOrUpdateNode(c *gin.Context) {
	var req apimodel.RouteNodesRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	pathId, _ := strconv.Atoi(c.Param("path_id"))
	req.PathID = pathId
	req.NodeName = ""
	err = req.Valid(apimodel.ValidOptCreateOrUpdate)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgCreateOrUpdate, err)
		return
	}

	err = handler.Operator.CreateOrUpdateNode(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgCreateOrUpdate, err)
		return
	}
	app.Success(c, nil)
}

func (handler *RestHandler) ListMapNodes(c *gin.Context) {
	req := apimodel.RouteNodesRequest{
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
	resp, err := handler.Operator.ListMapNodes(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgListData, err)
		return
	}
	app.Success(c, resp)
}

func (handler *RestHandler) DeleteMapNodes(c *gin.Context) {
	var req apimodel.RouteNodesRequest
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
	err = handler.Operator.DeleteMapNodes(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgDeleteData, err)
		return
	}
	app.Success(c, nil)
}

// CreateOrUpdateMapRoutes 接收n个点位信息，将点位按照顺序存储并生成路径
func (handler *RestHandler) CreateOrUpdateMapRoutes(c *gin.Context) {
	var req apimodel.MapRoutesArrRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	pathId, _ := strconv.Atoi(c.Param("path_id"))
	for i := range req.Nodes {
		req.Nodes[i].PathID = pathId
		req.Nodes[i].NodeName = ""
		err = req.Nodes[i].Valid(apimodel.ValidOptCreateOrUpdate)
		if err != nil {
			app.SendParameterErrorResponse(c, err.Error())
			return
		}
	}
	err = handler.Operator.CreateOrUpdateMapRoute(&req)
	if err != nil {
		app.SendServerErrorResponse(c, err.Error(), err)
		return
	}
	app.Success(c, nil)
}
func (handler *RestHandler) ListMapRoutes(c *gin.Context) {
	req := apimodel.MapRoutesRequest{
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
	resp, err := handler.Operator.ListMapRoutes(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgListData, err)
		return
	}
	app.Success(c, resp)
}
func (handler *RestHandler) DeleteMapRoute(c *gin.Context) {
	var req apimodel.MapRoutesRequest
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
	err = handler.Operator.DeleteMapRoute(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgDeleteData, err)
		return
	}
	app.Success(c, nil)
}

func (handler *RestHandler) CreateOrUpdatePath(c *gin.Context) {
	var req apimodel.PathRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	err = req.Valid(apimodel.ValidOptCreateOrUpdate)
	if err != nil {
		app.SendParameterErrorResponse(c, err.Error())
		return
	}
	err = handler.Operator.CreateOrUpdatePath(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgCreateOrUpdate, err)
		return
	}
	app.Success(c, nil)
}

func (handler *RestHandler) ListPath(c *gin.Context) {
	req := apimodel.PathRequest{
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
	resp, err := handler.Operator.ListPath(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgListData, err)
		return
	}
	app.Success(c, resp)
}

func (handler *RestHandler) DeletePath(c *gin.Context) {
	var req apimodel.PathRequest
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
	err = handler.Operator.DeletePath(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgDeleteData, err)
		return
	}
	app.Success(c, nil)
}

// CheckRoute 导航路径校验
func (handler *RestHandler) CheckRoute(c *gin.Context) {
	var req apimodel.MapRoutesArrRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	for _, routes := range req.Routes {
		err = routes.Valid(apimodel.ValidOptCreateOrUpdate)
		if err != nil {
			app.SendParameterErrorResponse(c, err.Error())
			return
		}
	}
	err = handler.Operator.CheckRoute(&req)
	if err != nil {
		app.SendServerErrorResponse(c, err.Error(), err)
		return
	}

	app.Success(c, nil)
}

// ListMapInfo 获取地图切片上所有路径以及点位信息
func (handler *RestHandler) ListMapInfo(c *gin.Context) {
	var req apimodel.RouteNodesRequest
	pathID, _ := strconv.Atoi(c.Param("path_id"))
	req.PathID = pathID
	if req.PathID <= 0 {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	resp, err := handler.Operator.ListMapInfo(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgListData, err)
		return
	}
	app.Success(c, resp)
}

func (handler *RestHandler) BatchDeleteMapNodes(c *gin.Context) {
	var req apimodel.BatchDeleteNodes
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	if len(req.IDs) <= 0 {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgPrefixInvalidParameter)
		return
	}

	err = handler.Operator.BatchDeleteMapNodes(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgDeleteData, err)
		return
	}
	app.Success(c, nil)
}
