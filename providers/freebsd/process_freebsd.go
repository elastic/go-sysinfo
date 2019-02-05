//+build,freebsd,cgo

package freebsd

//#include <sys/sysctl.h>
//#include <sys/time.h>
import "C"

import (
	"strconv"
	"strings"
	"syscall"

	"github.com/pkg/errors"
)

const kernCptimeMIB = "kern.cp_time"
const kernClockrateMIB = "kern.clockrate"

func Cptime() (map[string]uint64, error) {
	var clock clockInfo

	if err := sysctlByName(kernClockrateMIB, &clock); err != nil {
		return make(map[string]uint64), errors.Wrap(err, "failed to get kern.clockrate")
	}

	cptime, err := syscall.Sysctl(kernCptimeMIB)
	if err != nil {
		return make(map[string]uint64), errors.Wrap(err, "failed to get kern.cp_time")
	}

	cpMap := make(map[string]uint64)

	times := strings.Split(cptime, " ")
	names := [5]string{"User", "Nice", "System", "IRQ", "Idle"}

	for index, time := range times {
		i, err := strconv.ParseUint(time, 10, 64)

		if err != nil {
			return cpMap, errors.Wrap(err, "error parsing kern.cp_time")
		}

		cpMap[names[index]] = i * uint64(clock.Tick) * 1000
	}

	return cpMap, nil
}
