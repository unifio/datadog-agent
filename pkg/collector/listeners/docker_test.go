// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2017 Datadog, Inc.

package listeners

import (
	"testing"
	"time"

	"github.com/DataDog/datadog-agent/pkg/util/cache"
	"github.com/DataDog/datadog-agent/pkg/util/docker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
)

func TestGetConfigIDFromPs(t *testing.T) {
	co := types.Container{
		ID:    "deadbeef",
		Image: "test",
	}
	dl := DockerListener{}

	ids := dl.getConfigIDFromPs(co)
	assert.Len(t, ids, 1)
	assert.Equal(t, "test", ids[0])

	labeledCo := types.Container{
		ID:     "deadbeef",
		Image:  "test",
		Labels: map[string]string{"io.datadog.check.id": "w00tw00t"},
	}
	ids = dl.getConfigIDFromPs(labeledCo)
	assert.Len(t, ids, 1)
	assert.Equal(t, "w00tw00t", ids[0])
}

func TestGetHostsFromPs(t *testing.T) {
	dl := DockerListener{}

	co := types.Container{
		ID:    "foo",
		Image: "test",
	}

	assert.Empty(t, dl.getHostsFromPs(co))

	nets := make(map[string]*network.EndpointSettings)
	nets["bridge"] = &network.EndpointSettings{IPAddress: "172.17.0.2"}
	nets["foo"] = &network.EndpointSettings{IPAddress: "172.17.0.3"}
	networkSettings := types.SummaryNetworkSettings{
		Networks: nets}

	co = types.Container{
		ID:              "deadbeef",
		Image:           "test",
		NetworkSettings: &networkSettings,
		Ports:           []types.Port{types.Port{PrivatePort: 1337}, types.Port{PrivatePort: 42}},
	}
	hosts := dl.getHostsFromPs(co)

	assert.Equal(t, "172.17.0.2", hosts["bridge"])
	assert.Equal(t, "172.17.0.3", hosts["foo"])
	assert.Equal(t, 2, len(hosts))
}

func TestGetPortsFromPs(t *testing.T) {
	dl := DockerListener{}

	co := types.Container{
		ID:    "foo",
		Image: "test",
	}
	assert.Empty(t, dl.getPortsFromPs(co))

	co.Ports = make([]types.Port, 0)
	assert.Empty(t, dl.getPortsFromPs(co))

	co.Ports = append(co.Ports, types.Port{PrivatePort: 1234})
	co.Ports = append(co.Ports, types.Port{PrivatePort: 4321})
	ports := dl.getPortsFromPs(co)
	assert.Equal(t, 2, len(ports))
	assert.Contains(t, ports, 1234)
	assert.Contains(t, ports, 4321)
}

// TODO Refactor
// func TestGetADIdentifiers(t *testing.T) {
// 	co := types.ContainerJSON{
// 		ContainerJSONBase: &types.ContainerJSONBase{ID: "deadbeef", Image: "test"},
// 		Mounts:            make([]types.MountPoint, 0),
// 		Config:            &container.Config{},
// 		NetworkSettings:   &types.NetworkSettings{},
// 	}
// 	dl := DockerListener{}

// 	ids := dl.getConfigIDFromInspect(co)
// 	assert.Len(t, ids, 1)
// 	assert.Equal(t, "test", ids[0])

// 	labeledCo := types.ContainerJSON{
// 		ContainerJSONBase: &types.ContainerJSONBase{ID: "deadbeef", Image: "test"},
// 		Mounts:            make([]types.MountPoint, 0),
// 		Config:            &container.Config{Labels: map[string]string{"io.datadog.check.id": "w00tw00t"}},
// 		NetworkSettings:   &types.NetworkSettings{},
// 	}
// 	ids = dl.getConfigIDFromInspect(labeledCo)
// 	assert.Len(t, ids, 1)
// 	assert.Equal(t, "w00tw00t", ids[0])
// }

