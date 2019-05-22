package report

import (
	"github.com/lumosin/goc/tl/errt"
	"ibmp/mcb/devser"
	"ibmp/wifimesh"
	"time"
)

func init() {
	go onlineSetLoop()
	go onlineUnsetLoop()
}

func GetOnineDevices() map[string]devser.Device {

	devices := make(map[string]devser.Device, 0)
	for _, v := range devser.DeviceMap {
		if v.GetStatus().Online {
			devices[v.GetComData().Mac] = v
		}
	}
	return devices
}

func onlineSetLoop() {
	for {
		devices := wifimesh.GetGroup().RecvDevsHB()
		for _, mac := range devices {
			dev, ok := devser.DeviceMap[mac]
			if ok {

				lastOnline := dev.GetStatus().Online
				dev.GetStatus().Online = true
				dev.GetStatus().TimeStamp = time.Now().Unix()
				log.Debug("%+v", dev)
				if lastOnline == false {
					dev.GetStatus().StaChange = true
				}
			}
		}

	}
}

func onlineUnsetLoop() {
	hbTimeS, err := cfg.TakeInt("HBTimeS")
	log.Info("set device HB time [%d]s", hbTimeS)
	errt.Errpanic(err)
	for {
		for _, v := range devser.DeviceMap {
			timeMinus := time.Now().Unix() - v.GetStatus().TimeStamp
			if timeMinus > int64(hbTimeS) && v.GetStatus().Online == true {
				v.GetStatus().Online = false
				log.Debug("offline[%d] dev[%+v]", timeMinus, v)
				if v.GetStatus().Monitor {
					v.GetStatus().StaChange = true
				}

			}
		}
		time.Sleep(time.Second * 1)
	}
}
