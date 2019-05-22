package report

import (
	"github.com/lumosin/goc/logface"
	"github.com/lumosin/goc/tl/cfgt"
	"ibmp/mcb/devser"
	"ibmp/mcb/server"
	"time"
)

var log = logface.New(logface.InfoLevel)
var cfg = cfgtool.New("conf.json")

func init() {
	go monitorLoop()
}

func monitorLoop() {
	for {
		for _, v := range devser.DeviceMap {
			sta := v.GetStatus()
			if sta.Monitor && sta.StaChange {
				server.New().Report(v.GetJson())
				sta.StaChange = false
			}
		}
		time.Sleep(time.Second * 1)
	}

}
