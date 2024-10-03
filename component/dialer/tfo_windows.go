package dialer

import "github.com/kitty314/1.18.9/constant/features"

func init() {
	// According to MSDN, this option is available since Windows 10, 1607
	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms738596(v=vs.85).aspx
	if features.WindowsMajorVersion < 10 || (features.WindowsMajorVersion == 10 && features.WindowsBuildNumber < 14393) {
		DisableTFO = true
	}
}
