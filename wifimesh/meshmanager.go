package wifimesh

import (
	"github.com/b1b2b3b4b5b6/goc/tl/turnt"
	"github.com/b1b2b3b4b5b6/goc/tl/errt"
	"github.com/b1b2b3b4b5b6/goc/tl/jsont"
	"ibmp/wifimesh/meshifs"
	"io/ioutil"
)

var initMeshID = "00:00:00:00:00:00"

type foundMesh struct {
	conn meshifs.MeshCon
}

type deviceRecv struct {
	Typ  string
	Mac  string
	Ver  string
	Data interface{}
}

type devicesSend struct {
	macs []string
	data string
}

var chFoundMesh = make(chan *foundMesh, 1000)

var chMeshDataW = make(chan *devicesSend)
var chMeshDataR = make(chan *deviceRecv, 1000)
var chHeatBeat = make(chan []string, 1000)

func init() {
	go meshHandleLoop()
	go meshDataRLoop()
	go meshRenewLoop()
}

func foundMeshHandle(m *foundMesh) {
	mg.Add(m.conn)
}

func meshDataWHandle(m *devicesSend) {
	for _, mesh := range mg.GetMeshMap() {
		tar_macs := mesh.GetOwnMac(m.macs)
		if len(tar_macs) != 0 {
			go mesh.SendCustom(tar_macs, m.data)
		}
	}
}

func meshHandleLoop() {
	for {
		select {
		case m := <-chFoundMesh:
			foundMeshHandle(m)

		case m := <-chMeshDataW:
			meshDataWHandle(m)
		}
	}
}

func parseBin(raw *meshifs.RecvRaw) {
	log.Debug("recv bin[%d]", len(raw.Data))
	bin := &deviceRecv{}
	bin.Mac = turnt.Mac2Str(raw.Data[0:6])
	bin.Data = raw.Data[6:]
	chMeshDataR <- bin
}

func parseStr(raw *meshifs.RecvRaw) {
	if raw.Data[len(raw.Data)-1] == 0 {
		raw.Data = raw.Data[:len(raw.Data)-1] //delete 0 from tail of string
	}
	//log.Trace("recv str: %s", string(raw.Data))
	var m struct {
		Typ string
	}
	err := jsont.Decode(string(raw.Data), &m)
	errt.Errpanic(err)
	switch m.Typ {
	case "log":
		log.Logkit(string(raw.Data))

	case "brust":
		var m struct {
			Typ       string
			Mac       string
			ParentMac string
			Layer     int
			Version   string
			DeviceTyp string
		}
		jsont.Decode(string(raw.Data), &m)
		mesh := mg.GetMesh(raw.MeshID)
		if mesh == nil {
			return
		}
		recv := &deviceRecv{Mac: m.Mac, Ver: m.Version, Typ: m.DeviceTyp}
		recv.Data = string(raw.Data)

		chMeshDataR <- recv

	case "online":
		var m struct {
			Typ     string
			Devices []string
		}
		jsont.Decode(string(raw.Data), &m)
		mesh := mg.GetMesh(raw.MeshID)
		mesh.RefreshDevices(m.Devices)
		chHeatBeat <- m.Devices

	case "num":
		var m struct {
			Typ string
			Num int
		}
		jsont.Decode(string(raw.Data), &m)
		mesh := mg.GetMesh(raw.MeshID)
		mesh.num = m.Num

	case "cut_ota":
		cut_size := 1024
		var m struct {
			Typ string
			Seq int
		}
		buf, err := ioutil.ReadFile("D://Code/ESP32/mcb/build/" + "Lighting")
		errt.Errpanic(err)
		jsont.Decode(string(raw.Data), &m)
		seq_map := make(map[int][]byte)

		n := 0
		for {
			if len(buf) > cut_size {
				seq_map[n] = buf[:cut_size]
				buf = buf[cut_size:]
			} else {
				seq_map[n] = buf[0:]
				break
			}
			n++
		}
		mesh_target := mg.GetMesh(raw.MeshID)

		if seq_map[m.Seq] == nil {
			bin := make([]byte, 0)
			mesh_target.SendBin(bin)
		} else {
			bin := make([]byte, 0)
			bin = append(bin, seq_map[m.Seq]...)
			mesh_target.SendBin(bin)
		}
	default:
		log.Panic("undefined typ")
	}
}

const (
	OutBin = 0
	OutStr = 1
)

func meshDataRLoop() {
	for {
		raw := meshifs.RawRecv()
		switch raw.Typ {
		case OutBin:
			parseBin(raw)
		case OutStr:
			parseStr(raw)
		default:
			log.Panic("receive undefine raw typ")
		}
	}
}

func meshRenewLoop() {
	for {
		mc := meshifs.TopoRecv()
		chFoundMesh <- &foundMesh{conn: mc}
	}
}
