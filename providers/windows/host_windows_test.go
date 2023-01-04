package windows

import (
	"encoding/json"
	"testing"

	"github.com/elastic/go-sysinfo/internal/registry"
)

var _ registry.HostProvider = windowsSystem{}

func TestHost(t *testing.T) {
	host, err := windowsSystem{}.Host()
	if err != nil {
		t.Fatal(err)
	}

	info := host.Info()
	data, _ := json.MarshalIndent(info, "", "  ")
	t.Log(string(data))
}
