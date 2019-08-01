package config

import (
	"github.com/pip-services3-go/pip-services3-commons-go/config"
)

type ContainerConfig []*ComponentConfig

func NewContainerConfig(components ...*ComponentConfig) ContainerConfig {
	return components
}

func NewContainerConfigFromValue(value interface{}) ContainerConfig {
	config := config.NewConfigParamsFromValue(value)
	result, _ := ReadContainerConfigFromConfig(config)
	return result
}

func ReadContainerConfigFromConfig(config *config.ConfigParams) (ContainerConfig, error) {
	if config == nil {
		return []*ComponentConfig{}, nil
	}

	names := config.GetSectionNames()
	result := make([]*ComponentConfig, len(names))
	for i, v := range names {
		c := config.GetSection(v)
		componentConfig, err := ReadComponentConfigFromConfig(c)
		if err != nil {
			return nil, err
		}
		result[i] = componentConfig
	}

	return result, nil
}
