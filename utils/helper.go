package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func GetExecPath() (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error fetching executable path:", err)
		return "", err
	}

	return filepath.Dir(executablePath), nil
}

func GetLibsPath() string {
	execPath, err := GetExecPath()
	if err != nil {
		panic(err)
	}
	return path.Join(execPath, "pack_libs")
}
