package darwin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/go-sysinfo/internal/registry"
)

var _ registry.HostProvider = darwinSystem{}
var _ registry.ProcessProvider = darwinSystem{}

func TestKernProcInfo(t *testing.T) {
	var p process
	if err := kern_procargs(os.Getpid(), &p); err != nil {
		t.Fatal(err)
	}

	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, exe, p.exe)
	assert.Equal(t, os.Args, p.args)
}
