package test_refer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pip-services3-go/pip-services3-commons-go/v3/refer"
	"github.com/pip-services3-go/pip-services3-components-go/v3/log"
	crefer "github.com/pip-services3-go/pip-services3-container-go/v3/refer"
)

func TestAutoCreateComponent(t *testing.T) {
	refs := crefer.NewEmptyManagedReferences()

	factory := log.NewDefaultLoggerFactory()
	refs.Put(nil, factory)

	logger, err := refs.GetOneRequired(
		refer.NewDescriptor("*", "logger", "*", "*", "*"),
	)

	assert.Nil(t, err)
	assert.NotNil(t, logger)
}

func TestStringLocator(t *testing.T) {
	refs := crefer.NewEmptyManagedReferences()

	factory := log.NewDefaultLoggerFactory()
	refs.Put(nil, factory)

	logger := refs.GetOneOptional("ABC")

	assert.Nil(t, logger)
}

func TestNilLocator(t *testing.T) {
	refs := crefer.NewEmptyManagedReferences()

	factory := log.NewDefaultLoggerFactory()
	refs.Put(nil, factory)

	logger := refs.GetOneOptional(nil)

	assert.Nil(t, logger)
}