func TestGetHosts(t *testing.T) {
	id := "fooooooooooo"
	cBase := types.ContainerJSONBase{
		ID:    id,
		Image: "test",
	}
	cj := types.ContainerJSON{
		ContainerJSONBase: &cBase,
		Mounts:            make([]types.MountPoint, 0),
		Config:            &container.Config{Labels: map[string]string{"io.datadog.check.id": "w00tw00t"}},
		NetworkSettings:   &types.NetworkSettings{},
	}
	// add cj to the cache to avoir having to query docker in the test
	cacheKey := docker.GetInspectCacheKey(id)
	cache.Cache.Set(cacheKey, cj, 10*time.Second)

	co := docker.Container{ID: id}
	svc := DockerService{
		ID:        ID(id),
		container: co,
	}

	res, _ := svc.GetHosts()
	assert.Empty(t, res)

	nets := make(map[string]*network.EndpointSettings)
	nets["bridge"] = &network.EndpointSettings{IPAddress: "172.17.0.2"}
	nets["foo"] = &network.EndpointSettings{IPAddress: "172.17.0.3"}
	ports := make(nat.PortMap)
	p, _ := nat.NewPort("tcp", "1337")
	ports[p] = make([]nat.PortBinding, 0)
	p, _ = nat.NewPort("tcp", "42")
	ports[p] = make([]nat.PortBinding, 0)

	id = "deadbeefffff"
	cBase = types.ContainerJSONBase{
		ID:    id,
		Image: "test",
	}
	networkSettings := types.NetworkSettings{
		NetworkSettingsBase: types.NetworkSettingsBase{Ports: ports},
		Networks:            nets,
	}

	cj = types.ContainerJSON{
		ContainerJSONBase: &cBase,
		Mounts:            make([]types.MountPoint, 0),
		Config:            &container.Config{},
		NetworkSettings:   &networkSettings,
	}
	// update cj in the cache
	cacheKey = docker.GetInspectCacheKey(id)
	cache.Cache.Set(cacheKey, cj, 10*time.Second)

	co = docker.Container{ID: id}
	svc = DockerService{
		ID:        ID(id),
		container: co,
	}
	hosts, _ := svc.GetHosts()

	assert.Equal(t, "172.17.0.2", hosts["bridge"])
	assert.Equal(t, "172.17.0.3", hosts["foo"])
	assert.Equal(t, 2, len(hosts))
}

func TestGetPorts(t *testing.T) {
	cBase := types.ContainerJSONBase{
		ID:    "deadbeefffff",
		Image: "test",
	}

	ports := make(nat.PortMap)
	networkSettings := types.NetworkSettings{
		NetworkSettingsBase: types.NetworkSettingsBase{Ports: ports},
		Networks:            make(map[string]*network.EndpointSettings),
	}

	cj := types.ContainerJSON{
		ContainerJSONBase: &cBase,
		Mounts:            make([]types.MountPoint, 0),
		Config:            &container.Config{},
		NetworkSettings:   &networkSettings,
	}
	// add cj to the cache to avoir having to query docker in the test
	cacheKey := docker.GetInspectCacheKey("deadbeefffff")
	cache.Cache.Set(cacheKey, cj, 10*time.Second)

	co := docker.Container{ID: "deadbeefffff"}
	svc := DockerService{
		ID:        ID("deadbeefffff"),
		container: co,
	}
	svcPorts, _ := svc.GetPorts()
	assert.Empty(t, svcPorts)

	ports = make(nat.PortMap, 2)
	p, _ := nat.NewPort("tcp", "1234")
	ports[p] = nil
	p, _ = nat.NewPort("tcp", "4321")
	ports[p] = nil

	cj.NetworkSettings.Ports = ports
	// update cj in the cache
	svc.Ports = nil
	cache.Cache.Set(cacheKey, cj, 10*time.Second)

	svcPorts, _ = svc.GetPorts()
	assert.Equal(t, 2, len(svcPorts))
	assert.Contains(t, svcPorts, 1234)
	assert.Contains(t, svcPorts, 4321)
}
