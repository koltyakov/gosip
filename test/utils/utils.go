package utils

import (
	"fmt"
	"path"
	"runtime"
)

// ResolveCnfgPath : resolves config path (for auth startegy tests only)
func ResolveCnfgPath(relativePath string) string {
	_, filename, _, _ := runtime.Caller(1)
	fmt.Println(filename)
	return path.Join(path.Dir(filename), "../..", relativePath)
}
