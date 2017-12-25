package linux

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"

	"github.com/elastic/go-sysinfo/types"
)

func MachineID() (string, error) {
	id, err := ioutil.ReadFile("/etc/machine-id")
	if os.IsNotExist(err) {
		return "", types.ErrNotImplemented
	}
	id = bytes.TrimSpace(id)
	return string(id), errors.Wrap(err, "failed to read machine-id")
}
