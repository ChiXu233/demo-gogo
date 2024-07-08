package service

import (
	"demo-gogo/api/apimodel"
	"demo-gogo/database/model"
	"demo-gogo/httpserver/errcode"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/lib/pq"
	log "github.com/wonderivan/logger"
	"gorm.io/gorm"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"strconv"
	"strings"
)

func (operator *ResourceOperator) CreateOrUpdateMap(req *apimodel.MapRequest) error {
	var opt model.Map
	selector := make(map[string]interface{})
	// 名称唯一性
	selector[model.FieldName] = req.Name
	err := operator.Database.ListEntityByFilter(model.TableNameMap, selector, model.OneQuery, &opt)
	if err != nil {
		return err
	}
	if opt.ID != 0 && opt.ID != req.ID {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "地图")
	}
	if req.ID > 0 {
		err = operator.Database.GetEntityByID(model.TableNameMap, req.ID, &opt)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "待修改地图")
			}
			return err
		}
	}
	err = copier.Copy(&opt, req)
	if err != nil {
		return err
	}
	if req.ID > 0 {
		err = operator.Database.SaveEntity(model.TableNameMap, &opt)
		if err != nil {
			log.Error("地图数据更新失败. err:[%v]", err)
			return err
		}
	} else {
		err = operator.Database.CreateEntity(model.TableNameMap, &opt)
		if err != nil {
			log.Error("地图数据创建失败. err:[%v]", err)
			return err
		}
	}
	return nil
}

func (operator *ResourceOperator) ListMap(req *apimodel.MapRequest) (*apimodel.MapPageResponse, error) {
	var resp apimodel.MapPageResponse
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	if req.ID > 0 {
		selector[model.FieldID] = req.ID
	}
	if req.Name != "" {
		selector[model.FieldName] = req.Name
	}
	var count int64
	var maps []model.Map
	err := operator.Database.CountEntityByFilter(model.TableNameMap, selector, model.OneQuery, &count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		order := model.Order{
			Field:     model.FieldID,
			Direction: apimodel.OrderAsc,
		}
		queryParams.Orders = append(queryParams.Orders, order)
		if req.PageSize > 0 {
			queryParams.Limit = &req.PageSize
			offset := (req.PageNo - 1) * req.PageSize
			queryParams.Offset = &offset
		}
		err = operator.Database.ListEntityByFilter(model.TableNameMap, selector, queryParams, &maps)
		if err != nil {
			log.Error("地图数据查询失败. err:[%v]", err)
			return nil, err
		}
	}
	resp.Load(count, maps)
	return &resp, nil
}

func (operator *ResourceOperator) DeleteMap(req *apimodel.MapRequest) error {
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	selector[model.FieldID] = req.ID
	err := operator.Database.DeleteEntityByFilter(model.TableNameMap, selector, queryParams, &model.Map{})
	if err != nil {
		log.Error("地图数据删除失败. err:[%v]", err)
		return err
	}
	return nil
}

func (operator *ResourceOperator) CreateOrUpdateMapInfo(req *apimodel.MapInfoRequest) error {
	var opt model.MapInfo
	var mapDB model.Map
	selector := make(map[string]interface{})
	// 名称唯一性
	selector[model.FieldName] = req.Name
	err := operator.Database.ListEntityByFilter(model.TableNameMapInfo, selector, model.OneQuery, &opt)
	if err != nil {
		return err
	}
	if opt.ID != 0 && opt.ID != req.ID {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "地图信息")
	}

	//if (req.Height == 0 && req.Weight == 0) && strings.HasSuffix(path.Base(req.MapURL), ".png") {
	//	file, err := os.Open(req.MapURL)
	//	if err != nil {
	//		return err
	//	}
	//	defer file.Close()
	//	img, err := png.Decode(file)
	//	if err != nil {
	//		return err
	//	}
	//	bounds := img.Bounds()
	//	req.Weight = float64(bounds.Dx())
	//	req.Height = float64(bounds.Dy())
	//}
	if req.ID > 0 {
		err = operator.Database.GetEntityByID(model.TableNameMapInfo, req.ID, &opt)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "待修改地图信息")
			}
			return err
		}
	}
	err = operator.Database.GetEntityByID(model.TableNameMap, req.MapID, &mapDB)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "关联地图")
		}
		return err
	}
	err = copier.Copy(&opt, req)
	if err != nil {
		return err
	}
	if req.ID > 0 {
		err = operator.Database.SaveEntity(model.TableNameMapInfo, &opt)
		if err != nil {
			log.Error("地图信息数据更新失败. err:[%v]", err)
			return err
		}
	} else {
		err = operator.Database.CreateEntity(model.TableNameMapInfo, &opt)
		if err != nil {
			log.Error("地图信息数据创建失败. err:[%v]", err)
			return err
		}
	}
	return nil
}

