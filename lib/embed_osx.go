//go:build darwin

package lib

import _ "embed"

//go:embed jsr_darwin
var JsrEmbed []byte
