package refer

import (
	crefer "github.com/pip-services3-go/pip-services3-commons-go/refer"
)

type LinkReferencesDecorator struct {
	ReferencesDecorator
	opened bool
}

func NewLinkReferencesDecorator(nextReferences crefer.IReferences,
	topReferences crefer.IReferences) *LinkReferencesDecorator {
	return &LinkReferencesDecorator{
		ReferencesDecorator: *NewReferencesDecorator(nextReferences, topReferences),
	}
}

func (c *LinkReferencesDecorator) IsOpen() bool {
	return c.opened
}

func (c *LinkReferencesDecorator) Open(correlationId string) error {
	if !c.opened {
		c.opened = true
		components := c.GetAll()
		crefer.Referencer.SetReferences(c.ReferencesDecorator.TopReferences, components)
	}
	return nil
}

func (c *LinkReferencesDecorator) Close(correlationId string) error {
	if c.opened {
		c.opened = false
		components := c.GetAll()
		crefer.Referencer.UnsetReferences(components)
	}
	return nil
}

func (c *LinkReferencesDecorator) Put(locator interface{}, component interface{}) {
	c.ReferencesDecorator.Put(locator, component)

	if c.opened {
		crefer.Referencer.SetReferencesForOne(c.ReferencesDecorator.TopReferences, component)
	}
}

func (c *LinkReferencesDecorator) Remove(locator interface{}) interface{} {
	component := c.ReferencesDecorator.Remove(locator)

	if c.opened {
		crefer.Referencer.UnsetReferencesForOne(component)
	}

	return component
}

func (c *LinkReferencesDecorator) RemoveAll(locator interface{}) []interface{} {
	components := c.NextReferences.RemoveAll(locator)

	if c.opened {
		crefer.Referencer.UnsetReferences(components)
	}

	return components
}
