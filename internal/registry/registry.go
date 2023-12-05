// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package registry

import (
	"fmt"

	"github.com/elastic/go-sysinfo/types"
)

type HostFSCreator = func(string) HostProvider
type ProcessFSCreator = func(string) ProcessProvider

// HostProvider defines interfaces that provide host-specific metrics
type HostProvider interface {
	Host() (types.Host, error)
}

// ProcessProvider defines interfaces that provide process-specific metrics
type ProcessProvider interface {
	Processes() ([]types.Process, error)
	Process(pid int) (types.Process, error)
	Self() (types.Process, error)
}

var (
	hostProvider            HostProvider
	hostProviderWithRoot    HostFSCreator
	processProvider         ProcessProvider
	processProviderWithRoot ProcessFSCreator
)

// Register a metrics provider. `provider` should implement one or more of `ProcessProvider`, `HostProvider`, `HostFSCreator` or `ProcessFSCreator`
func Register(provider interface{}) {
	if h, ok := provider.(HostProvider); ok {
		if hostProvider != nil {
			panic(fmt.Sprintf("HostProvider already registered: %v", hostProvider))
		}
		hostProvider = h
	}

	if p, ok := provider.(ProcessProvider); ok {
		if processProvider != nil {
			panic(fmt.Sprintf("ProcessProvider already registered: %v", processProvider))
		}
		processProvider = p
	}

	if creator, ok := provider.(HostFSCreator); ok {
		if hostProviderWithRoot != nil {
			panic("hostProviderWithRoot already registered")
		}
		hostProviderWithRoot = creator
	}

	if creator, ok := provider.(ProcessFSCreator); ok {
		if processProviderWithRoot != nil {
			panic("processProviderWithRoot is already registered")
		}
		processProviderWithRoot = creator
	}
}

// GetHostProvider returns the HostProvider registered for the system. May return nil.
func GetHostProvider() HostProvider {
	if hostProviderWithRoot != nil {
		return hostProviderWithRoot("")
	}
	return hostProvider
}

// GetHostProviderWithRoot creates a host provider for the given sysfs root. May return nil.
func GetHostProviderWithRoot(hostFS string) HostProvider {
	if hostProviderWithRoot != nil {
		return hostProviderWithRoot(hostFS)
	}
	return hostProvider
}

// GetProcessProvider returns the ProcessProvider registered on the system. May return nil.
func GetProcessProvider() ProcessProvider {
	if processProviderWithRoot != nil {
		return processProviderWithRoot("")
	}
	return processProvider
}

// GetProcessProviderWithRoot creates a process provider for the sysfs root. May return nil
func GetProcessProviderWithRoot(hostFS string) ProcessProvider {
	if processProviderWithRoot != nil {
		return processProviderWithRoot(hostFS)
	}
	return processProvider
}
