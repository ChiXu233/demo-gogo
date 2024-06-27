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
	"strings"
)

func (operator *ResourceOperator) CreateOrUpdateBaseMap(req *apimodel.BaseMapRequest) error {
	var opt model.BaseMap
	var pathMap model.Path
	selector := make(map[string]interface{})
	// 名称唯一性
	selector[model.FieldName] = req.Name
	err := operator.Database.ListEntityByFilter(model.TableNameBaseMap, selector, model.OneQuery, &opt)
	if err != nil {
		return err
	}
	if opt.ID != 0 && opt.ID != req.ID {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "地图")
	}
	//if !utils.Exists(req.MapURL) {
	//	err = fmt.Errorf("路径文件不存在")
	//	return err
	//}
	if (req.Height == 0 && req.Weight == 0) && strings.HasSuffix(path.Base(req.MapURL), ".png") {
		file, err := os.Open(req.MapURL)
		if err != nil {
			return err
		}
		defer file.Close()
		img, err := png.Decode(file)
		if err != nil {
			return err
		}
		bounds := img.Bounds()
		req.Weight = float64(bounds.Dx())
		req.Height = float64(bounds.Dy())
	}
	if req.ID > 0 {
		err = operator.Database.GetEntityByID(model.TableNameBaseMap, req.ID, &opt)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "待修改地图")
			}
			return err
		}
	}
	err = operator.Database.GetEntityByID(model.TableNamePath, req.PathID, &pathMap)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "关联大路径")
		}
		return err
	}
	err = copier.Copy(&opt, req)
	if err != nil {
		return err
	}
	if req.ID > 0 {
		err = operator.Database.SaveEntity(model.TableNameBaseMap, &opt)
		if err != nil {
			log.Error("地图数据更新失败. err:[%v]", err)
			return err
		}
	} else {
		err = operator.Database.CreateEntity(model.TableNameBaseMap, &opt)
		if err != nil {
			log.Error("地图数据创建失败. err:[%v]", err)
			return err
		}
	}
	return nil
}

func (operator *ResourceOperator) ListBaseMap(req *apimodel.BaseMapRequest) (*apimodel.MapPageResponse, error) {
	var resp apimodel.MapPageResponse
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	if req.ID > 0 {
		selector[model.FieldID] = req.ID
	}
	if req.Name != "" {
		selector[model.FieldName] = req.Name
	}
	if req.PathID > 0 {
		selector[model.FieldPathId] = req.PathID
	}
	var count int64
	var maps []model.BaseMap
	err := operator.Database.CountEntityByFilter(model.TableNameBaseMap, selector, model.OneQuery, &count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		order := model.Order{
			Field:     req.OrderBy,
			Direction: req.Order,
		}
		queryParams.Orders = append(queryParams.Orders, order)
		if req.PageSize > 0 {
			queryParams.Limit = &req.PageSize
			offset := (req.PageNo - 1) * req.PageSize
			queryParams.Offset = &offset
		}
		err = operator.Database.ListEntityByFilter(model.TableNameBaseMap, selector, queryParams, &maps)
		if err != nil {
			log.Error("地图数据查询失败,err:[%v]", err)
			return nil, err
		}
	}
	resp.Load(count, maps)
	return &resp, nil
}

func (operator *ResourceOperator) DeleteBaseMap(req *apimodel.BaseMapRequest) error {
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	selector[model.FieldID] = req.ID
	err := operator.Database.DeleteEntityByFilter(model.TableNameBaseMap, selector, queryParams, &model.BaseMap{})
	if err != nil {
		log.Error("地图数据删除失败. err:[%v]", err)
		return err
	}
	return nil
}

