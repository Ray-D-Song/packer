package hooks

import (
	"ray-d-song.com/packer/dict"
	"ray-d-song.com/packer/utils"
)

func AfterSync() {
	if value, exists := dict.HooksMap["after_sync"]; exists {
		utils.Execute(value)
	}
}
