package ca

import (
	"github.com/kitty314/1.18.9/constant/features"
)

func init() {
	// crypto/x509: certificate validation in Windows fails to validate IP in SAN
	// https://github.com/golang/go/issues/37176
	// As far as I can tell this is still the case on most older versions of Windows (but seems to be fixed in 10)
	if features.WindowsMajorVersion < 10 && len(_CaCertificates) > 0 {
		DisableSystemCa = true
	}
}
