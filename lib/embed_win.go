//go:build windows

package lib

import _ "embed"

//go:embed jsr.exe.zip
var JsrEmbed []byte