func (operator *ResourceOperator) ListMapInfoPageResponse(req *apimodel.MapInfoRequest) (*apimodel.MapInfoPageResponse, error) {
	var resp apimodel.MapInfoPageResponse
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	if req.ID > 0 {
		selector[model.FieldID] = req.ID
	}
	if req.Name != "" {
		selector[model.FieldName] = req.Name
	}
	if req.MapID > 0 {
		selector[model.FieldMapId] = req.MapID
	}
	var count int64
	var maps []model.MapInfo
	err := operator.Database.CountEntityByFilter(model.TableNameMapInfo, selector, model.OneQuery, &count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		order := model.Order{
			Field:     req.OrderBy,
			Direction: apimodel.OrderAsc,
		}
		queryParams.Orders = append(queryParams.Orders, order)
		if req.PageSize > 0 {
			queryParams.Limit = &req.PageSize
			offset := (req.PageNo - 1) * req.PageSize
			queryParams.Offset = &offset
		}
		err = operator.Database.ListEntityByFilter(model.TableNameMapInfo, selector, queryParams, &maps)
		if err != nil {
			log.Error("地图数据查询失败,err:[%v]", err)
			return nil, err
		}
	}
	resp.Load(count, maps)
	return &resp, nil
}

func (operator *ResourceOperator) DeleteMapInfo(req *apimodel.MapInfoRequest) error {
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	selector[model.FieldID] = req.ID
	err := operator.Database.DeleteEntityByFilter(model.TableNameMapInfo, selector, queryParams, &model.MapInfo{})
	if err != nil {
		log.Error("地图数据删除失败. err:[%v]", err)
		return err
	}
	return nil
}

func (operator *ResourceOperator) CreateOrUpdateNode(req *apimodel.RouteNodesRequest) error {
	var opt model.MapRouteNodes
	var mapList model.MapInfo
	var startNode model.MapRouteNodes
	var routeCreate []model.MapRoutes
	selector := make(map[string]interface{})
	// 开启事务
	tx, err := operator.TransactionBegin()
	if err != nil {
		log.Error("CreateOrUpdateTrainType TransactionBegin Error.err[%v]", err)
		return err
	}
	defer func() {
		_ = tx.TransactionRollback()
	}()

	//验证map_id是否存在
	err = operator.Database.GetEntityByID(model.TableNameMapInfo, req.InfoID, &mapList)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "地图关联路径节点")
		}
		return err
	}
	//同一map_id中名称的唯一性
	selector[model.FieldInfoId] = req.InfoID
	selector[model.FieldName] = req.NodeName
	err = operator.Database.ListEntityByFilter(model.TableNameMapRouteNodes, selector, model.OneQuery, &opt)
	if err != nil {
		return err
	}
	if opt.ID != 0 && opt.ID != req.ID {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "地图路径节点")
	}

	selector = make(map[string]interface{})
	selector[model.FieldInfoId] = req.InfoID
	//维护nodeName字段
	params := model.OneQuery
	order := model.Order{
		Field:     model.FieldName,
		Direction: apimodel.OrderDesc,
	}
	params.Orders = append(params.Orders, order)
	err = tx.Database.ListEntityByFilter(model.TableNameMapRouteNodes, selector, params, &startNode)
	if err != nil {
		return err
	}
	nameNumber := 0
	if startNode.ID > 0 {
		nameNumber, _ = strconv.Atoi(strings.Split(startNode.NodeName, "Site")[1])
	}
	if nameNumber+1 < 10 {
		req.NodeName = "Site" + "000" + strconv.Itoa(nameNumber+1)
	} else if nameNumber+1 < 100 && nameNumber+1 >= 10 {
		req.NodeName = "Site" + "00" + strconv.Itoa(nameNumber+1)
	} else if nameNumber+1 < 1000 && nameNumber+1 >= 100 {
		req.NodeName = "Site" + "0" + strconv.Itoa(nameNumber+1)
	}
	if req.ID > 0 {
		err = operator.Database.GetEntityByID(model.TableNameMapRouteNodes, req.ID, &opt)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "待修改地图路径节点")
			}
			return err
		}
		req.NodeName = opt.NodeName
	}
	//nodeName不允许编辑

	err = copier.Copy(&opt, req)
	if err != nil {
		return err
	}

	//判断新增节点坐标是否处在任意一条线上
	var routes []model.MapRoutes
	var nodes, nodeRois []model.MapRouteNodes
	nameMap := make(map[string]struct{})
	selector = make(map[string]interface{})
	selector[model.FieldInfoId] = req.InfoID
	if err = operator.Database.ListEntityByFilter(model.TableNameMapRoutes, selector, model.QueryParams{}, &routes); err != nil {
		return err
	}
	err = operator.Database.GetEntityPluck(model.TableNameMapRouteNodes, selector, model.QueryParams{}, "name", &nodes)
	if err != nil {
		return err
	}
	err = operator.Database.GetEntityPluck(model.TableNameMapRouteNodes, selector, model.QueryParams{}, "roi", &nodeRois)
	if err != nil {
		return err
	}

	for i := range nodes {
		nodes[i].Roi = nodeRois[i].Roi
	}
	for _, v := range routes {
		nameMap[v.RoutesName] = struct{}{}
	}
	for _, v := range routes {
		//拿到路径的起始点、末尾点坐标
		var roiH, roiE pq.Float64Array
		nodeHeadName := v.Start
		nodeEndName := v.End
		for _, k := range nodes {
			if k.NodeName == nodeHeadName {
				roiH = k.Roi
			}
			if k.NodeName == nodeEndName {
				roiE = k.Roi
			}
		}

		if roiH != nil && roiE != nil {
			flag, point := apimodel.PointToLine(req.Roi, roiH, roiE)
			if flag {
				req.Roi = point
				//同一线段上
				routeA := model.MapRoutes{RoutesName: nodeHeadName + "-" + req.NodeName, InfoID: req.InfoID, PathRole: "双向", Start: nodeHeadName, End: req.NodeName, StartToEnd: "正向行走", EndToStart: "正向行走"}
				routeB := model.MapRoutes{RoutesName: req.NodeName + "-" + nodeEndName, InfoID: req.InfoID, PathRole: "双向", Start: req.NodeName, End: nodeEndName, StartToEnd: "正向行走", EndToStart: "正向行走"}
				if _, ok := nameMap[routeA.RoutesName]; ok {
					continue
				}
				nameMap[routeA.RoutesName] = struct{}{}
				routeCreate = append(routeCreate, routeA)
				if _, ok := nameMap[routeB.RoutesName]; ok {
					continue
				}
				nameMap[routeB.RoutesName] = struct{}{}
				routeCreate = append(routeCreate, routeB)
			}
		}
	}

	if routeCreate != nil {
		err = operator.BatchCreateEntity(model.TableNameMapRoutes, routeCreate)
		if err != nil {
			log.Error("地图路径节点创建失败,err:[%v]", err)
			return err
		}
	}

	if req.ID > 0 {
		err = operator.Database.SaveEntity(model.TableNameMapRouteNodes, &opt)
		if err != nil {
			log.Error("地图路径节点更新失败,err:[%v]", err)
			return err
		}
	} else {
		err = operator.Database.CreateEntity(model.TableNameMapRouteNodes, &opt)
		if err != nil {
			log.Error("地图路径节点创建失败,err:[%v]", err)
			return err
		}
	}
	return nil
}

