package apimodel

import (
	"demo-gogo/database/model"
	"demo-gogo/httpserver/errcode"
	"fmt"
)

type MapInfo struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	CreateAt       string  `json:"created_time"`
	UpdateAt       string  `json:"updated_time"`
	MapURL         string  `json:"map_url"`
	MapURLCompress string  `json:"map_url_compress"`
	Height         float64 `json:"height"`
	Weight         float64 `json:"weight"`
}

type MapRoutes struct {
	ID         int    `json:"id"`
	CreateAt   string `json:"created_time"`
	UpdateAt   string `json:"updated_time"`
	RoutesName string `json:"routes_name"` //路径名称
	MapID      uint   `json:"map_id"`      //对应mapID
	NodesID    []uint `json:"nodes_id"`    //包含路径节点ID
	Path       string `json:"path"`        //路径
	PathRole   string `json:"path_role"`   //路径运行规则
}

type MapRouteNodes struct {
	ID       int    `json:"id"`
	CreateAt string `json:"created_time"`
	UpdateAt string `json:"updated_time"`
	NodeName string `json:"node_name"`
	MapID    uint   `json:"map_id"`    //对应mapID
	Angle    uint   `json:"angle"`     //节点角度
	Comment  string `json:"comment"`   //标签
	RoiArray []uint `json:"roi_array"` //节点坐标,[33,66]=>(x,y)
}

type MapRequest struct {
	ID             int     `json:"id" uri:"id" form:"id"`
	Name           string  `json:"name" form:"name"`
	MapURL         string  `json:"map_url"`
	MapURLCompress string  `json:"map_url_compress"`
	Height         float64 `json:"height"`
	Weight         float64 `json:"weight"`
	PaginationRequest
}

func (m *MapInfo) Load(mapData model.BaseMap) {
	m.ID = mapData.ID
	m.Name = mapData.Name
	m.CreateAt = mapData.CreatedAt.String()
	m.UpdateAt = mapData.UpdatedAt.String()
	m.MapURL = mapData.MapURL
	m.MapURLCompress = mapData.MapURLCompress
	m.Height = mapData.Height
	m.Weight = mapData.Weight
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
		} else {
			orderByFields := []string{model.FieldID, model.FieldName, model.FieldCreatedTime, model.FieldUpdatedTime}
			return req.PaginationRequest.Valid(orderByFields)
		}
	}
	return nil
}

type MapPageResponse struct {
	List []MapInfo `json:"list"`
	PaginationResponse
}

func (resp *MapPageResponse) Load(total int64, list []model.BaseMap) {
	resp.List = make([]MapInfo, 0, len(list))
	for _, v := range list {
		info := MapInfo{}
		info.Load(v)
		resp.List = append(resp.List, info)
	}
	resp.TotalSize = int(total)
}
