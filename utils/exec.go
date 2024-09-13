package utils

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"

	"ray-d-song.com/packer/dict"
)

func Execute(path string) {
	cmd := exec.Command(dict.JSRuntimePath, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	cmd.Wait()
}
