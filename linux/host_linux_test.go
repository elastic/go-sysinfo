package linux

import (
	"encoding/json"
	"testing"

	"github.com/elastic/go-sysinfo/internal/registry"
)

var _ registry.HostProvider = linuxSystem{}

func TestHost(t *testing.T) {
	host, err := linuxSystem{}.Host()
	if err != nil {
		t.Fatal(err)
	}

	info := host.Info()
	data, _ := json.MarshalIndent(info, "", "  ")
	t.Log(string(data))
}
