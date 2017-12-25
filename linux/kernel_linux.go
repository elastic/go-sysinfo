package linux

import (
	"syscall"

	"github.com/pkg/errors"
)

func KernelVersion() (string, error) {
	var uname syscall.Utsname
	if err := syscall.Uname(&uname); err != nil {
		return "", errors.Wrap(err, "kernel version")
	}

	data := make([]byte, 0, len(uname.Release))
	for _, v := range uname.Release {
		if v == 0 {
			break
		}
		data = append(data, byte(v))
	}

	return string(data), nil
}
