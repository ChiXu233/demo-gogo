package handler

import "demo-gogo/service"

type RestHandler struct {
	Operator service.Operator
}

func NewHandler() *RestHandler {
	return &RestHandler{
		Operator: service.GetOperator(),
	}
}