func (operator *ResourceOperator) ListMapNodes(req *apimodel.RouteNodesRequest) (*apimodel.RouteNodesResponse, error) {
	var resp apimodel.RouteNodesResponse
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	if req.ID > 0 {
		selector[model.FieldID] = req.ID
	}
	if req.InfoID > 0 {
		selector[model.FieldInfoId] = req.InfoID
	}
	if req.NodeName != "" {
		selector[model.FieldName] = req.NodeName
	}
	var count int64
	var nodes []model.MapRouteNodes
	err := operator.Database.CountAllEntityByFilter(model.TableNameMapRouteNodes, selector, model.QueryParams{}, &count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		order := model.Order{
			Field:     req.OrderBy,
			Direction: apimodel.OrderAsc,
		}
		queryParams.Orders = append(queryParams.Orders, order)
		if req.PageSize > 0 {
			queryParams.Limit = &req.PageSize
			offset := (req.PageNo - 1) * req.PageSize
			queryParams.Offset = &offset
		}
		err = operator.Database.ListEntityByFilter(model.TableNameMapRouteNodes, selector, queryParams, &nodes)
		if err != nil {
			log.Error("地图关联节点查询失败,err:[%v]", err)
			return nil, err
		}
	}
	resp.Load(count, nodes)
	return &resp, nil
}

