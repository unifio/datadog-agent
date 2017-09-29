// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2017 Datadog, Inc.

package listeners

import (
	"github.com/DataDog/datadog-agent/pkg/util/docker"
)

// ID is the representation of the unique ID of a Service
type ID string

// DockerService reprensents a docker container we can run (a) check(s) against.
// It should be matched with a check template by the ConfigResolver using the
// ADIdentifiers field.
type DockerService struct {
	ID            ID                // unique ID
	ADIdentifiers []string          // identifiers on which templates will be matched
	Hosts         map[string]string // network --> IP address
	Ports         []int
	Tags          []string
	Pid           int // Process identifier
	container     docker.Container
}

// Service is an interface to get
type Service interface {
	GetID() ID
	GetADIdentifiers() ([]string, error)
	GetHosts() (map[string]string, error)
	GetPorts() ([]int, error)
	GetTags() ([]string, error)
	GetPid() (int, error)
}

// ServiceListener monitors running services and triggers check (un)scheduling
//
// It holds a cache of running services, listens to new/killed services and
// updates its cache, and the ConfigResolver with these events.
type ServiceListener interface {
	Listen(newSvc, delSvc chan<- Service)
	Stop()
}
