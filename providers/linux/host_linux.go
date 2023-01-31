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

package linux

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joeshaw/multierror"
	"github.com/prometheus/procfs"

	"github.com/elastic/go-sysinfo/internal/registry"
	"github.com/elastic/go-sysinfo/providers/shared"
	"github.com/elastic/go-sysinfo/types"
)

func init() {
	registry.Register(newLinuxSystem(""))
}

type linuxSystem struct {
	procFS procFS
}

func newLinuxSystem(hostFS string) linuxSystem {
	mountPoint := filepath.Join(hostFS, procfs.DefaultMountPoint)
	fs, _ := procfs.NewFS(mountPoint)
	return linuxSystem{
		procFS: procFS{FS: fs, mountPoint: mountPoint},
	}
}

func (s linuxSystem) Host() (types.Host, error) {
	return newHost(s.procFS)
}

type host struct {
	procFS procFS
	stat   procfs.Stat
	info   types.HostInfo
}

func (h *host) Info() types.HostInfo {
	return h.info
}

func (h *host) Memory() (*types.HostMemoryInfo, error) {
	content, err := ioutil.ReadFile(h.procFS.path("meminfo"))
	if err != nil {
		return nil, err
	}

	return parseMemInfo(content)
}

// VMStat reports data from /proc/vmstat on linux.
func (h *host) VMStat() (*types.VMStatInfo, error) {
	content, err := ioutil.ReadFile(h.procFS.path("vmstat"))
	if err != nil {
		return nil, err
	}

	return parseVMStat(content)
}

// LoadAverage reports data from /proc/loadavg on linux.
func (h *host) LoadAverage() (*types.LoadAverageInfo, error) {
	loadAvg, err := h.procFS.LoadAvg()
	if err != nil {
		return nil, err
	}

	return &types.LoadAverageInfo{
		One:     loadAvg.Load1,
		Five:    loadAvg.Load5,
		Fifteen: loadAvg.Load15,
	}, nil
}

// NetworkCounters reports data from /proc/net on linux
func (h *host) NetworkCounters() (*types.NetworkCountersInfo, error) {
	snmpRaw, err := ioutil.ReadFile(h.procFS.path("net/snmp"))
	if err != nil {
		return nil, err
	}
	snmp, err := getNetSnmpStats(snmpRaw)
	if err != nil {
		return nil, err
	}

	netstatRaw, err := ioutil.ReadFile(h.procFS.path("net/netstat"))
	if err != nil {
		return nil, err
	}
	netstat, err := getNetstatStats(netstatRaw)
	if err != nil {
		return nil, err
	}

	return &types.NetworkCountersInfo{SNMP: snmp, Netstat: netstat}, nil
}

func (h *host) CPUTime() (types.CPUTimes, error) {
	stat, err := h.procFS.Stat()
	if err != nil {
		return types.CPUTimes{}, err
	}

	return types.CPUTimes{
		User:    time.Duration(stat.CPUTotal.User * float64(time.Second)),
		System:  time.Duration(stat.CPUTotal.System * float64(time.Second)),
		Idle:    time.Duration(stat.CPUTotal.Idle * float64(time.Second)),
		IOWait:  time.Duration(stat.CPUTotal.Iowait * float64(time.Second)),
		IRQ:     time.Duration(stat.CPUTotal.IRQ * float64(time.Second)),
		Nice:    time.Duration(stat.CPUTotal.Nice * float64(time.Second)),
		SoftIRQ: time.Duration(stat.CPUTotal.SoftIRQ * float64(time.Second)),
		Steal:   time.Duration(stat.CPUTotal.Steal * float64(time.Second)),
	}, nil
}

func newHost(fs procFS) (*host, error) {
	stat, err := fs.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to read proc stat: %w", err)
	}

	h := &host{stat: stat, procFS: fs}
	r := &reader{}
	r.architecture(h)
	r.bootTime(h)
	r.containerized(h)
	r.hostname(h)
	r.domain(h)
	r.fqdn(h)
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
		return &multierror.MultiError{Errors: r.errs}
	}
	return nil
}