// DeleteMapNodes 删除路径节点并删除对应路径并生成新路径
func (operator *ResourceOperator) DeleteMapNodes(req *apimodel.RouteNodesRequest) error {
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	var routes []model.MapRoutes
	var routeCreate []model.MapRoutes
	var node model.MapRouteNodes
	// 开启事务
	tx, err := operator.TransactionBegin()
	if err != nil {
		log.Error("CreateOrUpdateTrainType TransactionBegin Error.err[%v]", err)
		return err
	}
	defer func() {
		_ = tx.TransactionRollback()
	}()
	selector[model.FieldID] = req.ID
	err = tx.Database.ListEntityByFilter(model.TableNameMapRouteNodes, selector, model.OneQuery, &node)
	if err != nil {
		log.Error("查找节点数据失败,err:[%v]", err)
		return err
	}
	if node.ID <= 0 {
		return fmt.Errorf("待删除节点不存在")
	}
	err = tx.Database.DeleteEntityByFilter(model.TableNameMapRouteNodes, selector, queryParams, &model.MapRouteNodes{})
	if err != nil {
		log.Error("删除地图关联节点失败,err:[%v]", err)
		return err
	}

	//删除关联的路径
	selector = make(map[string]interface{})
	selector[model.FieldInfoId] = node.InfoID
	err = tx.Database.ListEntityByFilter(model.TableNameMapRoutes, selector, queryParams, &routes)
	if err != nil {
		log.Error("查找节点关联路径失败,err:[%v]", err)
		return err
	}
	for i := range routes {
		if routes[i].Start == node.NodeName || routes[i].End == node.NodeName {
			selector[model.FieldID] = routes[i].ID
			err = tx.Database.DeleteEntityByFilter(model.TableNameMapRoutes, selector, queryParams, &model.MapRoutes{})
			if err != nil {
				return err
			}
		}
	}

	//生成新路径
	for _, v := range routes {
		if v.End == node.NodeName {
			for _, k := range routes {
				if k.Start == node.NodeName {
					route := model.MapRoutes{RoutesName: v.Start + "-" + k.End, InfoID: v.InfoID, PathRole: "双向", Start: v.Start, End: k.End, StartToEnd: "正向行走", EndToStart: "正向行走"}
					routeCreate = append(routeCreate, route)
				}
			}
		}
	}
	//生成新路径
	if len(routeCreate) > 0 {
		for _, v := range routeCreate {
			err = tx.Database.CreateEntity(model.TableNameMapRoutes, &v)
			if err != nil {
				log.Error("DeleteMapRoute TransactionCommit Error.err[%v]", err)
				return err
			}
		}
	}
	err = tx.TransactionCommit()
	if err != nil {
		log.Error("CreateOrUpdateTrainType TransactionCommit Error.err[%v]", err)
		return err
	}
	return nil
}

//接收n个点位信息，将点位按照顺序存储并生成路径

