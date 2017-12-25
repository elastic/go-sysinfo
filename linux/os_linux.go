package linux

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/pkg/errors"

	"github.com/elastic/go-sysinfo/types"
)

const (
	osRelease     = "/etc/os-release"
	redhatRelease = "/etc/redhat-release"
	gentooRelease = "/etc/gentoo-release"
)

func OperatingSystem() (*types.OSInfo, error) {
	f, err := findOSReleaseFile()
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	fields, err := parseOSRelease(content)
	if err != nil {
		return nil, err
	}

	return &types.OSInfo{
		Name:     fields["NAME"],
		Version:  fields["VERSION_ID"],
		Codename: fields["VERSION"],
		Platform: runtime.GOOS,
	}, nil
}

func findOSReleaseFile() (string, error) {
	matches, err := filepath.Glob("/etc/*-release")
	if err != nil {
		return "", errors.Wrap(err, "failed to read /etc")
	}

	for _, name := range matches {
		switch name {
		case osRelease, redhatRelease, gentooRelease:
			return name, nil
		}
	}

	return "", errors.New("os release file not found")
}

func parseOSRelease(content []byte) (map[string]string, error) {
	fields := map[string]string{}

	s := bufio.NewScanner(bytes.NewReader(content))
	for s.Scan() {
		parts := bytes.SplitN(s.Bytes(), []byte("="), 2)
		if len(parts) != 2 {
			continue
		}

		key := string(parts[0])
		value, err := strconv.Unquote(string(parts[1]))
		if err != nil {
			continue
		}

		fields[key] = value
	}

	return fields, s.Err()
}
