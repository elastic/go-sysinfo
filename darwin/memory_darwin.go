package darwin

import (
	"github.com/pkg/errors"
)

const hwMemsizeMIB = "hw.memsize"

func MemTotal() (uint64, error) {
	var size uint64
	if err := sysctlByName(hwMemsizeMIB, &size); err != nil {
		return 0, errors.Wrap(err, "failed to get mem total")
	}

	return size, nil
}
