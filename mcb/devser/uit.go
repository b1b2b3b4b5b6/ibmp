package devser

import (
	"github.com/lumosin/goc/tl/errt"
	"github.com/lumosin/goc/tl/jsont"
)

type UitDev struct {
	Status      DevStatus
	ComData     DevComdata
	BacnetID    int
	ImgAd       []string
	SW_1        int
	SW_2        int
	SW_3        int
	SW_4        int
	Rssi        int
	ImgProgress int
}

func (p *UitDev) SetComData(comData DevComdata) {
	p.ComData = comData
}

func (p *UitDev) GetComData() *DevComdata {
	return &p.ComData
}

func (p *UitDev) ReplyCmd(cmd string) (replyDev string, replySer string) {
	type WP struct {
		ImgAd []string
	}
	var device struct {
		Typ       string
		Mac       string
		WriteProp WP
	}
	err := jsontool.Decode(cmd, &device)
	errtool.Errpanic(err)

	var cusData struct {
		Cmd   string
		ImgAd []string
	}
	cusData.Cmd = "ImgAd"
	cusData.ImgAd = device.WriteProp.ImgAd

	return jsontool.Encode(cusData), ""
}

func (p *UitDev) ReplyFeed(feed interface{}) (replyDev string, replySer string) {
	switch feed.(type) {
	case []byte:
		log.Panic("can not support feed bin")
	case string:
		str := feed.(string)
		log.Trace("device[%+v] feed[%s]", p, str)
		var m struct {
			SW_1        int
			SW_2        int
			SW_3        int
			SW_4        int
			Rssi        int
			Version     string
			ImgProgress int
		}
		jsontool.Decode(str, &m)
		p.Status.Init = true
		p.SW_1 = m.SW_1
		p.SW_2 = m.SW_2
		p.SW_3 = m.SW_3
		p.SW_4 = m.SW_4
		p.Rssi = m.Rssi
		p.ComData.Ver = m.Version
		p.ImgProgress = m.ImgProgress
	}

	p.Status.StaChange = true
	return "", ""
}

func (p *UitDev) GetStatus() *DevStatus {
	return &p.Status
}

func (p *UitDev) GetJson() string {

	type readProp struct {
		ImgProgress int
	}

	type device struct {
		Typ      string
		Mac      string
		Ver      string
		Status   string
		ReadProp readProp
	}
	var m struct {
		Typ    string
		Device device
	}

	m.Typ = "status"
	m.Device.Typ = p.ComData.Typ
	m.Device.Ver = p.ComData.Ver
	m.Device.Mac = p.ComData.Mac

	if p.Status.Online {
		m.Device.Status = "online"
	} else {
		m.Device.Status = "offline"
	}
	m.Device.ReadProp.ImgProgress = p.ImgProgress

	str := jsontool.Encode(&m)
	return str
}
