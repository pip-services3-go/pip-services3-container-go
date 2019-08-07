package config

import (
	"path/filepath"

	"github.com/pip-services3-go/pip-services3-commons-go/config"
	"github.com/pip-services3-go/pip-services3-commons-go/errors"
	cconfig "github.com/pip-services3-go/pip-services3-components-go/config"
)

type TContainerConfigReader struct {}

var ContainerConfigReader = &TContainerConfigReader{}

func (c *TContainerConfigReader) ReadFromFile(correlationId string,
    path string, parameters *config.ConfigParams) (ContainerConfig, error) {
	if path == "" {
		return nil, errors.NewConfigError(correlationId, "NO_PATH", "Missing config file path")
	}

	ext := filepath.Ext(path)

	if ext == "json" {
		return c.ReadFromJsonFile(correlationId, path, parameters)
	}

	if ext == "yaml" || ext == "yml" {
		return c.ReadFromYamlFile(correlationId, path, parameters)
	}

	return c.ReadFromJsonFile(correlationId, path, parameters)
}

func (c *TContainerConfigReader) ReadFromJsonFile(correlationId string,
    path string, parameters *config.ConfigParams) (ContainerConfig, error) {
	config, err := cconfig.ReadJsonConfig(correlationId, path, parameters)
	if err != nil {
		return nil, err
	}
	return ReadContainerConfigFromConfig(config)
}

func (c *TContainerConfigReader) ReadFromYamlFile(correlationId string,
    path string, parameters *config.ConfigParams) (ContainerConfig, error) {
	config, err := cconfig.ReadYamlConfig(correlationId, path, parameters)
	if err != nil {
		return nil, err
	}
	return ReadContainerConfigFromConfig(config)
}