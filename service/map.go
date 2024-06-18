package service

import (
	"demo-gogo/api/apimodel"
	"demo-gogo/database/model"
	"demo-gogo/httpserver/errcode"
	"demo-gogo/utils"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	log "github.com/wonderivan/logger"
	"gorm.io/gorm"
	"image/png"
	"os"
)

func (operator *ResourceOperator) CreateOrUpdateMap(req *apimodel.MapRequest) error {
	var opt model.BaseMap
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
	if !utils.Exists(req.MapURL) {
		err = fmt.Errorf("路径文件不存在")
		return err
	}
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
	var maps []model.BaseMap
	err := operator.Database.CountEntityByFilter(model.TableNameMap, selector, model.OneQuery, &count)
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
		err = operator.Database.ListEntityByFilter(model.TableNameMap, selector, queryParams, &maps)
		if err != nil {
			log.Error("地图数据查询失败,err:[%v]", err)
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
