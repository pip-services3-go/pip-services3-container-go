package examples

import (
	"github.com/pip-services3-go/pip-services3-commons-go/v3/refer"
	"github.com/pip-services3-go/pip-services3-components-go/v3/build"
)

var ControllerDescriptor = refer.NewDescriptor("pip-services-dummies", "controller", "default", "*", "1.0")

func NewDummyFactory() *build.Factory {
	factory := build.NewFactory()

	factory.RegisterType(ControllerDescriptor, NewDummyController)

	return factory
}