func (r *reader) architecture(h *host) {
	v, err := Architecture()
	if r.addErr(err) {
		return
	}
	h.info.Architecture = v
}

func (r *reader) bootTime(h *host) {
	v, err := bootTime(h.procFS.FS)
	if r.addErr(err) {
		return
	}
	h.info.BootTime = v
}

func (r *reader) containerized(h *host) {
	v, err := IsContainerized()
	if r.addErr(err) {
		return
	}
	h.info.Containerized = &v
}

func (r *reader) hostname(h *host) {
	v, err := os.Hostname()
	if r.addErr(err) {
		return
	}
	h.info.Hostname = v
}

func (r *reader) domain(h *host) {
	v, err := domainname()
	if r.addErr(err) {
		return
	}
	h.info.Domain = v
}

const etcHosts = "/etc/hosts"

func (r *reader) fqdn(h *host) {
	f, err := os.Open(etcHosts)
	if err != nil {
		r.addErr(fmt.Errorf("could open %q to get FQDN: %w", etcHosts, err))
		return
	}

	hname, err := os.Hostname()
	if err != nil {
		r.addErr(fmt.Errorf("could get hostname to look for FQDN: %w", err))
		return
	}

	fqdn, err := fqdnFromHosts(hname, f)
	if err != nil {
		r.addErr(fmt.Errorf("error when looking for FQDN on %s: %w", etcHosts, err))
		return
	}

	if fqdn == "" {
		// FQDN not found on hosts file, fall back to net.Lookup?
		// add an error?
	}
	h.info.FQDN = fqdn
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

type procFS struct {
	procfs.FS
	mountPoint string
}

func (fs *procFS) path(p ...string) string {
	elem := append([]string{fs.mountPoint}, p...)
	return filepath.Join(elem...)
}

// fqdnFromHosts looks for the FQDN for hostname on hostFile.
// If successfully it returns FQDN, nil. If no FQDN for hostname is found
// it returns "", nil. It returns "", err if any error happens.
func fqdnFromHosts(hostname string, hostsFile fs.File) (string, error) {
	s := bufio.NewScanner(hostsFile)

	for s.Scan() {
		fqdn := findInHostsLine(hostname, s.Text())
		if fqdn != "" {
			return fqdn, nil
		}
	}
	if err := s.Err(); err != nil {
		return "", fmt.Errorf("error reading hosts file lines: %w", err)
	}

	return "", nil
}

// findInHostsLine takes a HOSTS(5) line and searches for an alias matching
// hostname, if found it returns the canonical_hostname. The canonical_hostname
// should be the FQDN, see HOSTNAME(1).
// TODO: check k8s: https://kubernetes.io/docs/tasks/network/customize-hosts-file-for-pods/
func findInHostsLine(hostname, hostsEntry string) string {
	line, _, _ := strings.Cut(hostsEntry, "#")
	if len(line) < 1 {
		fmt.Printf("skip comment or empty: %q\n", hostsEntry)
		return ""
	}

	fileds := strings.FieldsFunc(line, func(r rune) bool {
		return r == ' ' || r == '\t'
	})

	if len(fileds) < 2 {
		// invalid hostsEntry
		return ""
	}

	// fields[0] is the ip address
	cannonical, aliases := fileds[1], fileds[1:]

	// TODO: confirm: a name should not repeat on different addresses.
	if len(fileds) == 2 {
		if fileds[1] == hostname {
			return cannonical
		}

		// If hostname was not set as an alias for FQDN, but the fist name
		// before the dot is the hostname:
		//   192.168.1.10    foo.mydomain.org	#  foo
		if hname, _, _ := strings.Cut(cannonical, "."); hname == hostname {
			return cannonical
		}
	}

	for _, h := range aliases {
		if h == hostname {
			return cannonical
		}
	}

	return ""
}
