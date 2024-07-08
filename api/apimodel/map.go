package apimodel

import (
	"demo-gogo/database/model"
	"demo-gogo/httpserver/errcode"
	"fmt"
	"github.com/lib/pq"
	"math"
)

type MapInfo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_time"`
	UpdatedAt string `json:"updated_time"`
}

type MapInfoInfo struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	CreateAt       string  `json:"created_time"`
	UpdateAt       string  `json:"updated_time"`
	MapID          int     `json:"map_id"`
	MapURL         string  `json:"map_url"`
	MapURLCompress string  `json:"map_url_compress"`
	PointCloud     string  `json:"point_cloud"` //点云
	Origin         float64 `json:"origin"`      //z轴起点
	Destination    float64 `json:"destination"` //z轴终点
}
type RouteNodesInfo struct {
	ID       int             `json:"id"`
	CreateAt string          `json:"created_time"`
	UpdateAt string          `json:"updated_time"`
	NodeName string          `json:"name"`
	InfoID   int             `json:"info_id"`
	Angle    float64         `json:"angle"`   //节点角度
	Comment  string          `json:"comment"` //标签
	Roi      pq.Float64Array `json:"roi"`     //节点坐标,[33,66]=>(x,y)
}
type MapRoutesInfo struct {
	ID         int             `json:"id"`
	CreateAt   string          `json:"created_time"`
	UpdateAt   string          `json:"updated_time"`
	RoutesName string          `json:"name"` //路径名称
	InfoID     int             `json:"info_id"`
	Start      string          `json:"start"`      //起点
	End        string          `json:"end" `       //终点
	StartToEnd string          `json:"start_end"`  //运行方向
	EndToStart string          `json:"end_start"`  //运行方向
	PathRole   string          `json:"path_role"`  //路径运行规则
	StartRoi   pq.Float64Array `json:"start_roi" ` //起点坐标
	EndRoi     pq.Float64Array `json:"end_roi"`    //终点坐标
}

type MapInfoRequest struct {
	ID             int     `json:"id" uri:"id" form:"id"`
	Name           string  `json:"name" form:"name"`
	MapURL         string  `json:"map_url"`
	MapURLCompress string  `json:"map_url_compress"`
	PointCloud     string  `json:"point_cloud" gorm:"column:point_cloud"` //点云
	MapID          int     `json:"map_id" form:"map_id"`
	Origin         float64 `json:"origin"`      //z轴起点
	Destination    float64 `json:"destination"` //z轴终点
	PaginationRequest
}

type RouteNodesRequest struct {
	ID       int             `json:"id" uri:"id" form:"id"`
	NodeName string          `json:"name" form:"name"`
	InfoID   int             `json:"info_id" form:"info_id"`
	Angle    float64         `json:"angle"`   //节点角度
	Comment  string          `json:"comment"` //标签
	Roi      pq.Float64Array `json:"roi"`     //节点坐标,[33,66]=>(x,y)
	PaginationRequest
}

type MapRoutesArrRequest struct {
	Nodes  []RouteNodesRequest `json:"nodes" form:"nodes"`
	Routes []MapRoutesRequest  `json:"routes" form:"routes"`
}

type MapRequest struct {
	ID   int    `json:"id" uri:"id" form:"id"`
	Name string `json:"name" form:"name"`
	PaginationRequest
}

type BatchDeleteNodes struct {
	IDs []int `json:"node_ids"`
}

type MapRoutesRequest struct {
	ID         int    `json:"id" uri:"id" form:"id"`
	RoutesName string `json:"name" form:"name"` //路径名称
	InfoID     int    `json:"info_id" form:"info_id"`
	PathRole   string `json:"path_role"`                         //路径运行规则
	Start      string `json:"start"`                             //起点
	End        string `json:"end" `                              //终点
	StartToEnd string `json:"start_end" gorm:"column:start_end"` //运行方向
	EndToStart string `json:"end_start" gorm:"column:end_start"` //运行方向
	PaginationRequest
}

