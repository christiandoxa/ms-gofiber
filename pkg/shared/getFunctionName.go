package shared

import (
	"runtime"
	"strings"
)

func GetFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	parts := strings.Split(fn.Name(), ".")
	return parts[len(parts)-1]
}
