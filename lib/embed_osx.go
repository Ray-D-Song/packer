//go:build darwin

package lib

import _ "embed"

//go:embed jsr_darwin_aarch64.zip
var JsrEmbed []byte
