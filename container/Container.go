package container

import (
	cconfig "github.com/pip-services3-go/pip-services3-commons-go/config"
	"github.com/pip-services3-go/pip-services3-commons-go/errors"
	crefer "github.com/pip-services3-go/pip-services3-commons-go/refer"
	"github.com/pip-services3-go/pip-services3-container-go/build"
	"github.com/pip-services3-go/pip-services3-container-go/config"
	"github.com/pip-services3-go/pip-services3-container-go/refer"
	cbuild "github.com/pip-services3-go/pip-services3-components-go/build"
	"github.com/pip-services3-go/pip-services3-components-go/info"
	"github.com/pip-services3-go/pip-services3-components-go/log"
)

type Container struct {
	logger log.ILogger
	factories *cbuild.CompositeFactory
	info *info.ContextInfo
	config config.ContainerConfig
	references *refer.ContainerReferences
	referenceable crefer.IReferenceable
	unreferenceable crefer.IUnreferenceable
}

func NewEmptyContainer() *Container {
	return &Container {
		logger: log.NewNullLogger(),
		factories: build.NewDefaultContainerFactory(),
		info: info.NewContextInfo(),
	}
}

func NewContainer(name string, description string) *Container {
	c := NewEmptyContainer()

	c.info.Name = name
	c.info.Description = description

	return c
}

func InheritContainer(name string, description string,
	referenceable crefer.IReferenceable) *Container {
	c := NewEmptyContainer()

	c.info.Name = name
	c.info.Description = description
	c.referenceable = referenceable
	c.unreferenceable, _ = referenceable.(crefer.IUnreferenceable)

	return c
}

func (c *Container) Configure(conf *cconfig.ConfigParams) {
	c.config, _ = config.ReadContainerConfigFromConfig(conf)
}

func (c *Container) ReadConfigFromFile(correlationId string,
	path string, parameters *cconfig.ConfigParams) error {

	var err error
	c.config, err = config.ContainerConfigReader.ReadFromFile(correlationId, path, parameters)
	//c.logger.Trace(correlationId, config.String())
	return err
}

func (c *Container) initReferences(references crefer.IReferences) {
	existingInfo, ok := references.GetOneOptional(
		crefer.NewDescriptor("pip-services", "context-info", "*", "*", "1.0"),
	).(*info.ContextInfo)
	if !ok {
		references.Put(
			crefer.NewDescriptor("pip-services", "context-info", "default", "default", "1.0"),
			c.info,
		)
	} else {
		c.info = existingInfo
	}

	references.Put(
		crefer.NewDescriptor("pip-services", "factory", "container", "default", "1.0"),
		c.factories,
	)
}

func (c *Container) Logger() log.ILogger {
	return c.logger
}

func (c *Container) SetLogger(logger log.ILogger) {
	c.logger = logger
}

func (c *Container) Info() *info.ContextInfo {
	return c.info
}

func (c *Container) AddFactory(factory cbuild.IFactory) {
	c.factories.Add(factory)
}

func (c *Container) IsOpen() bool {
	return c.references != nil
}

func (c *Container) Open(correlationId string) error {
	var err error

	if c.references != nil {
		return errors.NewInvalidStateError(
			correlationId, "ALREADY_OPENED", "Container was already opened",
		)
	}

	defer func() {
		if r := recover(); r != nil {
			err, _ = r.(error)
			c.logger.Error(correlationId, err, "Failed to start container")
			c.Close(correlationId)
		}
	}()

	c.logger.Trace(correlationId, "Starting container.");

	// Create references with configured components
	c.references = refer.NewContainerReferences()
	c.initReferences(c.references)
	c.references.PutFromConfig(c.config)

	if c.referenceable != nil {
		c.referenceable.SetReferences(c.references);
	}

	// Get custom description if available
	infoDescriptor := crefer.NewDescriptor("*", "context-info", "*", "*", "*")
	info, ok := c.references.GetOneOptional(infoDescriptor).(*info.ContextInfo)
	if ok {
		c.info = info
	}

	// Get reference to logger
	c.logger = log.NewCompositeLoggerFromReferences(c.references)

	// Open references
	err = c.references.Open(correlationId)
	if err == nil {
		c.logger.Info(correlationId, "Container %s started", c.info.Name)
	} else {
		c.logger.Fatal(correlationId, err, "Failed to start container")
		c.Close(correlationId)
	}

	return err
}

func (c *Container) Close(correlationId string) error {
	// Skip if container wasn't opened
	if c.references == nil {
		return nil
	}

	var err error

	defer func() {
		if r := recover(); r != nil {
			err, _ = r.(error)
			c.logger.Error(correlationId, err, "Failed to stop container")
		}
	}()

	c.logger.Trace(correlationId, "Stopping %s container", c.info.Name);

	// Unset references for child container
	if c.unreferenceable != nil {
		c.unreferenceable.UnsetReferences()
	}

	// Close and dereference components
	err = c.references.Close(correlationId);

	c.references = nil

	if err == nil {
		c.logger.Info(correlationId, "Container %s stopped", c.info.Name)
	} else {
		c.logger.Error(correlationId, err, "Failed to stop container")
	}

	return err
}