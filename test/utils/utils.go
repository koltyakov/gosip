package utils

import (
	"path"
	"runtime"
)

// ResolveCnfgPath : resolves config path
func ResolveCnfgPath(relativePath string) string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "../..", relativePath)
}
