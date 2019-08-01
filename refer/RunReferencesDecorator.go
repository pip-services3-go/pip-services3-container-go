package refer

import (
	crefer "github.com/pip-services3-go/pip-services3-commons-go/refer"
	"github.com/pip-services3-go/pip-services3-commons-go/run"
)

type RunReferencesDecorator struct {
	ReferencesDecorator
	opened bool
}

func NewRunReferencesDecorator(nextReferences crefer.IReferences,
	topReferences crefer.IReferences) *RunReferencesDecorator {
	return &RunReferencesDecorator{
		ReferencesDecorator: *NewReferencesDecorator(nextReferences, topReferences),
	}
}

func (c *RunReferencesDecorator) IsOpen() bool {
	return c.opened
}

func (c *RunReferencesDecorator) Open(correlationId string) error {
	if !c.opened {
		components := c.GetAll()
		err := run.Opener.Open(correlationId, components)
		c.opened = err == nil
		return err
	}
	return nil
}

func (c *RunReferencesDecorator) Close(correlationId string) error {
	if c.opened {
		components := c.GetAll()
		err := run.Closer.Close(correlationId, components)
		c.opened = false
		return err
	}
	return nil
}

func (c *RunReferencesDecorator) Put(locator interface{}, component interface{}) {
	c.ReferencesDecorator.Put(locator, component)

	if c.opened {
		run.Opener.OpenOne("", component)
	}
}

func (c *RunReferencesDecorator) Remove(locator interface{}) interface{} {
	component := c.ReferencesDecorator.Remove(locator)

	if c.opened {
		run.Closer.CloseOne("", component)
	}

	return component
}

func (c *RunReferencesDecorator) RemoveAll(locator interface{}) []interface{} {
	components := c.NextReferences.RemoveAll(locator)

	if c.opened {
		run.Closer.Close("", components)
	}

	return components
}
