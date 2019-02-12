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

// +build,freebsd,cgo

package freebsd

// #cgo LDFLAGS: -lkvm
//#include <kvm.h>
//#include <sys/vmmeter.h>
import "C"

import (
	"os"
	"time"

	"github.com/joeshaw/multierror"
	"github.com/pkg/errors"

	"github.com/dbolcsfoldi/go-sysinfo/internal/registry"
	"github.com/elastic/go-sysinfo/providers/shared"
	"github.com/elastic/go-sysinfo/types"
)

func init() {
	registry.Register(freebsdSystem{})
}

type freebsdSystem struct{}

func (s freebsdSystem) Host() (types.Host, error) {
	return newHost()
}

type host struct {
	info types.HostInfo
}

func (h *host) Info() types.HostInfo {
	return h.info
}

func (h *host) CPUTime() (types.CPUTimes, error) {
	cpu := types.CPUTimes{}
	r := &reader{}
	r.cpuTime(&cpu)

	return cpu, nil
}

func (h *host) Memory() (*types.HostMemoryInfo, error) {
	m := &types.HostMemoryInfo{}
	r := &reader{}
	r.memInfo(m)
	return m, r.Err()
}

func newHost() (*host, error) {
	h := &host{}
	r := &reader{}
	r.architecture(h)
	r.bootTime(h)
	r.hostname(h)
	r.network(h)
	r.kernelVersion(h)
	r.os(h)
	r.time(h)
	r.uniqueID(h)
	return h, r.Err()
}

type reader struct {
	errs []error
}

func (r *reader) addErr(err error) bool {
	if err != nil {
		if errors.Cause(err) != types.ErrNotImplemented {
			r.errs = append(r.errs, err)
		}
		return true
	}
	return false
}

func (r *reader) Err() error {
	if len(r.errs) > 0 {
		return &multierror.MultiError{Errors: r.errs}
	}
	return nil
}

func (r *reader) cpuTime(cpu *types.CPUTimes) {
	cptime, err := Cptime()

	if r.addErr(err) {
		return
	}

	cpu.User = time.Duration(cptime["User"])
	cpu.System = time.Duration(cptime["System"])
	cpu.Idle = time.Duration(cptime["Idle"])
	cpu.Nice = time.Duration(cptime["Nice"])
	cpu.IRQ = time.Duration(cptime["IRQ"])
}

func (r *reader) memInfo(m *types.HostMemoryInfo) {
	pageSize, err := PageSize()

	if r.addErr(err) {
		return
	}

	totalMemory, err := TotalMemory()
	if r.addErr(err) {
		return
	}

	m.Total = totalMemory

	vm, err := VmTotal()
	if r.addErr(err) {
		return
	}

	m.Free = uint64(vm.Free) * uint64(pageSize)
	m.Used = m.Total - m.Free

	numFreeBuffers, err := NumFreeBuffers()
	if r.addErr(err) {
		return
	}

	m.Available = m.Free + (uint64(numFreeBuffers) * uint64(pageSize))

	swap, err := KvmGetSwapInfo()
	if r.addErr(err) {
		return
	}

	swapMaxPages, err := SwapMaxPages()
	if r.addErr(err) {
		return
	}

	if swap.Total > swapMaxPages {
		swap.Total = swapMaxPages
	}

	m.VirtualTotal = uint64(swap.Total) * uint64(pageSize)
	m.VirtualUsed = uint64(swap.Used) * uint64(pageSize)
	m.VirtualFree = m.VirtualTotal - m.VirtualUsed
}

func (r *reader) architecture(h *host) {
	v, err := Architecture()
	if r.addErr(err) {
		return
	}
	h.info.Architecture = v
}

func (r *reader) bootTime(h *host) {
	v, err := BootTime()
	if r.addErr(err) {
		return
	}

	h.info.BootTime = v
}

func (r *reader) hostname(h *host) {
	v, err := os.Hostname()
	if r.addErr(err) {
		return
	}
	h.info.Hostname = v
}

func (r *reader) network(h *host) {
	ips, macs, err := shared.Network()
	if r.addErr(err) {
		return
	}
	h.info.IPs = ips
	h.info.MACs = macs
}

func (r *reader) kernelVersion(h *host) {
	v, err := KernelVersion()
	if r.addErr(err) {
		return
	}
	h.info.KernelVersion = v
}

func (r *reader) os(h *host) {
	v, err := OperatingSystem()
	if r.addErr(err) {
		return
	}
	h.info.OS = v
}

func (r *reader) time(h *host) {
	h.info.Timezone, h.info.TimezoneOffsetSec = time.Now().Zone()
}

func (r *reader) uniqueID(h *host) {
	v, err := MachineID()
	if r.addErr(err) {
		return
	}
	h.info.UniqueID = v
}
