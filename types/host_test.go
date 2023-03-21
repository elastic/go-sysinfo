package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHostInfo_FQDNAwareHostname(t *testing.T) {
	hostInfo := HostInfo{
		Hostname: "foo",
		FQDN: "foo.bar.baz",
	}

	tests := map[string]struct{
		wantFQDN bool
		expected string
	}{
		"want_fqdn": {
			wantFQDN: true,
			expected: "foo.bar.baz",
		},
		"no_fqdn": {
			wantFQDN: false,
			expected: "foo",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := hostInfo.FQDNAwareHostname(test.wantFQDN)
			require.Equal(t, test.expected, actual)
		})
	}
}
