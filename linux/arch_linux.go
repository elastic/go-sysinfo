package linux

import (
	"syscall"

	"github.com/pkg/errors"
)

func Architecture() (string, error) {
	var uname syscall.Utsname
	if err := syscall.Uname(&uname); err != nil {
		return "", errors.Wrap(err, "architecture")
	}

	data := make([]byte, 0, len(uname.Machine))
	for _, v := range uname.Machine {
		if v == 0 {
			break
		}
		data = append(data, byte(v))
	}

	return string(data), nil
}
