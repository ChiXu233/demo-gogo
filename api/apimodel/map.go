package apimodel

import (
	"demo-gogo/database/model"
	"demo-gogo/httpserver/errcode"
	"fmt"
	"github.com/lib/pq"
	"math"
)

type BaseMapInfo struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	CreateAt       string  `json:"created_time"`
	UpdateAt       string  `json:"updated_time"`
	PathID         int     `json:"path_id"`
	MapURL         string  `json:"map_url"`
	MapURLCompress string  `json:"map_url_compress"`
	PointCloud     string  `json:"point_cloud"` //点云
	Height         float64 `json:"height"`
	Weight         float64 `json:"weight"`
	Origin         float64 `json:"origin"`      //z轴起点
	Destination    float64 `json:"destination"` //z轴终点
}
type RouteNodesInfo struct {
	ID       int             `json:"id"`
	CreateAt string          `json:"created_time"`
	UpdateAt string          `json:"updated_time"`
	NodeName string          `json:"name"`
	MapID    int             `json:"map_id"`  //对应mapID
	Angle    string          `json:"angle"`   //节点角度
	Comment  string          `json:"comment"` //标签
	Roi      pq.Float64Array `json:"roi"`     //节点坐标,[33,66]=>(x,y)
}
type MapRoutesInfo struct {
	ID         int    `json:"id"`
	CreateAt   string `json:"created_time"`
	UpdateAt   string `json:"updated_time"`
	RoutesName string `json:"name"`      //路径名称
	MapID      int    `json:"map_id"`    //对应mapID
	PathRole   string `json:"path_role"` //路径运行规则
}

type PathInfo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	CreateAt string `json:"created_time"`
	UpdateAt string `json:"updated_time"`
}

type PathRequest struct {
	ID   int    `json:"id" uri:"id" form:"id"`
	Name string `json:"name" form:"name"`
	PaginationRequest
}

type BaseMapRequest struct {
	ID             int     `json:"id" uri:"id" form:"id"`
	Name           string  `json:"name" form:"name"`
	MapURL         string  `json:"map_url"`
	MapURLCompress string  `json:"map_url_compress"`
	PointCloud     string  `json:"point_cloud" gorm:"column:point_cloud"` //点云
	PathID         int     `json:"path_id" form:"path_id"`
	Height         float64 `json:"height"`
	Weight         float64 `json:"weight"`
	Origin         float64 `json:"origin"`      //z轴起点
	Destination    float64 `json:"destination"` //z轴终点
	PaginationRequest
}

type RouteNodesRequest struct {
	ID       int             `json:"id" uri:"id" form:"id"`
	NodeName string          `json:"name" form:"name"`
	MapID    int             `json:"map_id" form:"map_id"` //对应mapID
	Angle    string          `json:"angle"`                //节点角度
	Comment  string          `json:"comment"`              //标签
	Roi      pq.Float64Array `json:"roi"`                  //节点坐标,[33,66]=>(x,y)
	PaginationRequest
}

type MapRoutesArrRequest struct {
	Nodes  []RouteNodesRequest `json:"nodes" form:"nodes"`
	Routes []MapRoutesRequest  `json:"routes" form:"routes"`
}

type MapRoutesRequest struct {
	ID         int    `json:"id" uri:"id" form:"id"`
	RoutesName string `json:"name" form:"name"`     //路径名称
	MapID      int    `json:"map_id" form:"map_id"` //对应mapID
	PathRole   string `json:"path_role"`            //路径运行规则
	PaginationRequest
}

func (m *BaseMapInfo) Load(mapData model.BaseMap) {
	m.ID = mapData.ID
	m.PathID = mapData.PathID
	m.Name = mapData.Name
	m.CreateAt = mapData.CreatedAt.String()
	m.UpdateAt = mapData.UpdatedAt.String()
	m.MapURL = mapData.MapURL
	m.MapURLCompress = mapData.MapURLCompress
	m.Height = mapData.Height
	m.Weight = mapData.Weight
	m.Origin = mapData.Origin
	m.Destination = mapData.Destination
	m.PointCloud = mapData.PointCloud
}

func (m *RouteNodesInfo) Load(nodeData model.MapRouteNodes) {
	m.ID = nodeData.ID
	m.NodeName = nodeData.NodeName
	m.MapID = nodeData.MapID
	m.Angle = nodeData.Angle
	m.Comment = nodeData.Comment
	m.Roi = nodeData.Roi
	m.CreateAt = nodeData.CreatedAt.String()
	m.UpdateAt = nodeData.UpdatedAt.String()
}

func (m *MapRoutesInfo) Load(routeData model.MapRoutes) {
	m.ID = routeData.ID
	m.RoutesName = routeData.RoutesName
	m.MapID = routeData.MapID
	m.PathRole = routeData.PathRole
	m.CreateAt = routeData.CreatedAt.String()
	m.UpdateAt = routeData.UpdatedAt.String()
}