func (operator *ResourceOperator) CreateOrUpdateMapRoute(req *apimodel.MapRoutesArrRequest) error {
	selector := make(map[string]interface{})
	var createNodes []model.MapRouteNodes
	var updateNodes []model.MapRouteNodes
	var nodes []model.MapRouteNodes
	var createRoutes []model.MapRoutes
	var updateRoutes []model.MapRoutes
	var nameStart model.MapRouteNodes

	//对相邻且相同坐标节点过滤
	for i := 0; i+1 < len(req.Nodes); i++ {
		if req.Nodes[i].Roi[0] == req.Nodes[i+1].Roi[0] && req.Nodes[i].Roi[1] == req.Nodes[i+1].Roi[1] {
			if i+2 >= len(req.Nodes) {
				req.Nodes = req.Nodes[:i+1]
			} else {
				req.Nodes = append(req.Nodes[:i+1], req.Nodes[i+2:]...)
			}
		}
	}
	// 开启事务
	tx, err := operator.TransactionBegin()
	if err != nil {
		log.Error("CreateOrUpdateTrainType TransactionBegin Error.err[%v]", err)
		return err
	}
	defer func() {
		_ = tx.TransactionRollback()
	}()

	//自动生成name字段

	params := model.OneQuery
	order := model.Order{
		Field:     model.FieldName,
		Direction: apimodel.OrderDesc,
	}
	params.Orders = append(params.Orders, order)
	err = tx.Database.ListEntityByFilter(model.TableNameMapRouteNodes, selector, params, &nameStart)
	if err != nil {
		return err
	}
	nameNumber := 0
	if nameStart.ID > 0 {
		nameNumber, _ = strconv.Atoi(strings.Split(nameStart.NodeName, "Site")[1])
	}

	if req.Nodes != nil {
		index := 1
		for _, v := range req.Nodes {
			var node model.MapRouteNodes
			//验证路径节点正确性
			selector[model.FieldName] = v.NodeName
			selector[model.FieldInfoId] = v.InfoID
			err = tx.Database.ListEntityByFilter(model.TableNameMapRouteNodes, selector, model.OneQuery, &node)
			if err != nil {
				return err
			}
			if node.ID != 0 && node.ID != v.ID {
				return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "该地图节点")
			}
			if v.ID > 0 {
				err = tx.Database.GetEntityByID(model.TableNameMapRouteNodes, v.ID, &node)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "待复用路径节点")
					}
					return err
				}
				//不允许编辑名称
				v.NodeName = node.NodeName
				err = copier.Copy(&node, v)
				if err != nil {
					return err
				}
				updateNodes = append(updateNodes, node)
			} else {
				err = copier.Copy(&node, v)
				if err != nil {
					return err
				}
				if nameNumber+index < 10 {
					node.NodeName = "Site" + "000" + strconv.Itoa(nameNumber+index)
				} else if nameNumber+index < 100 && nameNumber+index >= 10 {
					node.NodeName = "Site" + "00" + strconv.Itoa(nameNumber+index)
				} else if nameNumber+index < 1000 && nameNumber+index >= 100 {
					node.NodeName = "Site" + "0" + strconv.Itoa(nameNumber+index)
				}
				createNodes = append(createNodes, node)
				index++
			}
			nodes = append(nodes, node)
		}

		if updateNodes != nil {
			for _, v := range updateNodes {
				err = tx.Database.SaveEntity(model.TableNameMapRouteNodes, &v)
				if err != nil {
					log.Error("地图路径节点更新失败. err:[%v]", err)
				}
			}
		}

		err = tx.Database.BatchCreateEntity(model.TableNameMapRouteNodes, createNodes)
		if err != nil {
			log.Error("地图路径节点创建失败. err:[%v]", err)
			return err
		}
	}

	//不传routes，自动生成对应路径
	for i := 0; i < len(nodes)-1; i++ {
		route := model.MapRoutes{
			RoutesName: nodes[i].NodeName + "-" + nodes[i+1].NodeName,
			InfoID:     req.Nodes[0].InfoID,
			//Path:       utils.String(nodes[i].ID) + "-" + utils.String(nodes[i+1].ID),
			Start:      nodes[i].NodeName,
			End:        nodes[i+1].NodeName,
			PathRole:   "双向",
			StartToEnd: "正向行走",
			EndToStart: "正向行走",
		}

		var routeIndex model.MapRoutes
		selector = make(map[string]interface{})
		selector[model.FieldName] = route.RoutesName
		selector[model.FieldInfoId] = route.InfoID
		err = tx.Database.ListEntityByFilter(model.TableNameMapRoutes, selector, model.OneQuery, &routeIndex)
		if err != nil {
			return err
		}
		if routeIndex.ID != 0 {
			continue
		}
		createRoutes = append(createRoutes, route)
	}

	if len(req.Routes) != 0 {
		//以传入的routes为准
		createRoutes = nil
		//编辑路径名称以及路径规则
		for _, v := range req.Routes {
			var route model.MapRoutes
			//名称唯一性
			selector = make(map[string]interface{})
			selector[model.FieldName] = v.RoutesName
			selector[model.FieldInfoId] = v.InfoID
			err = tx.Database.ListEntityByFilter(model.TableNameMapRoutes, selector, model.OneQuery, &route)
			if err != nil {
				return err
			}
			if route.ID != 0 && route.ID != v.ID {
				return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "该地图路径")
			}
			if v.ID > 0 {
				err = tx.Database.GetEntityByID(model.TableNameMapRoutes, v.ID, &route)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "待编辑路径")
					}
					return err
				}
				err = copier.Copy(&route, v)
				if err != nil {
					return err
				}
				updateRoutes = append(updateRoutes, route)
			} else {
				err = copier.Copy(&route, v)
				if err != nil {
					return err
				}
				createRoutes = append(createRoutes, route)
			}
		}
	}

	if req.Nodes != nil {
		//校正是否有节点交叉左右联通
		nameMap := make(map[string]struct{})
		roiMap := make(map[string]model.MapRouteNodes)
		var routes []model.MapRoutes
		var nodeNames, nodeRois []model.MapRouteNodes
		selector = make(map[string]interface{})
		selector[model.FieldInfoId] = req.Nodes[0].InfoID
		if err = tx.Database.ListEntityByFilter(model.TableNameMapRoutes, selector, model.QueryParams{}, &routes); err != nil {
			return err
		}
		err = tx.Database.GetEntityPluck(model.TableNameMapRouteNodes, selector, model.QueryParams{}, "name", &nodeNames)
		if err != nil {
			return err
		}
		err = tx.Database.GetEntityPluck(model.TableNameMapRouteNodes, selector, model.QueryParams{}, "roi", &nodeRois)
		if err != nil {
			return err
		}
		//
		for i := range nodeNames {
			nodeNames[i].Roi = nodeRois[i].Roi
		}
		routes = append(routes, createRoutes...)
		routes = append(routes, updateRoutes...)
		for _, route := range routes {
			nameMap[route.RoutesName] = struct{}{}
		}
		for _, node := range nodeNames {
			roiMap[node.NodeName] = node
		}
		for i := range nodes {
			for _, v := range routes {
				//拿到路径的起始点、末尾点坐标
				var roiH, roiE pq.Float64Array
				nodeHeadName := v.Start
				//strings.Split(v.RoutesName, "-")[0]
				nodeEndName := v.End
				//strings.Split(v.RoutesName, "-")[1]
				roiH = roiMap[nodeHeadName].Roi
				roiE = roiMap[nodeEndName].Roi
				if roiH != nil && roiE != nil {
					flag, point := apimodel.PointToLine(nodes[i].Roi, roiH, roiE)
					if flag {
						nodes[i].Roi = point
						//同一线段上
						routeA := model.MapRoutes{RoutesName: nodeHeadName + "-" + nodes[i].NodeName, InfoID: nodes[i].InfoID, PathRole: "双向", Start: nodeHeadName, End: nodes[i].NodeName, StartToEnd: "正向行走", EndToStart: "正向行走"}
						routeB := model.MapRoutes{RoutesName: nodes[i].NodeName + "-" + nodeEndName, InfoID: nodes[i].InfoID, PathRole: "双向", Start: nodes[i].NodeName, End: nodeEndName, StartToEnd: "正向行走", EndToStart: "正向行走"}
						if _, ok := nameMap[routeA.RoutesName]; ok {
							continue
						}
						nameMap[routeA.RoutesName] = struct{}{}
						createRoutes = append(createRoutes, routeA)
						if _, ok := nameMap[routeB.RoutesName]; ok {
							continue
						}
						nameMap[routeB.RoutesName] = struct{}{}
						createRoutes = append(createRoutes, routeB)
					}
				}
			}
		}
	}
	for _, v := range updateRoutes {
		err = tx.Database.SaveEntity(model.TableNameMapRoutes, &v)
		if err != nil {
			log.Error("地图路径更新失败. err:[%v]", err)
			return err
		}
	}

	err = tx.Database.BatchCreateEntity(model.TableNameMapRoutes, createRoutes)

	if err != nil {
		log.Error("地图路径创建失败. err:[%v]", err)
		return err
	}

	err = tx.TransactionCommit()
	if err != nil {
		log.Error("CreateOrUpdateTrainType TransactionCommit Error.err[%v]", err)
		return err
	}
	return nil
}

