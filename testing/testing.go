package testing

import (
	"os"
	"path"
	"runtime"
)

// to correctly load config files, keys etc.
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
