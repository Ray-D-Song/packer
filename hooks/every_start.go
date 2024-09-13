package hooks

import (
	"bytes"
	"io"
	"os"
	"runtime"

	"ray-d-song.com/packer/dict"
	"ray-d-song.com/packer/lib"
	"ray-d-song.com/packer/utils"
)

func EveryStart() {
	ensureJSRuntime()
}

func ensureJSRuntime() {
	if !utils.CheckExists(dict.JSRuntimePath) {
		dst, err := os.CreateTemp("", "jsr-*.zip")
		if err != nil {
			panic(err)
		}
		defer os.Remove(dst.Name())

		if _, err := io.Copy(dst, bytes.NewReader(lib.JsrEmbed)); err != nil {
			panic(err)
		}
		utils.Unzip(dst.Name(), dict.LibsDir)

		if runtime.GOOS != "windows" {
			os.Chmod(dict.JSRuntimePath, 0755)
		}
	}
}
