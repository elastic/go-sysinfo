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

//go:build freebsd

package freebsd

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/elastic/go-sysinfo/types"

	"golang.org/x/sys/unix"
)

var tickDuration = sync.OnceValues(func() (time.Duration, error) {
	const mib = "kern.clockrate"

	c, err := unix.SysctlClockinfo(mib)
	if err != nil {
		return 0, fmt.Errorf("failed to get %s: %w", mib, err)
	}
	return time.Duration(c.Tick) * time.Microsecond, nil
})

var pageSizeBytes = sync.OnceValues(func() (uint64, error) {
	const mib = "vm.stats.vm.v_page_size"

	v, err := unix.SysctlUint32(mib)
	if err != nil {
		return 0, fmt.Errorf("failed to get %s: %w", mib, err)
	}

	return uint64(v), nil
})

func activePageCount() (uint64, error) {
	const mib = "vm.stats.vm.v_active_count"

	v, err := unix.SysctlUint32(mib)
	if err != nil {
		return 0, fmt.Errorf("failed to get %s: %w", mib, err)
	}
	return uint64(v), nil
}

func architecture() (string, error) {
	const mib = "hw.machine"

	arch, err := unix.Sysctl(mib)
	if err != nil {
		return "", fmt.Errorf("failed to get architecture: %w", err)
	}

	return arch, nil
}

func bootTime() (time.Time, error) {
	const mib = "kern.boottime"

	tv, err := unix.SysctlTimeval(mib)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get host uptime: %w", err)
	}

	bootTime := time.Unix(tv.Sec, tv.Usec*int64(time.Microsecond))
	return bootTime, nil
}

// buffersUsedBytes returns the number memory bytes used as disk cache.
func buffersUsedBytes() (uint64, error) {
	const mib = "vfs.bufspace"

	v, err := unix.SysctlUint64(mib)
	if err != nil {
		return 0, fmt.Errorf("failed to get %s: %w", mib, err)
	}

	return v, nil
}

func cachePageCount() (uint64, error) {
	const mib = "vm.stats.vm.v_cache_count"

	v, err := unix.SysctlUint32(mib)
	if err != nil {
		return 0, fmt.Errorf("failed to get %s: %w", mib, err)
	}

	return uint64(v), nil
}

const sizeOfUint64 = int(unsafe.Sizeof(uint64(0)))

// cpuStateTimes uses sysctl kern.cp_time to get the amount of time spent in
// different CPU states.
func cpuStateTimes() (*types.CPUTimes, error) {
	tickDuration, err := tickDuration()
	if err != nil {
		return nil, err
	}

	const mib = "kern.cp_time"
	buf, err := unix.SysctlRaw("kern.cp_time")
	if err != nil {
		return nil, fmt.Errorf("failed to get %s: %w", mib, err)
	}

	var clockTicks [unix.CPUSTATES]uint64
	if len(buf) < len(clockTicks)*sizeOfUint64 {
		return nil, fmt.Errorf("kern.cp_time data is too short (got %d bytes)", len(buf))
	}
	for i := range clockTicks {
		val := *(*uint64)(unsafe.Pointer(&buf[sizeOfUint64*i]))
		clockTicks[i] = val
	}

	return &types.CPUTimes{
		User:   time.Duration(clockTicks[unix.CP_USER]) * tickDuration,
		System: time.Duration(clockTicks[unix.CP_SYS]) * tickDuration,
		Idle:   time.Duration(clockTicks[unix.CP_IDLE]) * tickDuration,
		IRQ:    time.Duration(clockTicks[unix.CP_INTR]) * tickDuration,
		Nice:   time.Duration(clockTicks[unix.CP_NICE]) * tickDuration,
	}, nil
}

func freePageCount() (uint64, error) {
	const mib = "vm.stats.vm.v_free_count"

	v, err := unix.SysctlUint32(mib)
	if err != nil {
		return 0, fmt.Errorf("failed to get %s: %w", mib, err)
	}

	return uint64(v), nil
}

func inactivePageCount() (uint64, error) {
	const mib = "vm.stats.vm.v_inactive_count"

	v, err := unix.SysctlUint32(mib)
	if err != nil {
		return 0, fmt.Errorf("failed to get %s: %w", mib, err)
	}

	return uint64(v), nil
}

func kernelVersion() (string, error) {
	const mib = "kern.osrelease"

	version, err := unix.Sysctl(mib)
	if err != nil {
		return "", fmt.Errorf("failed to get kernel version: %w", err)
	}

	return version, nil
}

func machineID() (string, error) {
	const mib = "kern.hostuuid"

	uuid, err := unix.Sysctl(mib)
	if err != nil {
		return "", fmt.Errorf("failed to get machine id: %w", err)
	}

	return uuid, nil
}

func operatingSystem() (*types.OSInfo, error) {
	info := &types.OSInfo{
		Type:     "freebsd",
		Family:   "freebsd",
		Platform: "freebsd",
	}

	osType, err := unix.Sysctl("kern.ostype")
	if err != nil {
		return info, err
	}
	info.Name = osType

	// Example: 13.0-RELEASE-p11
	osRelease, err := unix.Sysctl("kern.osrelease")
	if err != nil {
		return info, err
	}
	info.Version = osRelease

	releaseParts := strings.Split(osRelease, "-")

	majorMinor := strings.Split(releaseParts[0], ".")
	if len(majorMinor) > 0 {
		info.Major, _ = strconv.Atoi(majorMinor[0])
	}
	if len(majorMinor) > 1 {
		info.Minor, _ = strconv.Atoi(majorMinor[1])
	}

	if len(releaseParts) > 1 {
		info.Build = releaseParts[1]
	}
	if len(releaseParts) > 2 {
		info.Patch, _ = strconv.Atoi(strings.TrimPrefix(releaseParts[2], "p"))
	}

	return info, nil
}

func totalPhysicalMem() (uint64, error) {
	const mib = "hw.physmem"

	v, err := unix.SysctlUint64(mib)
	if err != nil {
		return 0, fmt.Errorf("failed to get %s: %w", mib, err)
	}
	return v, nil
}

func wirePageCount() (uint64, error) {
	const mib = "vm.stats.vm.v_wire_count"

	v, err := unix.SysctlUint32(mib)
	if err != nil {
		return 0, fmt.Errorf("failed to get %s: %w", mib, err)
	}
	return uint64(v), nil
}
