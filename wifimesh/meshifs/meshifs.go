package meshifs

/*
keepalive test:
	client break: pass
	server break: pass
recv&send test: pass
*/

import (
	"goc/logface"
	"goc/toolcom/cfgtool"
)

var log = logface.New(logface.InfoLevel)
var cfg = cfgtool.New("conf.json")

type RecvRaw struct {
	MeshID string
	Typ    uint8
	Data   []byte
}

const (
	MaxifsLen = 1472
)

var chRecvRaw = make(chan *RecvRaw, 1000)
var chMeshCon = make(chan MeshCon, 1000)

type MeshCon interface {
	Send([]byte) error
	GetMeshID() string
	Destroy()
}

func RawRecv() *RecvRaw {
	return <-chRecvRaw
}

func TopoRecv() MeshCon {
	return <-chMeshCon
}
