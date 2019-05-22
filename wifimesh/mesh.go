package wifimesh

import (
	"encoding/binary"
	"encoding/json"
	"github.com/lumosin/goc/tl/errt"
	"github.com/lumosin/goc/tl/jsont"
	"ibmp/wifimesh/meshifs"
)

const (
	sendBin = 0
	sendStr = 1
)

type mesh struct {
	meshID     string
	deviceList []string
	chSend     chan []byte
	conn       meshifs.MeshCon
	num        int
}

func newMesh(conn meshifs.MeshCon) *mesh {
	mesh := &mesh{meshID: conn.GetMeshID(), conn: conn}
	mesh.chSend = make(chan []byte)
	go mesh.sendDataLopp()
	return mesh
}

func (p *mesh) ReportGather() {
	p.chSend <- createStrW("ReportGather", nil, "")
}

func (p *mesh) ReportNum() {
	p.chSend <- createStrW("ReportNum", nil, "")
}

func (p *mesh) Ota(macs []string, typ string, peroidMs int, isHttp int) {
	var m struct {
		AppUrl   string
		PeroidMs int
		IsHttp   int
	}
	m.AppUrl = getBinUrl(typ)
	if m.AppUrl == "" {
		log.Error("bin[%s] not exist", typ)
		return
	}
	m.PeroidMs = peroidMs
	m.IsHttp = isHttp
	p.chSend <- createStrW("Ota", macs, string(jsont.Encode(m)))
	log.Debug("Ota[%s][%d] macs: %+v", typ, peroidMs, macs)
}

func (p *mesh) Add(macs []string, targeMac string) {
	var m struct {
		Mac string
	}
	p.chSend <- createStrW("Add", macs, string(jsont.Encode(m)))
	log.Debug("Add macs: %+v", macs)
}

func (p *mesh) Delete(macs []string) {
	p.chSend <- createStrW("Delete", macs, "")
	log.Debug("Delete macs: %+v", macs)
}

func (p *mesh) ChangeConfig(macs []string, MESHID string, SSID string, PASSWD string) {
	var m struct {
		SSID   string
		PASSWD string
		MESHID string
	}
	m.SSID = SSID
	m.PASSWD = PASSWD
	m.MESHID = MESHID
	p.chSend <- createStrW("ChangeConfig", macs, string(jsont.Encode(m)))
	log.Debug("ChangeConfig[SSID: %s, PASSWS: %s MeshID: %s] macs: %+v", SSID, PASSWD, MESHID, macs)
}

func (p *mesh) SendCustom(macs []string, argCusData interface{}) {
	p.chSend <- createStrW("SendCustom", macs, argCusData)
	log.Debug("SendCustom[%s] macs: %+v", argCusData, macs)
}

func (p *mesh) SendBin(bin []byte) {
	p.chSend <- createBinW(bin)
	log.Debug("SendBin[%d]", len(bin))
}

func (p *mesh) Destroy() {
	close(p.chSend)
	p.Delete(nil)
	p.conn.Destroy()
}

//following func just for code convenient

func (p *mesh) GetOwnMac(macs []string) []string {
	own_macs := make([]string, 0, 100)
	for _, mac := range macs {
		for _, val := range p.deviceList {
			if val == mac {
				own_macs = append(own_macs, val)
			}
		}
	}
	return own_macs
}

func createStrW(typ string, macs []string, cusDateW interface{}) []byte {

	var data struct {
		Typ     string
		Devcies []string
		CusData interface{}
	}
	data.Typ = typ
	data.Devcies = macs
	data.CusData = cusDateW
	jsonByte, err := json.Marshal(data)
	jsonByte = append(jsonByte, 0) // add 0 to the tail of string
	errt.Errpanic(err)

	strByte := append([]byte{0, 0, 0, 0, 0, 0, sendStr}, jsonByte...)
	lenByte := make([]byte, 2)
	binary.BigEndian.PutUint16(lenByte, uint16(len(strByte)))
	strByte = append(lenByte, strByte...)
	return strByte
}

func createBinW(bin []byte) []byte {
	binByte := append([]byte{0, 0, 0, 0, 0, 0, sendBin}, bin...)
	len_byte := make([]byte, 2)
	binary.BigEndian.PutUint16(len_byte, uint16(len(binByte)))
	binByte = append(len_byte, binByte...)
	return binByte
}

func (p *mesh) sendDataLopp() {
	for {
		data := <-p.chSend
		p.conn.Send(data)
	}
}

func (p *mesh) RefreshDevices(macs []string) {
	p.deviceList = macs
}
