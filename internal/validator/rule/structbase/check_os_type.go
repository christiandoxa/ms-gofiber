package structbase

import "github.com/go-playground/validator/v10"

var osTypeMap = map[string]struct{}{
	"IOS":     {},
	"ANDROID": {},
	"OTHER":   {},
}

func checkOsType(sl validator.StructLevel, terminalType, osType, osVersion string) {
	if terminalType == "APP" || terminalType == "WAP" {
		if _, ok := osTypeMap[osType]; !ok {
			sl.ReportError(osType, "OsType", "osType", "invalid", "")
		}
		return
	}

	if osType != "" {
		sl.ReportError(osType, "OsType", "osType", "invalid", "")
	}
	if osVersion != "" {
		sl.ReportError(osVersion, "OsVersion", "osVersion", "invalid", "")
	}
}
