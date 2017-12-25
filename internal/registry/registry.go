package registry

import (
	"github.com/pkg/errors"

	"github.com/elastic/go-sysinfo/types"
)

var (
	hostProvider    HostProvider
	processProvider ProcessProvider
)

type HostProvider interface {
	Host() (types.Host, error)
}

type ProcessProvider interface {
	Process(pid int) (types.Process, error)
	Processes() ([]types.Process, error)
}

func Register(provider interface{}) {
	if h, ok := provider.(HostProvider); ok {
		if hostProvider != nil {
			panic(errors.Errorf("HostProvider already registered: %v", hostProvider))
		}
		hostProvider = h
	}

	if p, ok := provider.(ProcessProvider); ok {
		if processProvider != nil {
			panic(errors.Errorf("ProcessProvider already registered: %v", processProvider))
		}
		processProvider = p
	}
}

func GetHostProvider() HostProvider       { return hostProvider }
func GetProcessProvider() ProcessProvider { return processProvider }
