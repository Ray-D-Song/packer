//go:build darwin

package lib

import _ "embed"

//go:embed qjs_darwin
var QjsEmbed []byte
