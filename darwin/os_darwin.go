package darwin

import (
	"io/ioutil"
	"runtime"

	"github.com/pkg/errors"
	"howett.net/plist"

	"github.com/elastic/go-sysinfo/types"
)

const (
	systemVersionPlist = "/System/Library/CoreServices/SystemVersion.plist"

	plistProductName         = "ProductName"
	plistProductVersion      = "ProductVersion"
	plistProductBuildVersion = "ProductBuildVersion"
)

func OperatingSystem() (*types.OSInfo, error) {
	data, err := ioutil.ReadFile(systemVersionPlist)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read plist file")
	}

	attrs := map[string]string{}
	if _, err = plist.Unmarshal(data, &attrs); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal plist data")
	}

	productName, found := attrs[plistProductName]
	if !found {
		return nil, errors.Wrapf(err, "plist key %v not found", plistProductName)
	}

	version, found := attrs[plistProductVersion]
	if !found {
		return nil, errors.Wrapf(err, "plist key %v not found", plistProductVersion)
	}

	build, found := attrs[plistProductBuildVersion]
	if !found {
		return nil, errors.Wrapf(err, "plist key %v not found", plistProductBuildVersion)
	}

	return &types.OSInfo{
		Platform: runtime.GOOS,
		Name:     productName,
		Version:  version,
		Build:    build,
	}, nil
}
