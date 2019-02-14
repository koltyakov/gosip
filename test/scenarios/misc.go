package scenarios

import (
	"path"
	"runtime"
)

func resolveCnfgPath(relativePath string) string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "../..", relativePath)
}
