package freebsd

import (
	"github.com/elastic/go-sysinfo/internal/registry"
)

var _ registry.HostProvider = freebsdSystem{}
var _ registry.ProcessProvider = freebsdSystem{}
