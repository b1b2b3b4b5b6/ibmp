package devser

import (
	"github.com/b1b2b3b4b5b6/goc/logface"
	"github.com/b1b2b3b4b5b6/goc/tl/debt"
)

var DeviceMap map[string]Device
var log = logface.New(logface.DebugLevel)

type DevComdata struct {
	Mac string
	Typ string
	Ver string
}

type DevStatus struct {
	Monitor   bool
	TimeStamp int64
	Online    bool
	StaChange bool
	Init      bool
}

func init() {
	DeviceMap = make(map[string]Device)
	debt.AddFunc("print_devices", print_devices)
}

type Device interface {
	SetComData(comData DevComdata)
	GetComData() *DevComdata
	ReplyCmd(cmd string) (string, string)
	ReplyFeed(interface{}) (string, string)
	GetStatus() *DevStatus
	GetJson() string
}

type DevJson struct {
	Typ       string
	Mac       string
	WriteProp interface{}
}

var bacnetId = 1

func Add(comData DevComdata) Device {
	switch comData.Typ {
	case "UIT":
		uit := &UitDev{ComData: comData}
		uit.BacnetID = bacnetId

		bacnetId++
		DeviceMap[comData.Mac] = uit

		log.Debug("devser add device[%+v]", DeviceMap[comData.Mac])
		return DeviceMap[comData.Mac]
	case "XFS":
		xfs := &XfsDev{ComData: comData}
		xfs.BacnetID = bacnetId

		bacnetId++
		DeviceMap[comData.Mac] = xfs

		log.Debug("devser add device[%+v]", DeviceMap[comData.Mac])
		return DeviceMap[comData.Mac]
	}
	log.Panic("undefined device typ[%+v]", comData)
	return nil
}

func GetTypDev(typ string) []Device {
	resList := make([]Device, 0)
	for _, v := range DeviceMap {
		if v.GetComData().Typ == typ {
			resList = append(resList, v)
		}
	}
	return resList
}

func print_devices() map[string]Device {
	return DeviceMap
}
