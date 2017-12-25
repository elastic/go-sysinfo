package darwin

import (
	"syscall"

	"github.com/pkg/errors"
)

const hardwareMIB = "hw.machine"

func Architecture() (string, error) {
	arch, err := syscall.Sysctl(hardwareMIB)
	if err != nil {
		return "", errors.Wrap(err, "failed to get architecture")
	}

	return arch, nil
}
