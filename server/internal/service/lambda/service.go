package lambda

import (
	"context"
)

type (
	Service interface {
		Register(c context.Context)
	}
	ServiceConductor struct{}
)

var Conductor Service

func init() {
	if Conductor == nil {
		Conductor = &ServiceConductor{}
	}
}

func (sc *ServiceConductor) Register(c context.Context) {

}
