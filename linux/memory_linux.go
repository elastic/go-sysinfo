package linux

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func MemTotal() (uint64, error) {
	v, err := findValue("/proc/meminfo", ":", "MemTotal")
	if err != nil {
		return 0, errors.Wrap(err, "failed to get mem total")
	}

	parts := strings.Fields(v)
	if len(parts) != 2 && parts[1] == "kB" {
		return 0, errors.Errorf("failed to parse mem total '%v'", v)
	}

	kB, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to parse mem total '%v'", parts[0])
	}

	return kB * 1024, nil
}
