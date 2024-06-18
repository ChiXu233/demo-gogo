package database

import (
	"demo-gogo/database/model"
	"demo-gogo/utils"
	"errors"
	"fmt"
	log "github.com/wonderivan/logger"
	"gorm.io/gorm"
	"reflect"
)

func ProcessQueryParams(db *gorm.DB, queryParams model.QueryParams) *gorm.DB {
	// 模糊查询和OR查询
	for _, subquery := range queryParams.SubQueries {
		if subquery != nil && len(subquery.Keywords) > 0 {
			db = db.Where(subquery.Statement(), subquery.Parameters()...)
		}

		if subquery != nil && len(subquery.OrConditions) > 0 {
			db = db.Where(subquery.OrConstruct())
		}
	}

	// 比较查询  多样比较 等于 不等于 大于 小于
	for _, compareQuery := range queryParams.CompareQueries {
		if compareQuery != nil && compareQuery.Field != "" && len(compareQuery.Parameters()) > 0 {
			db = db.Where(compareQuery.Statement(), compareQuery.Parameters()[0])
		}
	}

	// in查询
	for _, in := range queryParams.InQueries {
		if in != nil && in.Statement() != "" {
			db = db.Where(in.Statement(), in.Parameters())
		}
	}

	// not in 查询
	for _, notIn := range queryParams.NotInQueries {
		if notIn != nil {
			db = db.Not(notIn)
		}
	}

	// between查询
	for _, rangeQuery := range queryParams.RangeQueries {
		if rangeQuery != nil && rangeQuery.Statement() != "" {
			db = db.Where(rangeQuery.Statement(), rangeQuery.Parameters()...)
		}
	}

	// 排序
	for _, order := range queryParams.Orders {
		db = db.Order(fmt.Sprintf("%s %s", order.Field, order.Direction))
	}

	// 偏移
	if queryParams.Offset != nil {
		db = db.Offset(*queryParams.Offset)
	}

	// 分页
	if queryParams.Limit != nil {
		db = db.Limit(*queryParams.Limit)
	}
	// 去重! 去重后只有单列返回
	for _, field := range queryParams.DistinctQueries {
		db = db.Distinct(field)
	}

	for _, field := range queryParams.GroupFields {
		db = db.Group(field)
	}
	return db
}

// GetEntityByID 根据ID查询实体，entity必须是指针
func (db *OrmDB) GetEntityByID(table string, id int, entity interface{}) error {
	if reflect.ValueOf(entity).Kind() != reflect.Ptr {
		return errors.New("GetEntityByID [entity] Kind Must Ptr")
	}
	defer utils.TimeCost()(fmt.Sprintf("[%s]GetEntityByID_timeCost", table))
	if err := db.Table(table).Where("id = ?", id).First(entity).Error; err != nil {
		log.Error("[%s]GetEntityByID Error.id[%#v] err[%#v]", table, id, err)
		return err
	}
	return nil
}

// GetEntityForUpdate 根据ID查询实体，entity必须是指针,锁定数据行
func (db *OrmDB) GetEntityForUpdate(table string, id int, entity interface{}) error {
	if reflect.ValueOf(entity).Kind() != reflect.Ptr {
		return errors.New("GetEntityByID [entity] Kind Must Ptr")
	}
	defer utils.TimeCost()(fmt.Sprintf("[%s]GetEntityForUpdate_timeCost", table))
	if err := db.Set("gorm:query_option", "FOR UPDATE").Table(table).Where("id = ?", id).First(entity).Error; err != nil {
		log.Error("[%s]GetEntityByID Error.id[%#v] err[%#v]", table, id, err)
		return err
	}
	return nil
}

// AssertRowExist 断言记录存在，如果不存在则报错，存在直接返回
func (db *OrmDB) AssertRowExist(table string, filter map[string]interface{}, params model.QueryParams, entity interface{}) error {
	if reflect.ValueOf(entity).Kind() != reflect.Ptr {
		return errors.New("AssertRowExist [entity] Kind Must Ptr")
	}
	defer utils.TimeCost()(fmt.Sprintf("[%s]AssertRowExist_timeCost", table))
	tx := db.Table(table).Where(filter)
	if err := ProcessQueryParams(tx, params).First(entity).Error; err != nil {
		log.Error("[%s]AssertRowExist Error.filter[%#v] params[%#v] err[%#v]", table, filter, params, err)
		return err
	}
	return nil
}

