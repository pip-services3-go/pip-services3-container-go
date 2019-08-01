package config

import (
	"github.com/pip-services3-go/pip-services3-commons-go/config"
	"github.com/pip-services3-go/pip-services3-commons-go/errors"
	"github.com/pip-services3-go/pip-services3-commons-go/refer"
	"github.com/pip-services3-go/pip-services3-commons-go/reflect"
)

type ComponentConfig struct {
	Descriptor *refer.Descriptor
	Type       *reflect.TypeDescriptor
	Config     *config.ConfigParams
}

func NewComponentConfigFromDescriptor(descriptor *refer.Descriptor,
	config *config.ConfigParams) *ComponentConfig {
	return &ComponentConfig{
		Descriptor: descriptor,
		Config:     config,
	}
}

func NewComponentConfigFromType(typ *reflect.TypeDescriptor,
	config *config.ConfigParams) *ComponentConfig {
	return &ComponentConfig{
		Type:   typ,
		Config: config,
	}
}

func ReadComponentConfigFromConfig(config *config.ConfigParams) (result *ComponentConfig, err error) {
	descriptor, err1 := refer.ParseDescriptorFromString(config.GetAsString("descriptor"))
	if err1 != nil {
		return nil, err1
	}

	typ, err2 := reflect.ParseTypeDescriptorFromString(config.GetAsString("type"))
	if err2 != nil {
		return nil, err2
	}

	if descriptor == nil && typ == nil {
		err = errors.NewConfigError("", "BAD_CONFIG", "Component configuration must have descriptor or type")
		return nil, err
	}

	return &ComponentConfig{
		Descriptor: descriptor,
		Type:       typ,
		Config:     config,
	}, nil
}
