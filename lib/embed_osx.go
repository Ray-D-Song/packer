//go:build darwin

package lib

import _ "embed"

//go:embed qjs_darwin_aarch64
var QjsEmbed []byte