// ListEntityByFilter 多条件查询实体 entities是一个实体对象切片对象
func (db *OrmDB) ListEntityByFilter(table string, filter map[string]interface{}, params model.QueryParams, entities interface{}) error {
	if reflect.ValueOf(entities).Kind() != reflect.Ptr {
		return errors.New("ListEntityByFilter [entities] Kind Must Ptr")
	}
	defer utils.TimeCost()(fmt.Sprintf("[%s]ListEntityByFilter_timeCost", table))
	tx := db.Table(table).Where(filter)
	if err := ProcessQueryParams(tx, params).Find(entities).Error; err != nil {
		log.Error("[%s]ListEntityByFilter Error.filter[%#v] params[%#v] err[%#v]", table, filter, params, err)
		return err
	}
	return nil
}

// GetEntityPluck 获取表的某一列数据 cols需要是切片的地址
func (db *OrmDB) GetEntityPluck(table string, filter map[string]interface{}, params model.QueryParams, column string, cols interface{}) error {
	if reflect.ValueOf(cols).Kind() != reflect.Ptr {
		return errors.New("GetEntityPluck [entities] Kind Must ptr")
	}
	defer utils.TimeCost()(fmt.Sprintf("[%s]GetEntityPluck_timeCost", table))
	tx := db.Table(table).Where(filter).Where(model.FieldDeletedTime, nil)
	if err := ProcessQueryParams(tx, params).Pluck(column, cols).Error; err != nil {
		log.Error("[%s]GetFramePluck Error.filter[%#v] params[%#v] colum[%#v] err[%#v]", table, filter, params, column, err)
		return err
	}
	return nil
}

// CountEntityByFilter 分页查询计数 mode
func (db *OrmDB) CountEntityByFilter(table string, filter map[string]interface{}, params model.QueryParams, count *int64) error {
	defer utils.TimeCost()(fmt.Sprintf("[%s]CountEntityByFilter_timeCost", table))
	tx := db.Table(table).Where(filter).Where(model.FieldDeletedTime, nil)
	if err := ProcessQueryParams(tx, params).Count(count).Error; err != nil {
		log.Error("[%s]CountEntityByFilter Error.filter[%#v] params[%#v] err[%#v]", table, filter, params, err)
		return err
	}
	return nil
}

// CountAllEntityByFilter 分页查询计数 mode
func (db *OrmDB) CountAllEntityByFilter(table string, filter map[string]interface{}, params model.QueryParams, count *int64) error {
	defer utils.TimeCost()(fmt.Sprintf("[%s]CountAllEntityByFilter_timeCost", table))
	tx := db.Table(table).Where(filter)
	if err := ProcessQueryParams(tx, params).Count(count).Error; err != nil {
		log.Error("[%s]CountAllEntityByFilter Error.filter[%#v] params[%#v] err[%#v]", table, filter, params, err)
		return err
	}
	return nil
}

// CreateEntity 创建实体对象，并将创建的ID填充到entity
func (db *OrmDB) CreateEntity(table string, entity interface{}) error {
	if reflect.ValueOf(entity).Kind() != reflect.Ptr {
		return errors.New("CreateEntity [entity] Kind Must Ptr")
	}
	defer utils.TimeCost()(fmt.Sprintf("[%s]CreateEntity_timeCost", table))
	if err := db.Table(table).Create(entity).Error; err != nil {
		log.Error("[%s]CreateEntity Error.entity[%#v] err[%#v]", table, entity, err)
		return err
	}
	return nil
}

// BatchCreateEntity 批量创建实体对象，并将创建的ID填充到entities entities需要是切片
// entities slice本身就是引用传递，不涉及slice扩容，可以不需要用指针
func (db *OrmDB) BatchCreateEntity(table string, entities interface{}) error {
	if reflect.ValueOf(entities).Kind() != reflect.Slice {
		return errors.New("BatchCreateEntity [entity] Kind Must slice")
	}
	defer utils.TimeCost()(fmt.Sprintf("[%s]BatchCreateEntity_timeCost", table))
	batch := 1
	insertLen := reflect.ValueOf(entities).Len()
	if insertLen > 100 {
		batch = insertLen/100 + 1
		if insertLen%100 > 0 {
			batch += 1
		}
	}
	if err := db.Table(table).CreateInBatches(entities, batch).Error; err != nil {
		log.Error("[%s]BatchCreateEntity Error.entities[%#v] err[%#v]", table, entities, err)
		return err
	}
	return nil
}

