package report

import (
	"github.com/b1b2b3b4b5b6/goc/tl/jsont"
	"ibmp/mcb/devser"
	"ibmp/wifimesh"
	"time"
)

func init() {
	go getInfoLoop()
}

func getInfoLoop() {
	for {
		for _, v := range devser.DeviceMap {
			sta := v.GetStatus()
			if sta.Online && !sta.Init {
				var cusData struct {
					Cmd string
				}
				cusData.Cmd = "Brust"
				wifimesh.GetGroup().Send([]string{v.GetComData().Mac}, jsont.Encode(&cusData))
				time.Sleep(time.Second * 1)
				continue
			}
		}
		time.Sleep(time.Second * 5)
	}
}
