package types

import "time"

type Host interface {
	Info() HostInfo
}

type HostInfo struct {
	Architecture      string    `json:"architecture"`            // Hardware architecture (e.g. x86_64, arm, ppc, mips).
	BootTime          time.Time `json:"boot_time"`               // Host boot time.
	Containerized     *bool     `json:"containerized,omitempty"` // Is the process containerized.
	Hostname          string    `json:"hostname"`                // Hostname
	IPs               []string  `json:"ips,omitempty"`           // List of all IPs.
	KernelVersion     string    `json:"kernel_version"`          // Kernel version.
	MACs              []string  `json:"mac_addresses"`           // List of MAC addresses.
	OS                *OSInfo   `json:"os"`                      // OS information.
	Timezone          string    `json:"timezone"`                // System timezone.
	TimezoneOffsetSec int       `json:"timezone_offset_sec"`     // Timezone offset (seconds from UTC).
	UniqueID          string    `json:"id,omitempty"`            // Unique ID of the host (optional).
}

func (host HostInfo) Uptime() time.Duration {
	return time.Since(host.BootTime)
}

type OSInfo struct {
	Platform string `json:"platform"`           // OS platform (e.g. linux, darwin, win32, freebsd).
	Name     string `json:"name"`               // OS Name (e.g. Mac OS X).
	Version  string `json:"version"`            // OS version (e.g. 10.12.6).
	Codename string `json:"codename,omitempty"` // OS codename (e.g. jessie).
	Build    string `json:"build,omitempty"`    // Build (e.g. 16G1114).
}

type LoadAverager interface {
	LoadAverage() LoadAverage
}

type LoadAverage struct {
	One     float64 `json:"one_min"`
	Five    float64 `json:"five_min"`
	Fifteen float64 `json:"fifteen_min"`
}
