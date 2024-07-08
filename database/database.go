package database

import (
	"demo-gogo/config"
	"demo-gogo/database/model"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/wonderivan/logger"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var ormDB Database

type OrmDB struct {
	*gorm.DB
}

type Database interface {
	GetEntityByID(table string, id int, entity interface{}) error
	GetEntityForUpdate(table string, id int, entity interface{}) error
	AssertRowExist(table string, filter map[string]interface{}, params model.QueryParams, entity interface{}) error
	ListEntityByFilter(table string, filter map[string]interface{}, params model.QueryParams, entities interface{}) error
	CountEntityByFilter(table string, filter map[string]interface{}, params model.QueryParams, count *int64) error
	CountAllEntityByFilter(table string, filter map[string]interface{}, params model.QueryParams, count *int64) error
	GetEntityPluck(table string, filter map[string]interface{}, params model.QueryParams, column string, cols interface{}) error
	CreateEntity(table string, entity interface{}) error
	BatchCreateEntity(table string, entities interface{}) error
	SaveEntity(table string, updater interface{}) error
	UpdateEntityByFilter(table string, filter map[string]interface{}, params model.QueryParams, updater interface{}) error
	DeleteEntityByFilter(table string, filter map[string]interface{}, params model.QueryParams, mode interface{}) error

	DeleteEntity(mode interface{}) error

	DeleteUnscopedEntityByFilter(table string, filter map[string]interface{}, params model.QueryParams, mode interface{}) error

	Begin() (Database, error)
	Commit() error
	Rollback() error
}

func InitDB() error {
	DBConfig := config.Conf.DB
	sqlConnection := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s TimeZone=Asia/Shanghai",
		DBConfig.Host, DBConfig.Port, DBConfig.User, DBConfig.Name, DBConfig.Password)
	db, err := gorm.Open(postgres.Open(sqlConnection), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		IgnoreRelationshipsWhenMigrating:         true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}
	// sqlDB.SetMaxIdleConns(DBConfig.MaxIdleConnects)
	// sqlDB.SetMaxOpenConns(DBConfig.MaxOpenConnects)
	if config.Conf.DB.InitTable {
		initTable(db)
	}

	switch config.Conf.APP.Mode {
	case gin.ReleaseMode:
		db.Logger = db.Config.Logger.LogMode(logger.Error)
	case gin.TestMode:
		db.Logger = db.Config.Logger.LogMode(logger.Warn)
	case gin.DebugMode:
		db.Logger = db.Config.Logger.LogMode(logger.Info)
	}
	ormDB = &OrmDB{
		DB: db,
	}
	return nil
}

func GetDatabase() Database {
	return ormDB
}

// SetMockDatabase for unit test
func SetMockDatabase(mockDB Database) {
	ormDB = mockDB
}

func initTable(db *gorm.DB) {
	err := db.AutoMigrate(&model.Map{})
	if err != nil {
		log.Error("init table[%s] error.[%s]", model.TableNameMap, err.Error())
	}
	err = db.AutoMigrate(&model.MapInfo{})
	if err != nil {
		log.Error("init table[%s] error.[%s]", model.TableNameMap, err.Error())
	}
	err = db.AutoMigrate(&model.MapRoutes{})
	if err != nil {
		log.Error("init table[%s] error.[%s]", model.TableNameMapRoutes, err.Error())
	}
	err = db.AutoMigrate(&model.MapRouteNodes{})
	if err != nil {
		log.Error("init table[%s] error.[%s]", model.TableNameMapRouteNodes, err.Error())
	}
}

func (db *OrmDB) Begin() (Database, error) {
	tx := db.DB.Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}
	return &OrmDB{DB: tx}, nil
}

func (db *OrmDB) Commit() error {
	tx := db.DB.Commit()
	if err := tx.Error; err != nil {
		return err
	}
	return nil
}

func (ormdb *OrmDB) Rollback() error {
	tx := ormdb.DB.Rollback()
	if err := tx.Error; err != nil {
		return err
	}
	return nil
}

// NewSqliteDatabase for unit test
func NewSqliteDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Error("failed to connect sqlite database")
		return nil, err
	}
	initTable(db)
	return db, nil
}

// NewPostgresDatabase for unit test
func NewPostgresDatabase(host, user, password, dbName string, port int) (*gorm.DB, error) {
	sqlConnection := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s TimeZone=Asia/Shanghai",
		host, port, user, dbName, password)
	db, err := gorm.Open(postgres.Open(sqlConnection), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		IgnoreRelationshipsWhenMigrating:         true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Error("failed to connect postgres database")
		return nil, err
	}
	return db, nil
}