func (operator *ResourceOperator) ListMapRoutes(req *apimodel.MapRoutesRequest) (*apimodel.MapRoutesResponse, error) {
	var resp apimodel.MapRoutesResponse
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	if req.ID > 0 {
		selector[model.FieldID] = req.ID
	}
	if req.RoutesName != "" {
		selector[model.FieldName] = req.RoutesName
	}
	if req.InfoID != 0 {
		selector[model.FieldInfoId] = req.InfoID
	}
	var count int64
	var maps []model.MapRoutes
	err := operator.Database.CountEntityByFilter(model.TableNameMapRoutes, selector, model.OneQuery, &count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		order := model.Order{
			Field:     req.OrderBy,
			Direction: apimodel.OrderAsc,
		}
		queryParams.Orders = append(queryParams.Orders, order)
		if req.PageSize > 0 {
			queryParams.Limit = &req.PageSize
			offset := (req.PageNo - 1) * req.PageSize
			queryParams.Offset = &offset
		}
		err = operator.Database.ListEntityByFilter(model.TableNameMapRoutes, selector, queryParams, &maps)
		if err != nil {
			log.Error("地图数据查询失败. err:[%v]", err)
			return nil, err
		}
	}
	resp.Load(count, maps)
	return &resp, nil
}

// DeleteMapRoute 删除路径并删除对应路径节点
func (operator *ResourceOperator) DeleteMapRoute(req *apimodel.MapRoutesRequest) error {
	var routes []model.MapRoutes
	var route model.MapRoutes
	selector := make(map[string]interface{})
	exitMap := make(map[string]struct{})
	var routeCreate []model.MapRoutes
	queryParams := model.QueryParams{}
	// 开启事务
	tx, err := operator.TransactionBegin()
	if err != nil {
		log.Error("CreateOrUpdateTrainType TransactionBegin Error.err[%v]", err)
		return err
	}
	defer func() {
		_ = tx.TransactionRollback()
	}()
	selector[model.FieldID] = req.ID
	err = tx.Database.ListEntityByFilter(model.TableNameMapRoutes, selector, model.OneQuery, &route)
	if err != nil {
		log.Error("查找路径数据失败, err:[%v]", err)
		return err
	}
	if route.ID <= 0 {
		return fmt.Errorf("待删除路径不存在")
	}
	err = tx.Database.DeleteEntityByFilter(model.TableNameMapRoutes, selector, queryParams, &model.MapRoutes{})
	if err != nil {
		log.Error("地图数据删除失败. err:[%v]", err)
		return err
	}
	selector = make(map[string]interface{})
	selector[model.FieldInfoId] = route.InfoID
	err = tx.Database.ListEntityByFilter(model.TableNameMapRoutes, selector, queryParams, &routes)
	if err != nil {
		log.Error("查找路径数据失败,err:[%v]", err)
		return err
	}
	for _, v := range routes {
		exitMap[v.Start] = struct{}{}
		exitMap[v.End] = struct{}{}
	}
	if _, ok := exitMap[route.Start]; !ok {
		selector[model.FieldName] = route.Start
		err = tx.Database.DeleteEntityByFilter(model.TableNameMapRouteNodes, selector, queryParams, &model.MapRouteNodes{})
		if err != nil {
			log.Error("删除路径数据失败,err:[%v]", err)
			return err
		}
	}
	//if _, ok := exitMap[route.End]; !ok {
	selector[model.FieldName] = route.End
	err = tx.Database.DeleteEntityByFilter(model.TableNameMapRouteNodes, selector, queryParams, &model.MapRouteNodes{})
	if err != nil {
		log.Error("删除路径节点数据失败,err:[%v]", err)
		return err
	}
	//}
	//删除以route.end为起点的route
	for _, v := range routes {
		if v.Start == route.End {
			tempRoute := model.MapRoutes{RoutesName: route.Start + "-" + v.End, InfoID: route.InfoID, PathRole: "双向", Start: route.Start, End: v.End, StartToEnd: "正向行走", EndToStart: "正向行走"}
			routeCreate = append(routeCreate, tempRoute)
		}
	}
	selector = make(map[string]interface{})
	selector["start"] = route.End
	err = tx.Database.DeleteEntityByFilter(model.TableNameMapRoutes, selector, queryParams, &model.MapRouteNodes{})
	if err != nil {
		log.Error("删除路径数据失败,err:[%v]", err)
		return err
	}
	//生成新路径
	if len(routeCreate) > 0 {
		for _, v := range routeCreate {
			err = tx.Database.CreateEntity(model.TableNameMapRoutes, &v)
			if err != nil {
				log.Error("DeleteMapRoute TransactionCommit Error.err[%v]", err)
				return err
			}
		}
	}
	err = tx.TransactionCommit()
	if err != nil {
		log.Error("DeleteMapRoute TransactionCommit Error.err[%v]", err)
		return err
	}
	return nil
}

