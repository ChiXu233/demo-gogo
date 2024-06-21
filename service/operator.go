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
	CreateOrUpdateBaseMap(req *apimodel.BaseMapRequest) error
	ListBaseMap(req *apimodel.BaseMapRequest) (*apimodel.MapPageResponse, error)
	DeleteBaseMap(req *apimodel.BaseMapRequest) error
	CreateOrUpdateNode(req *apimodel.RouteNodesRequest) error
	ListMapNodes(req *apimodel.RouteNodesRequest) (*apimodel.RouteNodesResponse, error)
	DeleteMapNodes(req *apimodel.RouteNodesRequest) error
	CreateOrUpdateMapRoute(req *apimodel.MapRoutesArrRequest) error
	ListMapRoutes(req *apimodel.MapRoutesRequest) (*apimodel.MapRoutesResponse, error)
	DeleteMapRoute(req *apimodel.MapRoutesRequest) error
	CreateOrUpdatePath(req *apimodel.PathRequest) error
	ListPath(req *apimodel.PathRequest) (*apimodel.PathResponse, error)
	DeletePath(req *apimodel.PathRequest) error
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
