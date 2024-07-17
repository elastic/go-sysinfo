package linux

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMachineIDLookup(t *testing.T) {
	path := "testdata/fedora30"
	known := "144d62edb0f142458f320852f495b72c"
	id, err := MachineIDHostfs(path)
	require.NoError(t, err)
	require.Equal(t, known, id)
}
