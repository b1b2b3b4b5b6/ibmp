package devser

import (
	"github.com/lumosin/goc/tl/errt"
	"github.com/lumosin/goc/tl/jsont"
)

type XfsDev struct {
	Status   DevStatus
	ComData  DevComdata
	BacnetID int
	IP       string
}

func (p *XfsDev) SetComData(comData DevComdata) {
	p.ComData = comData
}

func (p *XfsDev) GetComData() *DevComdata {
	return &p.ComData
}

func (p *XfsDev) ReplyCmd(cmd string) (replyDev string, replySer string) {
	var device struct {
		Typ string
		Mac string
	}
	err := jsontool.Decode(cmd, &device)
	errtool.Errpanic(err)

	var cusData struct {
		Cmd string
	}

	return jsontool.Encode(cusData), ""
}

func (p *XfsDev) ReplyFeed(feed interface{}) (replyDev string, replySer string) {
	switch feed.(type) {
	case []byte:
		log.Panic("can not support feed bin")
	case string:
		str := feed.(string)
		log.Trace("device[%+v] feed[%s]", p, str)
		var m struct {
			Version string
			IP      string
		}
		jsontool.Decode(str, &m)
		p.Status.Init = true
		p.ComData.Ver = m.Version
		p.IP = m.IP
	}

	p.Status.StaChange = true
	return "", ""
}

func (p *XfsDev) GetStatus() *DevStatus {
	return &p.Status
}

func (p *XfsDev) GetJson() string {

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

	str := jsontool.Encode(&m)
	return str
}
