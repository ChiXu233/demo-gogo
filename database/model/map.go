package model

import "github.com/lib/pq"

type BaseMap struct {
	Model
	Name           string  `json:"name" gorm:"column:name"`
	MapURL         string  `json:"map_url" gorm:"column:map_url"`
	MapURLCompress string  `json:"map_url_compress" gorm:"column:map_url_compress"`
	Height         float64 `json:"height"`
	Weight         float64 `json:"weight"`
}

type MapRoutes struct {
	Model
	RoutesName string `json:"name" gorm:"column:name"`           //路径名称
	MapID      int    `json:"map_id" gorm:"column:map_id"`       //对应mapID
	PathRole   string `json:"path_role" gorm:"column:path_role"` //路径运行规则
}

type MapRouteNodes struct {
	Model
	NodeName string          `json:"name" gorm:"column:name" `          //节点名称
	MapID    int             `json:"map_id" gorm:"column:map_id"`       //对应mapID
	Angle    string          `json:"angle" gorm:"column:angle"`         //节点角度
	Comment  string          `json:"comment" gorm:"column:comment"`     //标签
	Roi      pq.Float64Array `gorm:"column:roi;type:float8[]" json:"-"` //节点坐标,[33,66]=>(x,y)
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
