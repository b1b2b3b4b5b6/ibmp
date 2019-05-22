package mcb

import (
	"ibmp/wifimesh"
)

func sort2dev(sendMap map[string]string) {
	macList := make([]string, 1)
	for k, v := range sendMap {
		macList[0] = k
		wifimesh.GetGroup().Send(macList, v)
	}
}
