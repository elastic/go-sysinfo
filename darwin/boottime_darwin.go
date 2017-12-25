package darwin

import (
	"syscall"
	"time"

	"github.com/pkg/errors"
)

const kernBoottimeMIB = "kern.boottime"

func BootTime() (time.Time, error) {
	var tv syscall.Timeval32
	if err := sysctlByName(kernBoottimeMIB, &tv); err != nil {
		return time.Time{}, errors.Wrap(err, "failed to get host uptime")
	}

	bootTime := time.Unix(int64(tv.Sec), int64(tv.Usec)*int64(time.Microsecond))
	return bootTime, nil
}
