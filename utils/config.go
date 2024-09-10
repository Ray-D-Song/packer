package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func CheckConfigExists() bool {
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error fetching executable path:", err)
		return false
	}

	execDir := filepath.Dir(executablePath)
	configPath := path.Join(execDir, "pack.yml")

	_, err = os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println("Error checking file existence:", err)
		return false
	}
	return true
}
