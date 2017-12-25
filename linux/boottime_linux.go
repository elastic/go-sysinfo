package linux

import (
	"sync"
	"time"

	"github.com/prometheus/procfs"
)

var (
	bootTime     time.Time
	bootTimeLock sync.Mutex
)

func BootTime() (time.Time, error) {
	bootTimeLock.Lock()
	defer bootTimeLock.Unlock()

	if !bootTime.IsZero() {
		return bootTime, nil
	}

	stat, err := procfs.NewStat()
	if err != nil {
		return time.Time{}, nil
	}

	bootTime = time.Unix(int64(stat.BootTime), 0)
	return bootTime, nil
}