// SaveEntity 如果包含主键id，执行更新。更新实体对象，updater需要是对象，全字段更新（默认值也会），更新结果为0会执行插入。如果没有id，执行插入
func (db *OrmDB) SaveEntity(table string, updater interface{}) error {
	if reflect.ValueOf(updater).Kind() != reflect.Ptr && reflect.ValueOf(updater).Kind() != reflect.Slice {
		return errors.New("SaveEntity [updater] Kind Must Ptr Or Slice")
	}
	defer utils.TimeCost()(fmt.Sprintf("[%s]SaveEntity_timeCost", table))
	tx := db.Table(table).Save(updater)
	if err := tx.Error; err != nil {
		log.Error("[%s]SaveEntity Error.updater[%#v] err[%#v]", table, updater, err)
		return err
	}
	if tx.RowsAffected == 0 {
		log.Warn("[%s]SaveEntity Warn.updater[%#v]", table, updater)
	}
	return nil
}

// UpdateEntityByFilter 根据条件更新实体对象，updater是结构体指针，默认值字段不更新。为map指针，字段会更新（map不会触发逻辑更新操作！）
func (db *OrmDB) UpdateEntityByFilter(table string, filter map[string]interface{}, params model.QueryParams, updater interface{}) error {
	if reflect.ValueOf(updater).Kind() != reflect.Ptr {
		return errors.New("UpdateEntityByFilter [updater] Kind Must Ptr")
	}
	defer utils.TimeCost()(fmt.Sprintf("[%s]UpdateEntityByFilter_timeCost", table))
	tx := db.Table(table).Where(filter)
	if err := ProcessQueryParams(tx, params).Updates(updater).Error; err != nil {
		log.Error("[%s]UpdateEntityByFilter Error.filter[%#v] params[%#v] updater[%#v] err[%#v]", table, filter, params, updater, err)
		return err
	}
	if tx.RowsAffected == 0 {
		log.Warn("[%s]UpdateEntityByFilter Warn.filter[%#v] params[%#v] updater[%#v]", table, filter, params, updater)
	}
	return nil
}

// DeleteEntityByFilter 逻辑删除 mode应该是空结构体指针，用来触发逻辑删除
func (db *OrmDB) DeleteEntityByFilter(table string, filter map[string]interface{}, params model.QueryParams, mode interface{}) error {
	if reflect.ValueOf(mode).Kind() != reflect.Ptr {
		return errors.New("DeleteEntityByFilter [mode] Kind Must Ptr")
	}
	defer utils.TimeCost()(fmt.Sprintf("[%s]DeleteEntityByFilter_timeCost", table))
	tx := db.Table(table).Where(filter)
	if err := ProcessQueryParams(tx, params).Delete(mode).Error; err != nil {
		log.Error("[%s]DeleteEntityByFilter Error.filter[%#v] params[%#v] err[%#v]", table, filter, params, err)
		return err
	}
	if tx.RowsAffected == 0 {
		log.Warn("[%s]DeleteEntityByFilter Warn.filter[%#v] params[%#v]", table, filter, params)
	}
	return nil
}

// DeleteEntity 逻辑删除 mode应该是空结构体指针
func (db *OrmDB) DeleteEntity(mode interface{}) error {
	if reflect.ValueOf(mode).Kind() != reflect.Ptr {
		return errors.New("DeleteEntity [mode] Kind Must Ptr")
	}
	defer utils.TimeCost()(fmt.Sprintf("[%v]DeleteEntity_timeCost", mode))
	tx := db.Delete(mode)
	if err := tx.Error; err != nil {
		log.Error("[%#v]DeleteEntity Error.err[%#v]", mode, err)
		return err
	}
	if tx.RowsAffected == 0 {
		log.Warn("[%#v]DeleteEntity Warn", mode)
	}
	return nil
}

// DeleteUnscopedEntityByFilter 硬删除 mode应该是空结构体指针，用来触发逻辑删除
func (db *OrmDB) DeleteUnscopedEntityByFilter(table string, filter map[string]interface{}, params model.QueryParams, mode interface{}) error {
	if reflect.ValueOf(mode).Kind() != reflect.Ptr {
		return errors.New("DeleteUnscopedEntityByFilter [mode] Kind Must Ptr")
	}
	defer utils.TimeCost()(fmt.Sprintf("[%s]DeleteUnscopedEntityByFilter_timeCost", table))
	tx := db.Unscoped().Table(table).Where(filter)
	if err := ProcessQueryParams(tx, params).Delete(mode).Error; err != nil {
		log.Error("[%s]DeleteUnscopedEntityByFilter Error.filter[%#v] params[%#v] err[%#v]", table, filter, params, err)
		return err
	}
	if tx.RowsAffected == 0 {
		log.Warn("[%s]DeleteUnscopedEntityByFilter Warn.filter[%#v] params[%#v]", table, filter, params)
	}
	return nil
}
