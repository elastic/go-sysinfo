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

//go:build freebsd && cgo

package freebsd

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/prometheus/procfs"

	"github.com/elastic/go-sysinfo/internal/registry"
	"github.com/elastic/go-sysinfo/providers/shared"
	"github.com/elastic/go-sysinfo/types"
)

func init() {
	registry.Register(newFreeBSDSystem())
}

type freebsdSystem struct{}

func newFreeBSDSystem() freebsdSystem {
	return freebsdSystem{}
}

func (s freebsdSystem) Host() (types.Host, error) {
	return newHost()
}

type host struct {
	procFS procFS
	info   types.HostInfo
}

func (h *host) Info() types.HostInfo {
	return h.info
}

func (h *host) CPUTime() (types.CPUTimes, error) {
	cpu := types.CPUTimes{}
	r := &reader{}
	r.cpuTime(&cpu)
	return cpu, r.Err()
}

func (h *host) Memory() (*types.HostMemoryInfo, error) {
	m := &types.HostMemoryInfo{}
	r := &reader{}
	r.memInfo(m)
	return m, r.Err()
}

func (h *host) FQDNWithContext(ctx context.Context) (string, error) {
	return shared.FQDNWithContext(ctx)
}

func (h *host) FQDN() (string, error) {
	return h.FQDNWithContext(context.Background())
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
		if !errors.Is(err, types.ErrNotImplemented) {
			r.errs = append(r.errs, err)
		}
		return true
	}
	return false
}

func (r *reader) Err() error {
	if len(r.errs) > 0 {
		return errors.Join(r.errs...)
	}
	return nil
}

func (r *reader) cpuTime(cpu *types.CPUTimes) {
	times, err := cpuStateTimes()
	if r.addErr(err) {
		return
	}
	*cpu = *times
}

func (r *reader) memInfo(m *types.HostMemoryInfo) {
	ps, err := pageSizeBytes()
	if r.addErr(err) {
		return
	}
	pageSize := uint64(ps)

	m.Total, err = totalPhysicalMem()
	if r.addErr(err) {
		return
	}

	activePages, err := activePageCount()
	if r.addErr(err) {
		return
	}
	m.Metrics = make(map[string]uint64, 6)
	m.Metrics["active_bytes"] = uint64(activePages) * pageSize

	wirePages, err := wirePageCount()
	if r.addErr(err) {
		return
	}
	m.Metrics["wired_bytes"] = uint64(wirePages) * pageSize

	inactivePages, err := inactivePageCount()
	if r.addErr(err) {
		return
	}
	m.Metrics["inactive_bytes"] = uint64(inactivePages) * pageSize

	cachePages, err := cachePageCount()
	if r.addErr(err) {
		return
	}
	m.Metrics["cache_bytes"] = uint64(cachePages) * pageSize

	freePages, err := freePageCount()
	if r.addErr(err) {
		return
	}
	m.Metrics["free_bytes"] = uint64(freePages) * pageSize

	buffers, err := buffersUsedBytes()
	if r.addErr(err) {
		return
	}
	m.Metrics["buffer_bytes"] = buffers

	m.Used = uint64(activePages+wirePages) * pageSize
	m.Free = uint64(freePages) * pageSize
	m.Available = uint64(inactivePages+cachePages+freePages)*pageSize + buffers

	// Virtual (swap) Memory
	swap, err := kvmGetSwapInfo()
	if r.addErr(err) {
		return
	}

	m.VirtualTotal = uint64(swap.Total) * pageSize
	m.VirtualUsed = uint64(swap.Used) * pageSize
	m.VirtualFree = m.VirtualTotal - m.VirtualUsed
}

func (r *reader) architecture(h *host) {
	v, err := architecture()
	if r.addErr(err) {
		return
	}
	h.info.Architecture = v
}

func (r *reader) bootTime(h *host) {
	v, err := bootTime()
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
	v, err := kernelVersion()
	if r.addErr(err) {
		return
	}
	h.info.KernelVersion = v
}

func (r *reader) os(h *host) {
	v, err := operatingSystem()
	if r.addErr(err) {
		return
	}
	h.info.OS = v
}

func (r *reader) time(h *host) {
	h.info.Timezone, h.info.TimezoneOffsetSec = time.Now().Zone()
}

func (r *reader) uniqueID(h *host) {
	v, err := machineID()
	if r.addErr(err) {
		return
	}
	h.info.UniqueID = v
}

type procFS struct {
	procfs.FS
	mountPoint string
}

func (fs *procFS) path(p ...string) string {
	elem := append([]string{fs.mountPoint}, p...)
	return filepath.Join(elem...)
}
