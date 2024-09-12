//go:build windows

package lib

import _ "embed"

//go:embed jsr_win_x64.exe
var JsrEmbed []byte