func (operator *ResourceOperator) CheckRoute(req *apimodel.MapRoutesArrRequest) error {
	selector := make(map[string]interface{})
	var nodeHead, nodeEnd model.MapRouteNodes
	var mapInfo model.MapInfo
	selector[model.FieldID] = req.Routes[0].InfoID
	err := operator.Database.ListEntityByFilter(model.TableNameMapInfo, selector, model.OneQuery, &mapInfo)
	if err != nil {
		log.Error("查询地图信息失败,err:[%v]", err)
		return err
	}
	mapInfo.MapURL = "/Users/dg2023/Desktop/test_photo.png"
	file, err := os.Open(mapInfo.MapURL)
	if err != nil {
		log.Error("读取文件失败", err)
		return err
	}
	defer file.Close()

	for _, route := range req.Routes {
		var mapRoute model.MapRoutes
		selector = make(map[string]interface{})
		selector[model.FieldInfoId] = route.InfoID
		selector[model.FieldName] = route.RoutesName
		err = operator.Database.ListEntityByFilter(model.TableNameMapRoutes, selector, model.OneQuery, &mapRoute)
		if err != nil {
			log.Error("路径数据查找失败,err:[%v]", err)
			return err
		}
		if mapRoute.ID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "待校验路径")
		}
		nodeHeadName := strings.Split(route.RoutesName, "-")[0]
		nodeEndName := strings.Split(route.RoutesName, "-")[1]
		selector[model.FieldName] = nodeHeadName
		err = operator.Database.ListEntityByFilter(model.TableNameMapRouteNodes, selector, model.OneQuery, &nodeHead)
		if err != nil {
			log.Error("节点数据查询失败,err:[%v]", err)
			return err
		}
		selector[model.FieldName] = nodeEndName
		err = operator.Database.ListEntityByFilter(model.TableNameMapRouteNodes, selector, model.OneQuery, &nodeEnd)
		if err != nil {
			log.Error("节点数据查询失败,err:[%v]", err)
			return err
		}
		points := apimodel.PointsAbove(nodeHead.Roi, nodeEnd.Roi)

		if strings.HasSuffix(path.Base(mapInfo.MapURL), ".png") {
			img, err := png.Decode(file)
			if err != nil {
				log.Error("读取png文件失败,err:[%v]", err)
				return err
			}
			for _, v := range points {
				grayColor := color.GrayModel.Convert(img.At(int(v[0]), int(v[1]))).(color.Gray)
				if grayColor.Y > 150 {
					log.Error("路径：[%v] 校验未通过", route.RoutesName)
					return fmt.Errorf("路径：[%v] 校验未通过", route.RoutesName)
				}
			}
		} else if strings.HasSuffix(path.Base(mapInfo.MapURL), ".jpg") {
			img, err := jpeg.Decode(file)
			if err != nil {
				log.Error("读取png文件失败,err:[%v]", err)
				return fmt.Errorf("读取png文件失败")
			}
			for _, v := range points {
				grayColor := color.GrayModel.Convert(img.At(int(v[0]), int(v[1]))).(color.Gray)
				if grayColor.Y > 150 {
					log.Error("路径：[%v] 校验未通过", route.RoutesName)
					return fmt.Errorf("路径：[%v] 校验未通过", route.RoutesName)
				}
			}
		}
	}
	return nil
}

