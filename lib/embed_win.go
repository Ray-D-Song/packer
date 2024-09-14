//go:build windows

package lib

import _ "embed"

//go:embed jsr_win_amd64.zip
var JsrEmbed []byte
