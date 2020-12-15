package mcb

import (
	"github.com/b1b2b3b4b5b6/goc/tl/errt"
	"github.com/b1b2b3b4b5b6/goc/tl/jsont"
	"ibmp/mcb/devser"
	"ibmp/mcb/report"
	"ibmp/mcb/server"
	"ibmp/wifimesh"
)

var ser = server.New()

func init() {
	go cmdHandleLoop()
	go feedHandleLoop()
}

func reportOnlineDevices(json string) {
	type device struct {
		Typ string
		Mac string
	}

	type initMsg struct {
		Typ     string
		Devices []device
	}
	im := initMsg{Typ: "init", Devices: make([]device, 0)}
	readDevices := report.GetOnineDevices()
	for _, v := range readDevices {
		im.Devices = append(im.Devices, device{Typ: v.GetComData().Typ, Mac: v.GetComData().Mac})
	}
	msg := jsont.Encode(&im)
	server.New().SendRequest(msg)
}

func controlDeviceInit(json string) {
	var m struct {
		Typ     string
		Devices []devser.DevComdata
	}
	err := jsont.Decode(json, &m)
	errt.Errpanic(err)

	for _, v := range devser.DeviceMap {
		v.GetStatus().Monitor = false
	}

	for _, initDev := range m.Devices {
		realDev, ok := devser.DeviceMap[initDev.Mac]
		if ok {
			realDev.GetStatus().Monitor = true
		} else {
			addVal := devser.Add(initDev)
			addVal.GetStatus().Monitor = true
		}
	}
}

func parseDeviceCmd(json string) {
	sendMap := make(map[string]string)
	var m struct {
		Typ     string
		Devices []devser.DevJson
	}
	err := jsont.Decode(json, &m)
	errt.Errpanic(err)
	for _, v := range m.Devices {
		dev, ok := devser.DeviceMap[v.Mac]
		if ok {
			cmd, _ := dev.ReplyCmd(jsont.Encode(&v))
			sendMap[v.Mac] = cmd
		} else {
			log.Error("unkonw device[%+v]", v)
		}
	}
	sort2dev(sendMap)
}

func feedHandleLoop() {
	for {
		recv := wifimesh.GetGroup().Recv()
		dev, ok := devser.DeviceMap[recv.Mac]
		if ok {
			dev.ReplyFeed(recv.Data)
			continue
		} else {
			devser.Add(devser.DevComdata{Mac: recv.Mac, Typ: recv.Typ, Ver: recv.Ver}).ReplyFeed(recv.Data)
		}
	}
}

func saveVersion(json string) {
	
}

func cmdHandleLoop() {
	for {
		recv := server.New().WaitRequest()

		var m struct {
			Typ string
		}
		jsont.Decode(recv, &m)
		switch m.Typ {
		case "askDevices":
			reportOnlineDevices(recv)

		case "init":
			controlDeviceInit(recv)

		case "status":
			parseDeviceCmd(recv)

		case "version":

		default:
			log.Panic("illegal cmd[%+v]", m)
		}
	}
}
