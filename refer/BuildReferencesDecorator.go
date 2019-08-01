package refer

import (
	"github.com/pip-services3-go/pip-services3-commons-go/refer"
	crefer "github.com/pip-services3-go/pip-services3-commons-go/refer"
	"github.com/pip-services3-go/pip-services3-components-go/build"
)

type BuildReferencesDecorator struct {
	ReferencesDecorator
}

func NewBuildReferencesDecorator(nextReferences crefer.IReferences,
	topReferences crefer.IReferences) *BuildReferencesDecorator {
	return &BuildReferencesDecorator{
		ReferencesDecorator: *NewReferencesDecorator(nextReferences, topReferences),
	}
}

func (c *BuildReferencesDecorator) FindFactory(locator interface{}) build.IFactory {
	components := c.GetAll()

	for _, component := range components {
		factory, ok := component.(build.IFactory)
		if ok && factory.CanCreate(locator) != nil {
			return factory
		}
	}

	return nil
}

func (c *BuildReferencesDecorator) Create(locator interface{},
	factory build.IFactory) interface{} {

	if factory == nil {
		return nil
	}

	var result interface{}

	defer func() {
		recover()
	}()

	result, _ = factory.Create(locator)

	return result
}

func (c *BuildReferencesDecorator) ClarifyLocator(locator interface{},
	factory build.IFactory) interface{} {

	if factory == nil {
		return nil
	}

	descriptor, ok := locator.(*refer.Descriptor)
	if !ok {
		return locator
	}

	anotherLocator := factory.CanCreate(locator)
	anotherDescriptor, ok1 := anotherLocator.(*refer.Descriptor)
	if !ok1 {
		return locator
	}

	group := descriptor.Group()
	if group == "" {
		group = anotherDescriptor.Group()
	}
	typ := descriptor.Type()
	if typ == "" {
		typ = anotherDescriptor.Type()
	}
	kind := descriptor.Kind()
	if kind == "" {
		kind = anotherDescriptor.Kind()
	}
	name := descriptor.Name()
	if name == "" {
		name = anotherDescriptor.Name()
	}
	version := descriptor.Version()
	if version == "" {
		version = anotherDescriptor.Version()
	}

	return refer.NewDescriptor(group, typ, kind, name, version)
}

func (c *BuildReferencesDecorator) GetOneOptional(locator interface{}) interface{} {
	components, err := c.Find(locator, false)
	if err != nil || len(components) == 0 {
		return nil
	}
	return components[0]
}

func (c *BuildReferencesDecorator) GetOneRequired(locator interface{}) (interface{}, error) {
	components, err := c.Find(locator, true)
	if err != nil || len(components) == 0 {
		return nil, err
	}
	return components[0], nil
}

func (c *BuildReferencesDecorator) GetOptional(locator interface{}) []interface{} {
	components, _ := c.Find(locator, false)
	return components
}

func (c *BuildReferencesDecorator) GetRequired(locator interface{}) ([]interface{}, error) {
	return c.Find(locator, true)
}

func (c *BuildReferencesDecorator) Find(locator interface{}, required bool) ([]interface{}, error) {
	components, _ := c.ReferencesDecorator.Find(locator, required)

	if required && len(components) == 0 {
		factory := c.FindFactory(locator)
		component := c.Create(locator, factory)
		if component != nil {
			locator = c.ClarifyLocator(locator, factory)
			c.ReferencesDecorator.TopReferences.Put(locator, component)
			components = append(components, component)
		}
	}

	if required && len(components) == 0 {
		err := refer.NewReferenceError("", locator)
		return nil, err
	}

	return components, nil
}
