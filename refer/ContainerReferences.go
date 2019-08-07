package refer

import (
	"fmt"

	cconfig "github.com/pip-services3-go/pip-services3-commons-go/config"
	"github.com/pip-services3-go/pip-services3-commons-go/refer"
	"github.com/pip-services3-go/pip-services3-commons-go/reflect"
	"github.com/pip-services3-go/pip-services3-components-go/build"
	"github.com/pip-services3-go/pip-services3-container-go/config"
)

type ContainerReferences struct {
	ManagedReferences
}

func NewContainerReferences() *ContainerReferences {
	return &ContainerReferences {
		ManagedReferences: *NewEmptyManagedReferences(),
	}
}

func (c *ContainerReferences) PutFromConfig(config config.ContainerConfig) error {
	var err error
	var locator interface{}
	var component interface{}

	defer func() {
		if r := recover(); r != nil {
			err, _ = r.(error)
		}
	}()

	for _, componentConfig := range config {
		if componentConfig.Type != nil {
			// Create component dynamically
			locator = componentConfig.Type
			component, err = reflect.TypeReflector.CreateInstanceByDescriptor(componentConfig.Type)
		} else if componentConfig.Descriptor != nil {
			// Or create component statically
			locator = componentConfig.Descriptor
			factory := c.ManagedReferences.Builder.FindFactory(locator)
			component = c.ManagedReferences.Builder.Create(locator, factory)
			if component == nil {
				return refer.NewReferenceError("", locator)
			}
			locator = c.ManagedReferences.Builder.ClarifyLocator(locator, factory)
		}

		// Check that component was created
		if component == nil {
			return build.NewCreateError(
				"CANNOT_CREATE_COMPONENT", 
				"Cannot create component",
			).WithDetails("config", config)
		}
		
		fmt.Printf("Created component %v\n", locator)

		// Add component to the list
		c.ManagedReferences.References.Put(locator, component)

		// Configure component
		configurable, ok := component.(cconfig.IConfigurable)
		if ok {
			configurable.Configure(componentConfig.Config)
		}

		// Set references to factories
		_, ok = component.(build.IFactory)
		if ok {
			referenceable, ok := component.(refer.IReferenceable)
			if ok {
				referenceable.SetReferences(c)
			}
		}
	}

	return err
}