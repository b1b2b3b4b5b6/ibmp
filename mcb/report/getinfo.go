package report

import (
	"goc/toolcom/jsontool"
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
				wifimesh.GetGroup().Send([]string{v.GetComData().Mac}, jsontool.Encode(&cusData))
				time.Sleep(time.Second * 1)
				continue
			}
		}
		time.Sleep(time.Second * 5)
	}
}
