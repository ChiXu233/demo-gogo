package service

import (
	"demo-gogo/api/apimodel"
	"demo-gogo/database"
)

var resourceOperator Operator

type ResourceOperator struct {
	database.Database
}

type Operator interface {
	CreateOrUpdateMap(req *apimodel.MapRequest) error
	ListMap(req *apimodel.MapRequest) (*apimodel.MapPageResponse, error)
	DeleteMap(req *apimodel.MapRequest) error
	CreateOrUpdateMapInfo(req *apimodel.MapInfoRequest) error
	ListMapInfoPageResponse(req *apimodel.MapInfoRequest) (*apimodel.MapInfoPageResponse, error)
	DeleteMapInfo(req *apimodel.MapInfoRequest) error
	CreateOrUpdateNode(req *apimodel.RouteNodesRequest) error
	ListMapNodes(req *apimodel.RouteNodesRequest) (*apimodel.RouteNodesResponse, error)
	DeleteMapNodes(req *apimodel.RouteNodesRequest) error
	CreateOrUpdateMapRoute(req *apimodel.MapRoutesArrRequest) error
	ListMapRoutes(req *apimodel.MapRoutesRequest) (*apimodel.MapRoutesResponse, error)
	DeleteMapRoute(req *apimodel.MapRoutesRequest) error
	CheckRoute(req *apimodel.MapRoutesArrRequest) error
	ListMapInfo(req *apimodel.RouteNodesRequest) (*apimodel.MapInfosResponse, error)
	BatchDeleteMapNodes(req *apimodel.BatchDeleteNodes) error
}

func GetOperator() Operator {
	if resourceOperator == nil {
		resourceOperator = &ResourceOperator{
			Database: database.GetDatabase(),
		}
	}
	return resourceOperator
}

func NewMockOperator() ResourceOperator {
	return ResourceOperator{
		Database: database.GetDatabase(),
	}
}

func (operator *ResourceOperator) TransactionBegin() (*ResourceOperator, error) {
	db, err := database.GetDatabase().Begin()
	if err != nil {
		return nil, err
	}
	return &ResourceOperator{
		Database: db,
	}, nil
}

func (operator *ResourceOperator) TransactionCommit() error {
	return operator.Database.Commit()
}

func (operator *ResourceOperator) TransactionRollback() error {
	return operator.Database.Rollback()
}
