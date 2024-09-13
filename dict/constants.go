package dict

import (
	"os/user"
	"path"
	"runtime"
)

var usr *user.User
var JSRuntimePath string
var UsrDir string
var PackerDir string
var StorageDir string
var LibsDir string

func init() {
	var err error
	usr, err = user.Current()
	if err != nil {
		panic(err)
	}

	UsrDir = usr.HomeDir
	PackerDir = path.Join(UsrDir, ".packer")
	StorageDir = path.Join(PackerDir, "storage")
	LibsDir = path.Join(PackerDir, "jsr")

	if runtime.GOOS == "windows" {
		JSRuntimePath = path.Join(PackerDir, "jsr", "jsr.exe")
	} else {
		JSRuntimePath = path.Join(PackerDir, "jsr", "jsr")
	}
}
