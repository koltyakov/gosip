package utils

import (
	"path"
	"runtime"
)

// ResolveCnfgPath : resolves config path (for auth strategy tests only)
func ResolveCnfgPath(relativePath string) string {
	_, filename, _, _ := runtime.Caller(1)
	// fmt.Println(filename)
	return path.Join(path.Dir(filename), "../..", relativePath)
}
