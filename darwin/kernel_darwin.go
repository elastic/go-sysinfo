package darwin

import (
	"syscall"

	"github.com/pkg/errors"
)

const kernelReleaseMIB = "kern.osrelease"

func KernelVersion() (string, error) {
	version, err := syscall.Sysctl(kernelReleaseMIB)
	if err != nil {
		return "", errors.Wrap(err, "failed to get kernel version")
	}

	return version, nil
}
