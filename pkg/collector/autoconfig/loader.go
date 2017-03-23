package autoconfig

import (
	"fmt"

	"github.com/DataDog/datadog-agent/pkg/collector/autoconfig/containers"
	"github.com/DataDog/datadog-agent/pkg/collector/check"
)

// ConfigLoader is responsible for interpolating Config objects' template variables
// with container and environment data, for loading resulting configs to create
// actual checks
type ConfigLoader struct {
	containerListener *ContainerListener
	configToChecks    map[check.ID][]check.ID // first ID is config ID, second is check ID
	containerToChecks map[containers.ID][]check.ID
	startCheckQ       chan check.Check
	stopCheckQ        chan check.ID
}

// NewConfigLoader creates a ConfigLoader
func NewConfigLoader(startQ chan check.Check, stopQ chan check.ID) *ConfigBuilder {
	configToChecks := make(map[check.ID][]check.ID)
	containerToChecks := make(map[containers.ID][]check.ID)

	c := make(chan containers.Container)
	ctrListener := NewContainerListener(c)

	cl := &ConfigLoader{ctrListener, configToChecks, containerToChecks, startQ, stopQ}

	go func() {
		for {
			ctr := <-c
			if ctr.IsAlive {
				configs, err := c.GetContainerConfigs(ctr)
				if err != nil {
					log.Error("TODO")
				}
				for _, config := range configs {
					check, err := cl.loadCheck(config)
					if err != nil {
						log.Error("TODO")
					}
					startQ <- check
				}
			} else {
				checks, err := c.GetContainerChecks(ctr)
				if err != nil {
					log.Error("TODO")
				}
				for _, check := range checks {
					stopQ <- check.ID
				}
			}
		}
	}()

	return cl
}

// Schedule takes a config, interpolates its templates variables if any,
// loads it into a check, and schedules it.
func (l *ConfigLoader) Schedule(config check.Config) check.Check {
	if isTemplate(config) == true && l.configToChecks[config.ID] == nil {
		l.configToChecks = append(l.templates, config)
		config = l.fillTemplate(config)
	}

	check := loadCheck(check.Config)

	startCheckQ <- check
	return check
}

// Invalidate is used to notify ConfigBuilder that a config has been deleted
// It takes care of unscheduling the checks associated with said config.
func (l *ConfigLoader) Invalidate(config check.Config) error {
	for configId, checkIds := range l.configToChecks {
		if configId == config.ID {
			for _, cId := range checkIds {
				stopCheckQ <- cId
			}
			delete(l.configToChecks, configId)
			return nil
		}
	}
	return fmt.Errorf("template not found")
}

func (l *ConfigLoader) GetContainerConfigs(c containers.Container) ([]check.Config, error) {
	// TODO
}

func (l *ConfigLoader) GetContainerChecks(c containers.Container) ([]check.Checks, error) {
	// TODO
}

// isTemplate introspects a Config and returns true if
// it has an identifier, a single instance, and some template variables.
func isTemplate(config check.Config) bool {
	if config.Identifier == "" || len(config.Instances) != 1 ||
		!hasTemplateVar(config.InitConfig) || !hasTemplateVar(config.Instances[0]) {
		return false
	}
	return true
}

// fillTemplate takes a templated Config and interpolates variables
// with data it pulls from the environment (docker lables, k8s annoations, etc.)
func (l *ConfigLoader) fillTemplate(template check.Config) check.Config {
	// TODO
}

func (l *ConfigLoader) loadCheck(check.Config) (check.Check, error) {
	// TODO
}

// AddContainer updates the internal container cache with a new container
// and tries and build new configs by matching this containers with existing templates
func (l *ConfigLoader) AddContainer(containers.Container) []check.Config {
	// TODO
	return []check.Config{}
}

// RemoveContainer deletes a container from the container cache
// and returns IDs of checks that need to be unscheduled
func (l *ConfigLoader) RemoveContainer(containers.Container) []check.ID {
	// TODO
	return []check.ID{}
}
