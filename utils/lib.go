package utils

import (
	"os"
	"os/user"
	"path"
	"path/filepath"
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

func Diff(a, b []string) []string {
	bMap := make(map[string]bool)
	for _, item := range b {
		bMap[item] = true
	}

	var result []string
	for _, item := range a {
		if !bMap[item] {
			result = append(result, item)
		}
	}

	return result
}

func DiffDeps(deps []string) []string {
	execPath, err := GetExecPath()
	if err != nil {
		panic(err)
	}
	libsPath := path.Join(execPath, "pack_libs")
	if !CheckExists(libsPath) {
		EnsureDirExists(libsPath)
		return deps
	}
	subs, err := GetSubdirectoryNames(libsPath)
	if err != nil {
		panic(err)
	}

	return Diff(deps, subs)
}

func GetSubdirectoryNames(dir string) ([]string, error) {
	var subdirNames []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != dir {
			subdirNames = append(subdirNames, filepath.Base(path))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return subdirNames, nil
}