func (m *MapInfoInfo) Load(mapData model.MapInfo) {
	m.ID = mapData.ID
	m.MapID = mapData.MapID
	m.Name = mapData.Name
	m.CreateAt = mapData.CreatedAt.String()
	m.UpdateAt = mapData.UpdatedAt.String()
	m.MapURL = mapData.MapURL
	m.MapURLCompress = mapData.MapURLCompress
	m.Origin = mapData.Origin
	m.Destination = mapData.Destination
	m.PointCloud = mapData.PointCloud
}

func (m *RouteNodesInfo) Load(nodeData model.MapRouteNodes) {
	m.ID = nodeData.ID
	m.NodeName = nodeData.NodeName
	m.InfoID = nodeData.InfoID
	m.Angle = nodeData.Angle
	m.Comment = nodeData.Comment
	m.Roi = nodeData.Roi
	m.CreateAt = nodeData.CreatedAt.String()
	m.UpdateAt = nodeData.UpdatedAt.String()
}

func (m *MapRoutesInfo) Load(routeData model.MapRoutes) {
	m.ID = routeData.ID
	m.RoutesName = routeData.RoutesName
	m.InfoID = routeData.InfoID
	m.PathRole = routeData.PathRole
	m.CreateAt = routeData.CreatedAt.String()
	m.UpdateAt = routeData.UpdatedAt.String()
	m.Start = routeData.Start
	m.End = routeData.End
	m.StartToEnd = routeData.StartToEnd
	m.EndToStart = routeData.EndToStart
	m.StartRoi = routeData.StartRoi
	m.EndRoi = routeData.EndRoi
}

func (req MapInfoRequest) Valid(opt string) error {
	if opt == ValidOptCreateOrUpdate {
		if req.ID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.Name == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "name")
		}
		if req.MapID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "map_id")
		}
	} else if opt == ValidOptDel {
		if req.ID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
	} else {
		orderByFields := []string{model.FieldID, model.FieldInfoId, model.FieldName, model.FieldCreatedTime, model.FieldUpdatedTime}
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
		if req.InfoID == 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "info_id")
		}
	} else if opt == ValidOptDel {
		if req.ID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
	} else {
		orderByFields := []string{model.FieldID, model.FieldName, model.FieldInfoId, model.FieldCreatedTime, model.FieldUpdatedTime}
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
		if req.InfoID == 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "info_id")
		}
	} else if opt == ValidOptDel {
		if req.ID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
	} else {
		orderByFields := []string{model.FieldID, model.FieldName, model.FieldInfoId, model.FieldCreatedTime, model.FieldUpdatedTime}
		return req.PaginationRequest.Valid(orderByFields)
	}
	return nil
}

