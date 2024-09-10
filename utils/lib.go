package utils

import (
	"os"
	"os/user"
	"path"
)

var StorageDir string

func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	StorageDir = path.Join(usr.HomeDir, ".packer", "storage")
}

func CheckExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func CheckLibExist(libName string) bool {
	libPath := path.Join(StorageDir, libName)

	return CheckExists(libPath)
}

func CheckLibVersionExist(libName string, version string) bool {
	libPath := path.Join(StorageDir, libName, version)

	return CheckExists(libPath)
}