func (operator *ResourceOperator) CreateOrUpdateNode(req *apimodel.RouteNodesRequest) error {
	var opt model.MapRouteNodes
	var mapList model.BaseMap
	var routeCreate []model.MapRoutes
	selector := make(map[string]interface{})
	//验证map_id是否存在
	selector[model.FieldID] = req.PathID
	err := operator.Database.GetEntityByID(model.TableNameMap, req.PathID, &mapList)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "地图关联路径节点")
		}
		return err
	}
	//同一map_id中名称的唯一性
	selector = make(map[string]interface{})
	selector[model.FieldPathId] = req.PathID
	selector[model.FieldName] = req.NodeName
	err = operator.Database.ListEntityByFilter(model.TableNameMapRouteNodes, selector, model.OneQuery, &opt)
	if err != nil {
		return err
	}
	if opt.ID != 0 && opt.ID != req.ID {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "地图路径节点")
	}

	if req.ID > 0 {
		err = operator.Database.GetEntityByID(model.TableNameMapRouteNodes, req.ID, &opt)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "待修改地图路径节点")
			}
			return err
		}
	}
	err = copier.Copy(&opt, req)
	if err != nil {
		return err
	}

	//判断新增节点坐标是否处在任意一条线上
	var routes []model.MapRoutes
	var nodes, nodeRois []model.MapRouteNodes
	nameMap := make(map[string]struct{})
	selector = make(map[string]interface{})
	selector[model.FieldPathId] = req.PathID
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
		nodeHeadName := strings.Split(v.RoutesName, "-")[0]
		nodeEndName := strings.Split(v.RoutesName, "-")[1]
		for _, k := range nodes {
			if k.NodeName == nodeHeadName {
				roiH = k.Roi
			}
			if k.NodeName == nodeEndName {
				roiE = k.Roi
			}
		}

		if (roiH != nil && roiE != nil) && apimodel.IsPointAbove(req.Roi, roiH, roiE) {
			//同一线段上
			routeA := model.MapRoutes{RoutesName: nodeHeadName + "-" + req.NodeName, PathID: req.PathID, PathRole: "双向"}
			routeB := model.MapRoutes{RoutesName: req.NodeName + "-" + nodeEndName, PathID: req.PathID, PathRole: "双向"}
			if _, ok := nameMap[routeA.RoutesName]; ok {
				continue
			}
			routeCreate = append(routeCreate, routeA)
			if _, ok := nameMap[routeB.RoutesName]; ok {
				continue
			}
			routeCreate = append(routeCreate, routeB)
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
	if req.PathID > 0 {
		selector[model.FieldPathId] = req.PathID
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
			Direction: req.Order,
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

func (operator *ResourceOperator) DeleteMapNodes(req *apimodel.RouteNodesRequest) error {
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	selector[model.FieldID] = req.ID
	err := operator.Database.DeleteEntityByFilter(model.TableNameMapRouteNodes, selector, queryParams, &model.MapRouteNodes{})
	if err != nil {
		log.Error("删除地图关联节点失败,err:[%v]", err)
		return err
	}
	return nil
}

func (operator *ResourceOperator) CreateOrUpdateMapRoute(req *apimodel.MapRoutesArrRequest) error {
	selector := make(map[string]interface{})
	var createNodes []model.MapRouteNodes
	var updateNodes []model.MapRouteNodes
	var nodes []model.MapRouteNodes
	var createRoutes []model.MapRoutes
	var updateRoutes []model.MapRoutes
	// 开启事务
	tx, err := operator.TransactionBegin()
	if err != nil {
		log.Error("CreateOrUpdateTrainType TransactionBegin Error.err[%v]", err)
		return err
	}
	defer func() {
		_ = tx.TransactionRollback()
	}()
	if req.Nodes != nil {
		for _, v := range req.Nodes {
			var node model.MapRouteNodes
			//验证路径节点正确性
			selector[model.FieldName] = v.NodeName
			selector[model.FieldPathId] = v.PathID
			err := operator.Database.ListEntityByFilter(model.TableNameMapRouteNodes, selector, model.OneQuery, &node)
			if err != nil {
				return err
			}
			if node.ID != 0 && node.ID != v.ID {
				return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "该地图节点")
			}
			if v.ID > 0 {
				err = operator.Database.GetEntityByID(model.TableNameMapRouteNodes, v.ID, &node)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "待复用路径节点")
					}
					return err
				}
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
				createNodes = append(createNodes, node)
			}
			nodes = append(nodes, node)
		}

		if updateNodes != nil {
			for _, v := range updateNodes {
				err = operator.Database.SaveEntity(model.TableNameMapRouteNodes, &v)
				if err != nil {
					log.Error("地图路径节点更新失败. err:[%v]", err)
				}
			}
		}

		err = operator.Database.BatchCreateEntity(model.TableNameMapRouteNodes, createNodes)
		if err != nil {
			log.Error("地图路径节点创建失败. err:[%v]", err)
			return err
		}
	}

	//不传routes，自动生成对应路径
	for i := 0; i < len(nodes)-1; i++ {
		route := model.MapRoutes{
			RoutesName: nodes[i].NodeName + "-" + nodes[i+1].NodeName,
			PathID:     req.Nodes[0].PathID,
			//Path:       utils.String(nodes[i].ID) + "-" + utils.String(nodes[i+1].ID),
			PathRole: "双向",
		}

		var routeIndex model.MapRoutes
		selector = make(map[string]interface{})
		selector[model.FieldName] = route.RoutesName
		selector[model.FieldPathId] = route.PathID
		err = operator.Database.ListEntityByFilter(model.TableNameMapRoutes, selector, model.OneQuery, &routeIndex)
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
			selector[model.FieldPathId] = v.PathID
			err = operator.Database.ListEntityByFilter(model.TableNameMapRoutes, selector, model.OneQuery, &route)
			if err != nil {
				return err
			}
			if route.ID != 0 && route.ID != v.ID {
				return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "该地图路径")
			}
			if v.ID > 0 {
				err = operator.Database.GetEntityByID(model.TableNameMapRoutes, v.ID, &route)
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
		selector[model.FieldPathId] = req.Nodes[0].PathID
		if err = operator.Database.ListEntityByFilter(model.TableNameMapRoutes, selector, model.QueryParams{}, &routes); err != nil {
			return err
		}
		err = operator.Database.GetEntityPluck(model.TableNameMapRouteNodes, selector, model.QueryParams{}, "name", &nodeNames)
		if err != nil {
			return err
		}
		err = operator.Database.GetEntityPluck(model.TableNameMapRouteNodes, selector, model.QueryParams{}, "roi", &nodeRois)
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
		for _, point := range nodes {
			for _, v := range routes {
				//拿到路径的起始点、末尾点坐标
				var roiH, roiE pq.Float64Array
				nodeHeadName := strings.Split(v.RoutesName, "-")[0]
				nodeEndName := strings.Split(v.RoutesName, "-")[1]
				roiH = roiMap[nodeHeadName].Roi
				roiE = roiMap[nodeEndName].Roi
				if (roiH != nil && roiE != nil) && apimodel.IsPointAbove(point.Roi, roiH, roiE) {
					//同一线段上
					routeA := model.MapRoutes{RoutesName: nodeHeadName + "-" + point.NodeName, PathID: point.PathID, PathRole: "双向"}
					routeB := model.MapRoutes{RoutesName: point.NodeName + "-" + nodeEndName, PathID: point.PathID, PathRole: "双向"}
					if _, ok := nameMap[routeA.RoutesName]; ok {
						continue
					}
					createRoutes = append(createRoutes, routeA)
					if _, ok := nameMap[routeB.RoutesName]; ok {
						continue
					}
					createRoutes = append(createRoutes, routeB)
				}
			}
		}
	}
	for _, v := range updateRoutes {
		err = operator.Database.SaveEntity(model.TableNameMapRoutes, &v)
		if err != nil {
			log.Error("地图路径更新失败. err:[%v]", err)
			return err
		}
	}

	err = operator.Database.BatchCreateEntity(model.TableNameMapRoutes, createRoutes)
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
	if req.PathID != 0 {
		selector[model.FieldPathId] = req.PathID
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
			Direction: req.Order,
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

func (operator *ResourceOperator) DeleteMapRoute(req *apimodel.MapRoutesRequest) error {
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	selector[model.FieldID] = req.ID
	err := operator.Database.DeleteEntityByFilter(model.TableNameMapRoutes, selector, queryParams, &model.MapRoutes{})
	if err != nil {
		log.Error("地图数据删除失败. err:[%v]", err)
		return err
	}
	return nil
}

func (operator *ResourceOperator) CreateOrUpdatePath(req *apimodel.PathRequest) error {
	var opt model.Path
	selector := make(map[string]interface{})
	// 名称唯一性
	selector[model.FieldName] = req.Name
	err := operator.Database.ListEntityByFilter(model.TableNamePath, selector, model.OneQuery, &opt)
	if err != nil {
		return err
	}
	if opt.ID != 0 && opt.ID != req.ID {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "路径")
	}
	if req.ID > 0 {
		err = operator.Database.GetEntityByID(model.TableNamePath, req.ID, &opt)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "待修改路径")
			}
			return err
		}
	}
	err = copier.Copy(&opt, req)
	if err != nil {
		return err
	}
	if req.ID > 0 {
		err = operator.Database.SaveEntity(model.TableNamePath, &opt)
		if err != nil {
			log.Error("路径数据更新失败. err:[%v]", err)
			return err
		}
	} else {
		err = operator.Database.CreateEntity(model.TableNamePath, &opt)
		if err != nil {
			log.Error("路径数据创建失败. err:[%v]", err)
			return err
		}
	}
	return nil
}

