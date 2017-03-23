package containers

import (
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

const (
	DockerBackend string = "docker"
)

type ID string

type Container struct {
	types.Container
	Backend string // only Docker for now
	IsAlive bool
	Err     error
}

// A ContainerListener listens to a stream of container-related
// events and create ConfigUpdate messages that it pushes on a channel
// for Auto Config to consume.
type ContainerListener struct {
	containerChanges chan Container
}

// NewContainerListener creates a ContainerListener.
// For now the only supported listener is the Docker one.
func NewContainerListener(containerChanges chan Container) *ContainerListener {
	// backend := config.Datadog.GetString("autoconf_backend")
	return &ContainerListener{containerChanges}
}

// Listen waits for Docker events and transmit relevant
// ones (see filters) to AutoConfig for config update.
// TODO: retrieve connect config (from docker check settings)?
func (listener *ContainerListener) Listen() {
	c, err := client.NewEnvClient()
	if err != nil {
		listener.containerChanges <- Container{nil, DockerBackend, err}
		return
	}

	// filters only match start/stop container events
	filters := filters.NewArgs()
	filters.Add("type", "container")
	filters.Add("event", "start")
	filters.Add("event", "die")
	eventOptions := types.EventsOptions{Since: time.Now().String(), Filters: filters}

	messages, errs := c.Events(context.Background(), eventOptions)

	for {
		select {
		case msg := <-messages:
			container := listener.GetContainer(msg)
			listener.containerChanges <- Container{listener.GetContainer(msg), DockerBackend, nil}
		case err := <-errs:
			if err != nil && err != io.EOF {
				listener.containerChanges <- Container{nil, DockerBackend, err}
				// TODO: do we need to break outta here?
			}
		}
	}
}

// GetContainer parses a docker Message and get the container concerned by it
func (listener *ContainerListener) GetContainer(events.Message) Container {
	return Container{}
}

// GetRunningContainers does pretty much what you would expect
// TODO: it probably shouldn't be there but in our Docker utils
func (listener *ContainerListener) GetRunningContainers() []Container {
	return []Container{}
}
