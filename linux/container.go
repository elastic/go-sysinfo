package linux

import (
	"bufio"
	"bytes"
	"io/ioutil"

	"github.com/pkg/errors"
)

const procOneCgroup = "/proc/1/cgroup"

// IsContainerized returns true if this process is containerized.
func IsContainerized() (bool, error) {
	data, err := ioutil.ReadFile(procOneCgroup)
	if err != nil {
		return false, errors.Wrap(err, "failed to read process cgroups")
	}

	return isContainerizedCgroup(data)
}

func isContainerizedCgroup(data []byte) (bool, error) {
	s := bufio.NewScanner(bytes.NewReader(data))
	for n := 0; s.Scan(); n++ {
		line := s.Bytes()
		if len(line) == 0 || line[len(line)-1] == '/' {
			continue
		}

		if bytes.HasSuffix(line, []byte("init.scope")) {
			return false, nil
		}
	}

	return true, s.Err()
}