func (operator *ResourceOperator) ListMapInfo(req *apimodel.RouteNodesRequest) (*apimodel.MapInfosResponse, error) {
	var resp apimodel.MapInfosResponse
	selector := make(map[string]interface{})
	var routes []model.MapRoutes
	var nodes []model.MapRouteNodes
	queryParams := model.QueryParams{}
	order := model.Order{
		Field:     model.FieldID,
		Direction: apimodel.OrderAsc,
	}
	queryParams.Orders = append(queryParams.Orders, order)
	selector[model.FieldInfoId] = req.InfoID
	err := operator.Database.ListEntityByFilter(model.TableNameMapRoutes, selector, queryParams, &routes)
	if err != nil {
		log.Error("路径数据查询失败. err:[%v]", err)
		return nil, err
	}

	err = operator.Database.ListEntityByFilter(model.TableNameMapRouteNodes, selector, queryParams, &nodes)
	if err != nil {
		log.Error("节点数据查询失败. err:[%v]", err)
		return nil, err
	}
	roiMap := make(map[string]pq.Float64Array)
	for _, v := range nodes {
		roiMap[v.NodeName] = v.Roi
	}
	for i := range routes {
		routes[i].StartRoi = roiMap[routes[i].Start]
		routes[i].EndRoi = roiMap[routes[i].End]
	}
	resp.Load(routes, nodes)
	return &resp, nil
}

func (operator *ResourceOperator) BatchDeleteMapNodes(req *apimodel.BatchDeleteNodes) error {
	var err error
	tx, err := operator.TransactionBegin()
	if err != nil {
		log.Error("BatchDeletePlan TransactionBegin Error.err:[%#v]", err)
		return err
	}
	defer func() {
		_ = tx.TransactionRollback()
	}()
	var nodes []model.MapRouteNodes
	var routes []model.MapRoutes
	var ids []int
	nameMap := make(map[string]struct{})
	filter := make(map[string]interface{})
	queryParams := model.QueryParams{}
	inQuery := model.InQuery{
		Field:  model.FieldID,
		Values: req.IDs,
	}
	queryParams.InQueries = append(queryParams.InQueries, &inQuery)
	err = tx.Database.ListEntityByFilter(model.TableNameMapRouteNodes, filter, queryParams, &nodes)
	if err != nil {
		return err
	}
	if len(nodes) <= 0 {
		return fmt.Errorf("查找节点数据失败")
	}
	//删除节点
	err = tx.Database.DeleteEntityByFilter(model.TableNameMapRouteNodes, filter, queryParams, &model.MapRouteNodes{})
	if err != nil {
		return err
	}
	//删除与节点关联路径
	filter[model.FieldInfoId] = nodes[0].InfoID
	err = tx.Database.ListEntityByFilter(model.TableNameMapRoutes, filter, model.QueryParams{}, &routes)
	if err != nil {
		return err
	}
	for _, v := range nodes {
		nameMap[v.NodeName] = struct{}{}
	}
	for _, v := range routes {
		if _, ok := nameMap[v.Start]; ok {
			ids = append(ids, v.ID)
		}
		if _, ok := nameMap[v.End]; ok {
			ids = append(ids, v.ID)
		}
	}
	queryParams = model.QueryParams{}
	inQuery = model.InQuery{
		Field:  model.FieldID,
		Values: ids,
	}
	queryParams.InQueries = append(queryParams.InQueries, &inQuery)
	err = tx.Database.DeleteEntityByFilter(model.TableNameMapRoutes, filter, queryParams, &model.MapRoutes{})
	if err != nil {
		return err
	}
	err = tx.TransactionCommit()
	if err != nil {
		log.Error("CreateOrUpdateTrainType TransactionCommit Error.err[%v]", err)
		return err
	}
	return nil
}
