package container

import (
	"os"
	"os/signal"
	"strings"
	"fmt"
	"time"

	cconfig "github.com/pip-services3-go/pip-services3-commons-go/config"
	crefer "github.com/pip-services3-go/pip-services3-commons-go/refer"
	"github.com/pip-services3-go/pip-services3-components-go/log"
)

type ProcessContainer struct {
	Container
	configPath string
}

func NewEmptyProcessContainer() *ProcessContainer {
	c := &ProcessContainer {
		Container: *NewEmptyContainer(),
		configPath: "./config/config.yml",
	}
	c.SetLogger(log.NewConsoleLogger())
	return c
}

func NewProcessContainer(name string, description string) *ProcessContainer {
	c := &ProcessContainer {
		Container: *NewContainer(name, description),
		configPath: "./config/config.yml",
	}
	c.SetLogger(log.NewConsoleLogger())
	return c
}

func InheritProcessContainer(name string, description string,
	referenceable crefer.IReferenceable) *ProcessContainer {
	c := &ProcessContainer {
		Container: *InheritContainer(name, description, referenceable),
		configPath: "./config/config.yml",
	}
	c.SetLogger(log.NewConsoleLogger())
	return c
}

func (c *ProcessContainer) SetConfigPath(configPath string) {
	c.configPath = configPath
}

func (c *ProcessContainer) getConfigPath(args []string) string {
	for index, arg := range args {
		nextArg := ""
		if index < len(args) - 1 {
			nextArg = args[index + 1]
			if strings.HasPrefix(nextArg, "-") {
				nextArg = ""
			}
		}

		if arg == "--config" || arg == "-c" {
			return nextArg
		}
	}

	return c.configPath
}

func (c *ProcessContainer) getParameters(args []string) *cconfig.ConfigParams {
	line := ""

	for index := 0; index < len(args); index++ {
		arg := args[index]
		nextArg := ""
		if index < len(args) - 1 {
			nextArg = args[index + 1]
			if strings.HasPrefix(nextArg, "-") {
				nextArg = ""
			}
		}

		if nextArg != "" {
			if arg == "--param" || arg == "--params" || arg == "-p" {
				if line != "" {
					line = line + ";"
				}
				line = line + nextArg
				index++
			}
		}
	}

	parameters := cconfig.NewConfigParamsFromString(line)

	for _, e := range os.Environ() {
		env := strings.Split(e, "=")
		parameters.SetAsObject(env[0], env[1])
	}

	return parameters
}

func (c *ProcessContainer) showHelp(args []string) bool {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			return true
		}
	}
	return false
}

func (c *ProcessContainer) printHelp() {
	fmt.Println("Pip.Services process container - http://www.github.com/pip-services/pip-services")
	fmt.Println("run [-h] [-c <config file>] [-p <param>=<value>]*")
}

func (c *ProcessContainer) captureErrors(correlationId string) {
	if r := recover(); r != nil {
		err, _ := r.(error)
		c.Logger().Fatal(correlationId, err, "Process is terminated")
		os.Exit(1)
	}
}

func (c *ProcessContainer) captureExit(correlationId string) {
	c.Logger().Info(correlationId, "Press Control-C to stop the microservice...")

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	go func() {
		select {
		case <- ch:
			c.Close(correlationId)
			c.Logger().Info(correlationId, "Googbye!")
			os.Exit(0)
		}
	}()
}

func (c *ProcessContainer) Run(args []string) {
	if c.showHelp(args) {
		c.printHelp()
		os.Exit(0)
		return
	}

	correlationId := c.Info().Name
	path := c.getConfigPath(args)
	parameters := c.getParameters(args)

	err := c.ReadConfigFromFile(correlationId, path, parameters)
	if err != nil {
		c.Logger().Fatal(correlationId, err, "Process is terminated")
		os.Exit(1)
		return
	}

	defer c.captureErrors(correlationId)
	c.captureExit(correlationId)

	err = c.Open(correlationId)
	if err != nil {
		c.Logger().Fatal(correlationId, err, "Process is terminated")
		os.Exit(1)
		return
	}

	for {
		time.Sleep(100)
	}
}