package linux

import (
	"testing"

	"github.com/elastic/go-sysinfo/types"
	"github.com/stretchr/testify/assert"
)

func TestOperatingSystem(t *testing.T) {
	t.Run("centos6", func(t *testing.T) {
		os, err := getOSInfo("testdata/centos6")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Family:   "redhat",
			Platform: "centos",
			Name:     "CentOS",
			Version:  "6.9 (Final)",
			Major:    6,
			Minor:    9,
			Codename: "Final",
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("centos7", func(t *testing.T) {
		os, err := getOSInfo("testdata/centos7")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Family:   "redhat",
			Platform: "centos",
			Name:     "CentOS Linux",
			Version:  "7 (Core)",
			Major:    7,
			Minor:    4,
			Patch:    1708,
			Codename: "Core",
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("debian9", func(t *testing.T) {
		os, err := getOSInfo("testdata/debian9")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Family:   "debian",
			Platform: "debian",
			Name:     "Debian GNU/Linux",
			Version:  "9 (stretch)",
			Major:    9,
			Codename: "stretch",
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("ubuntu1404", func(t *testing.T) {
		os, err := getOSInfo("testdata/ubuntu1404")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Family:   "debian",
			Platform: "ubuntu",
			Name:     "Ubuntu",
			Version:  "14.04.5 LTS, Trusty Tahr",
			Major:    14,
			Minor:    4,
			Patch:    5,
			Codename: "trusty",
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("ubuntu1710", func(t *testing.T) {
		os, err := getOSInfo("testdata/ubuntu1710")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Family:   "debian",
			Platform: "ubuntu",
			Name:     "Ubuntu",
			Version:  "17.10 (Artful Aardvark)",
			Major:    17,
			Minor:    10,
			Patch:    0,
			Codename: "artful",
		}, *os)
		t.Logf("%#v", os)
	})
}
