package refer

import (
	crefer "github.com/pip-services3-go/pip-services3-commons-go/refer"
)

type ManagedReferences struct {
	ReferencesDecorator
	References *crefer.References
	Builder    *BuildReferencesDecorator
	Linker     *LinkReferencesDecorator
	Runner     *RunReferencesDecorator
}

func NewManagedReferences(tuples []interface{}) *ManagedReferences {
	c := &ManagedReferences{
		ReferencesDecorator: *NewReferencesDecorator(nil, nil),
	}

	c.References = crefer.NewReferences(tuples)
	c.Builder = NewBuildReferencesDecorator(c.References, c)
	c.Linker = NewLinkReferencesDecorator(c.Builder, c)
	c.Runner = NewRunReferencesDecorator(c.Linker, c)

	c.ReferencesDecorator.NextReferences = c.Runner

	return c
}

func NewEmptyManagedReferences() *ManagedReferences {
	return NewManagedReferences([]interface{}{})
}

func NewManagedReferencesFromTuples(tuples ...interface{}) *ManagedReferences {
	return NewManagedReferences(tuples)
}

func (c *ManagedReferences) IsOpen() bool {
	return c.Linker.IsOpen() && c.Runner.IsOpen()
}

func (c *ManagedReferences) Open(correlationId string) error {
	err := c.Linker.Open(correlationId)
	if err == nil {
		err = c.Runner.Open(correlationId)
	}
	return err
}

func (c *ManagedReferences) Close(correlationId string) error {
	err := c.Runner.Close(correlationId)
	if err == nil {
		err = c.Linker.Close(correlationId)
	}
	return err
}
