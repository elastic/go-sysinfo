package linux

import "github.com/elastic/go-sysinfo/internal/registry"

var _ registry.HostProvider = linuxSystem{}
var _ registry.ProcessProvider = linuxSystem{}
