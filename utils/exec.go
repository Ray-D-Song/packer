package utils

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"

	"ray-d-song.com/packer/lib"
)

func Execute() {
	tmpFile, err := os.CreateTemp("", "njs")
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		return
	}
	defer os.Remove(tmpFile.Name())
	if _, err := io.Copy(tmpFile, bytes.NewReader(lib.QjsEmbed)); err != nil {
		fmt.Println("Error writing to temp file:", err)
		return
	}
	tmpFile.Close()
	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		fmt.Println("Error setting file permissions:", err)
		return
	}

	cmd := exec.Command(tmpFile.Name(), "./log.js")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for command to finish:", err)
		return
	}
}
