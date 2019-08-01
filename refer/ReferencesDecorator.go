package refer

import crefer "github.com/pip-services3-go/pip-services3-commons-go/refer"

type ReferencesDecorator struct {
	NextReferences crefer.IReferences
	TopReferences  crefer.IReferences
}

func NewReferencesDecorator(nextReferences crefer.IReferences,
	topReferences crefer.IReferences) *ReferencesDecorator {
	c := &ReferencesDecorator{
		NextReferences: nextReferences,
		TopReferences:  topReferences,
	}

	if c.NextReferences == nil {
		c.NextReferences = topReferences
	}
	if c.TopReferences == nil {
		c.TopReferences = nextReferences
	}

	return c
}

func (c *ReferencesDecorator) Put(locator interface{}, component interface{}) {
	c.NextReferences.Put(locator, component)
}

func (c *ReferencesDecorator) Remove(locator interface{}) interface{} {
	return c.NextReferences.Remove(locator)
}

func (c *ReferencesDecorator) RemoveAll(locator interface{}) []interface{} {
	return c.NextReferences.RemoveAll(locator)
}

func (c *ReferencesDecorator) GetAllLocators() []interface{} {
	return c.NextReferences.GetAllLocators()
}

func (c *ReferencesDecorator) GetAll() []interface{} {
	return c.NextReferences.GetAll()
}

func (c *ReferencesDecorator) GetOneOptional(locator interface{}) interface{} {
	var component interface{}

	defer func() {
		recover()
	}()

	components, err := c.Find(locator, false)
	if err == nil && len(components) > 0 {
		component = components[0]
	}

	return component
}

func (c *ReferencesDecorator) GetOneRequired(locator interface{}) (interface{}, error) {
	components, err := c.Find(locator, true)
	if err != nil || len(components) == 0 {
		return nil, err
	}
	return components[0], nil
}

func (c *ReferencesDecorator) GetOptional(locator interface{}) []interface{} {
	components := []interface{}{}

	defer func() {
		recover()
	}()

	components, _ = c.Find(locator, false)

	return components
}

func (c *ReferencesDecorator) GetRequired(locator interface{}) ([]interface{}, error) {
	return c.Find(locator, true)
}

func (c *ReferencesDecorator) Find(locator interface{}, required bool) ([]interface{}, error) {
	return c.NextReferences.Find(locator, required)
}
