package model

import "github.com/lib/pq"

type BaseMap struct {
	Model
	Name           string  `json:"name" gorm:"column:name"`
	MapURL         string  `json:"map_url" gorm:"column:map_url"`
	MapURLCompress string  `json:"map_url_compress" gorm:"column:map_url_compress"`
	PointCloud     string  `json:"point_cloud" gorm:"column:point_cloud"` //点云
	PathID         int     `json:"path_id" gorm:"column:path_id"`         //对应大路径id
	Height         float64 `json:"height"`
	Weight         float64 `json:"weight"`
	Origin         float64 `json:"origin" gorm:"column:origin"`           //z轴起点
	Destination    float64 `json:"destination" gorm:"column:destination"` //z轴终点
}

type MapRoutes struct {
	Model
	RoutesName string          `json:"name" gorm:"column:name"`           //路径名称
	PathID     int             `json:"path_id" gorm:"column:path_id"`     //对应大路径id
	PathRole   string          `json:"path_role" gorm:"column:path_role"` //路径运行规则
	Start      string          `json:"start" gorm:"column:start"`
	End        string          `json:"end" gorm:"column:end"`
	StartToEnd string          `json:"start_end" gorm:"column:start_end"` //运行方向
	EndToStart string          `json:"end_start" gorm:"column:end_start"` //运行方向
	StartRoi   pq.Float64Array `json:"start_roi" gorm:"column:start_roi;type:float8[]"`
	EndRoi     pq.Float64Array `json:"end_roi" gorm:"column:end_point;type:float8[]"`
}

type MapRouteNodes struct {
	Model
	NodeName string          `json:"name" gorm:"column:name" `          //节点名称
	PathID   int             `json:"path_id" gorm:"column:path_id"`     //对应大路径id
	Angle    float64         `json:"angle" gorm:"column:angle"`         //节点角度
	Comment  string          `json:"comment" gorm:"column:comment"`     //标签
	Roi      pq.Float64Array `gorm:"column:roi;type:float8[]" json:"-"` //节点坐标,[33,66]=>(x,y)
}

type Path struct {
	Model
	Name string `json:"name" gorm:"column:name"` //大路径名称
}

func (m *BaseMap) TableName() string {
	return TableNameMap
}
func (m *MapRoutes) TableName() string {
	return TableNameMapRoutes
}
func (m *MapRouteNodes) TableName() string {
	return TableNameMapRouteNodes
}
func (m *Path) TableName() string {
	return TableNamePath
}