type MapInfoPageResponse struct {
	List []MapInfoInfo `json:"list"`
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

type MapPageResponse struct {
	List []MapInfo `json:"list"`
	PaginationResponse
}

type MapInfosResponse struct {
	Nodes  []RouteNodesInfo `json:"nodes"`
	Routes []MapRoutesInfo  `json:"routes"`
}

func (m *MapInfo) Load(mapData model.Map) {
	m.ID = mapData.ID
	m.Name = mapData.Name
	m.CreatedAt = mapData.CreatedAt.String()
	m.UpdatedAt = mapData.UpdatedAt.String()
}

func (resp *MapInfoPageResponse) Load(total int64, list []model.MapInfo) {
	resp.List = make([]MapInfoInfo, 0, len(list))
	for _, v := range list {
		info := MapInfoInfo{}
		info.Load(v)
		resp.List = append(resp.List, info)
	}
	resp.TotalSize = int(total)
}

func (resp *MapPageResponse) Load(total int64, list []model.Map) {
	resp.List = make([]MapInfo, 0, len(list))
	for _, v := range list {
		info := MapInfo{}
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

func (resp *MapInfosResponse) Load(routes []model.MapRoutes, nodes []model.MapRouteNodes) {
	resp.Routes = make([]MapRoutesInfo, 0)
	resp.Nodes = make([]RouteNodesInfo, 0)
	for _, v := range routes {
		info := MapRoutesInfo{}
		info.Load(v)
		resp.Routes = append(resp.Routes, info)
	}
	for _, v := range nodes {
		info := RouteNodesInfo{}
		info.Load(v)
		resp.Nodes = append(resp.Nodes, info)
	}
}

// IsPointAbove 判断点是否在直线上
func IsPointAbove(p, p1, p2 pq.Float64Array) bool {
	if (p[0] == p1[0] && p[1] == p1[1]) || (p[0] == p2[0] && p[1] == p2[1]) {
		return false
	}
	k := (p2[1] - p1[1]) / (p2[0] - p1[0])
	b := p1[1] - k*p1[0]
	y := k*p[0] + b
	return p[1] == y
}

func PointToLine(p, p1, p2 pq.Float64Array) (bool, pq.Float64Array) {
	var point pq.Float64Array
	//px, py, x1, y1, x2, y2
	// 计算直线 Ax + By + C = 0 的 A, B, C
	//A := y2 - y1
	A := p2[1] - p1[1]
	//B := x1 - x2
	B := p1[0] - p2[0]
	//C := x2*y1 - x1*y2
	C := p2[0]*p1[1] - p1[0]*p2[1]

	// 计算距离公式的分子
	numerator := math.Abs(A*p[0] + B*p[1] + C)
	// 计算距离公式的分母
	denominator := math.Sqrt(A*A + B*B)
	// 计算距离
	distance := numerator / denominator

	if distance < 4 {
		point = perpendicularIntersection(p, p1, p2)
		return true, point
	} else {
		return false, nil
	}

}

// 计算点与直线距离最近点坐标
func perpendicularIntersection(p, p1, p2 pq.Float64Array) pq.Float64Array {
	//直线平行y轴方向
	if p1[0] == p2[0] {
		return pq.Float64Array{p1[0], p[1]}
		//return Point{x: p1.x, y: p.y}
	}
	// 直线平行x轴方向
	if p1[1] == p2[1] {
		return pq.Float64Array{p[0], p1[1]}
		//return Point{x: p.x, y: p1.y}
	}
	// 计算直线AB的斜率
	slope := (p2[1] - p1[1]) / (p2[0] - p1[0])
	// 垂线的斜率为 -1/m
	mPerpendicular := -1 / slope
	// 垂线方程的截距 b = y - mPerpendicular * x
	bPerpendicular := p[1] - mPerpendicular*p[0]
	// 直线AB方程的截距 b = y - m * x
	b := p1[1] - slope*p1[0]
	// 求交点的x坐标
	intersectX := (bPerpendicular - b) / (slope - mPerpendicular)
	// 求交点的y坐标
	intersectY := slope*intersectX + b
	return pq.Float64Array{intersectX, intersectY}
}

func IsPointOnLine(p, p1, p2 pq.Float64Array) bool {
	return (p[0]-p1[0])*(p2[1]-p1[1]) == (p[1]-p1[1])*(p2[0]-p1[0])
}

// PointsAbove 返回线段上的所有点
func PointsAbove(p1, p2 pq.Float64Array) [][]float64 {
	var x1, x2, y1, y2 float64
	var point [][]float64
	if p1[0] > p2[0] {
		x1, y1 = p2[0], p2[1]
		x2, y2 = p1[0], p1[1]
	} else {
		x1, y1 = p1[0], p1[1]
		x2, y2 = p2[0], p2[1]
	}

	if x1 == x2 {
		for y := y1; y <= y2; y++ {
			point = append(point, []float64{x1, y})
		}
	} else if y1 == y2 {
		for x := x1; x <= x2; x++ {
			point = append(point, []float64{x, y1})
		}
	} else {
		slope := (y2 - y1) / (x2 - x1)
		for x := x1; x <= x2; x++ {
			y := math.Round(y1 + slope*(x-x1))
			point = append(point, []float64{x, y})
		}
	}
	return point
}
