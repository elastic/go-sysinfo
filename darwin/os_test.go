package darwin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const SystemVersionPlist = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
        <key>ProductBuildVersion</key>
        <string>16G1114</string>
        <key>ProductCopyright</key>
        <string>1983-2017 Apple Inc.</string>
        <key>ProductName</key>
        <string>Mac OS X</string>
        <key>ProductUserVisibleVersion</key>
        <string>10.12.6</string>
        <key>ProductVersion</key>
        <string>10.12.6</string>
</dict>
</plist>
`

func TestOperatingSystem(t *testing.T) {
	osInfo, err := getOSInfo([]byte(SystemVersionPlist))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "darwin", osInfo.Family)
	assert.Equal(t, "darwin", osInfo.Platform)
	assert.Equal(t, "Mac OS X", osInfo.Name)
	assert.Equal(t, "10.12.6", osInfo.Version)
	assert.Equal(t, 10, osInfo.Major)
	assert.Equal(t, 12, osInfo.Minor)
	assert.Equal(t, 6, osInfo.Patch)
	assert.Equal(t, "16G1114", osInfo.Build)
}