func (operator *ResourceOperator) ListPath(req *apimodel.PathRequest) (*apimodel.PathResponse, error) {
	var resp apimodel.PathResponse
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	if req.ID > 0 {
		selector[model.FieldID] = req.ID
	}
	if req.Name != "" {
		selector[model.FieldName] = req.Name
	}
	var count int64
	var maps []model.Path
	err := operator.Database.CountEntityByFilter(model.TableNamePath, selector, model.OneQuery, &count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		order := model.Order{
			Field:     req.OrderBy,
			Direction: req.Order,
		}
		queryParams.Orders = append(queryParams.Orders, order)
		if req.PageSize > 0 {
			queryParams.Limit = &req.PageSize
			offset := (req.PageNo - 1) * req.PageSize
			queryParams.Offset = &offset
		}
		err = operator.Database.ListEntityByFilter(model.TableNamePath, selector, queryParams, &maps)
		if err != nil {
			log.Error("路径数据查询失败. err:[%v]", err)
			return nil, err
		}
	}
	resp.Load(count, maps)
	return &resp, nil
}

func (operator *ResourceOperator) DeletePath(req *apimodel.PathRequest) error {
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	selector[model.FieldID] = req.ID
	err := operator.Database.DeleteEntityByFilter(model.TableNamePath, selector, queryParams, &model.Path{})
	if err != nil {
		log.Error("路径数据删除失败. err:[%v]", err)
		return err
	}
	return nil
}

func (operator *ResourceOperator) CheckRoute(req *apimodel.MapRoutesArrRequest) error {
	selector := make(map[string]interface{})
	var nodeHead, nodeEnd model.MapRouteNodes
	var mapInfo model.BaseMap
	selector[model.FieldID] = req.Routes[0].PathID
	err := operator.Database.ListEntityByFilter(model.TableNameBaseMap, selector, model.OneQuery, &mapInfo)
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
		selector[model.FieldPathId] = route.PathID
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
		fmt.Println(points)

		if strings.HasSuffix(path.Base(mapInfo.MapURL), ".png") {
			img, err := png.Decode(file)
			if err != nil {
				log.Error("读取png文件失败,err:[%v]", err)
				return err
			}
			for _, v := range points {
				grayColor := color.GrayModel.Convert(img.At(int(v[0]), int(v[1]))).(color.Gray)
				fmt.Println(grayColor.Y, "灰度")
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
				fmt.Println(grayColor.Y, "灰度")
				if grayColor.Y > 150 {
					log.Error("路径：[%v] 校验未通过", route.RoutesName)
					return fmt.Errorf("路径：[%v] 校验未通过", route.RoutesName)
				}
			}
		}
	}
	return nil
}
