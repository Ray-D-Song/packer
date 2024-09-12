//go:build windows

package lib

import _ "embed"

//go:embed qjs_win_x64.exe
var QjsEmbed []byte
