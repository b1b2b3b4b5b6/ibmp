package report

import (
	"github.com/b1b2b3b4b5b6/goc/logface"
	"github.com/b1b2b3b4b5b6/goc/tl/cfgt"
	"ibmp/mcb/devser"
	"ibmp/mcb/server"
	"time"
)

var log = logface.New(logface.InfoLevel)
var cfg = cfgt.New("conf.json")

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
