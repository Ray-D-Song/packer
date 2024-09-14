//go:build linux

package lib

import _ "embed"

//go:embed jsr_linux_amd64.zip
var JsrEmbed []byte
