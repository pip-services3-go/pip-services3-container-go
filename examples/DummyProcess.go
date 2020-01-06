package examples

import "github.com/pip-services3-go/pip-services3-container-go/v3/container"

func NewDummyProcess() *container.ProcessContainer {
	c := container.NewProcessContainer("dummy", "Sample dummy process")
	c.SetConfigPath("./examples/dummy.yaml")
	c.AddFactory(NewDummyFactory())
	return c
}