func (m *PathInfo) Load(mapData model.Path) {
	m.ID = mapData.ID
	m.Name = mapData.Name
	m.CreateAt = mapData.CreatedAt.String()
	m.UpdateAt = mapData.UpdatedAt.String()
}

func (req PathRequest) Valid(opt string) error {
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
		orderByFields := []string{model.FieldID, model.FieldPathId, model.FieldName, model.FieldCreatedTime, model.FieldUpdatedTime}
		return req.PaginationRequest.Valid(orderByFields)
	}
	return nil
}

func (req BaseMapRequest) Valid(opt string) error {
	if opt == ValidOptCreateOrUpdate {
		if req.ID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.Name == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "name")
		}
		if req.PathID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "path_id")
		}
	} else if opt == ValidOptDel {
		if req.ID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
	} else {
		orderByFields := []string{model.FieldID, model.FieldPathId, model.FieldName, model.FieldCreatedTime, model.FieldUpdatedTime}
		return req.PaginationRequest.Valid(orderByFields)
	}
	return nil
}

func (req RouteNodesRequest) Valid(opt string) error {
	if opt == ValidOptCreateOrUpdate {
		if req.ID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "node_id")
		}
		if req.NodeName == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "name")
		}
		if req.MapID == 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "map_id")
		}
	} else if opt == ValidOptDel {
		if req.ID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
	} else {
		orderByFields := []string{model.FieldID, model.FieldName, model.FieldMapId, model.FieldCreatedTime, model.FieldUpdatedTime}
		return req.PaginationRequest.Valid(orderByFields)
	}
	return nil
}

func (req MapRoutesRequest) Valid(opt string) error {
	if opt == ValidOptCreateOrUpdate {
		if req.ID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.RoutesName == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "name")
		}
		if req.MapID == 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "map_id")
		}
	} else if opt == ValidOptDel {
		if req.ID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
	} else {
		orderByFields := []string{model.FieldID, model.FieldName, model.FieldMapId, model.FieldCreatedTime, model.FieldUpdatedTime}
		return req.PaginationRequest.Valid(orderByFields)
	}
	return nil
}

type MapPageResponse struct {
	List []BaseMapInfo `json:"list"`
	PaginationResponse
}

type PathResponse struct {
	List []PathInfo `json:"list"`
	PaginationResponse
}

type RouteNodesResponse struct {
	List []RouteNodesInfo `json:"list"`
	PaginationResponse
}

type MapRoutesResponse struct {
	List []MapRoutesInfo `json:"list"`
	PaginationResponse
}

func (resp *MapPageResponse) Load(total int64, list []model.BaseMap) {
	resp.List = make([]BaseMapInfo, 0, len(list))
	for _, v := range list {
		info := BaseMapInfo{}
		info.Load(v)
		resp.List = append(resp.List, info)
	}
	resp.TotalSize = int(total)
}

func (resp *RouteNodesResponse) Load(total int64, list []model.MapRouteNodes) {
	resp.List = make([]RouteNodesInfo, 0, len(list))
	for _, v := range list {
		info := RouteNodesInfo{}
		info.Load(v)
		resp.List = append(resp.List, info)
	}
	resp.TotalSize = int(total)
}

func (resp *MapRoutesResponse) Load(total int64, list []model.MapRoutes) {
	resp.List = make([]MapRoutesInfo, 0, len(list))
	for _, v := range list {
		info := MapRoutesInfo{}
		info.Load(v)
		resp.List = append(resp.List, info)
	}
	resp.TotalSize = int(total)
}

func (resp *PathResponse) Load(total int64, list []model.Path) {
	resp.List = make([]PathInfo, 0, len(list))
	for _, v := range list {
		info := PathInfo{}
		info.Load(v)
		resp.List = append(resp.List, info)
	}
	resp.TotalSize = int(total)
}

func IsPointAbove(p, p1, p2 pq.Float64Array) bool {
	//if (p[0] >= p1[0] && p[0] <= p2[0]) || (p[0] > p2[0] && p[0] < p1[0]) {
	//	if (p[1] >= p1[1] && p[1] <= p2[1]) || (p[1] > p2[1] && p[1] < p1[1]) {
	//		return true
	//	}
	//}
	//return false
	if (p[0] < p1[0] && p[0] > p2[0]) || (p[0] < p2[0] && p[0] > p1[0]) {
		return false
	}
	if (p[1] < p1[1] && p[1] > p2[1]) || (p[1] < p2[1] && p[1] > p1[1]) {
		return false
	}
	slope := math.Abs((p2[1] - p1[1]) / (p2[0] - p1[0]))
	intercept := p1[1] - slope*p1[0]
	expected := slope*p[1] + intercept
	return p[1] == expected

}

func IsPointOnLine(p, p1, p2 pq.Float64Array) bool {
	return (p[0]-p1[0])*(p2[1]-p1[1]) == (p[1]-p1[1])*(p2[0]-p1[0])
}
