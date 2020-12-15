package wifimesh

import (
	"github.com/b1b2b3b4b5b6/goc/logface"
	"github.com/b1b2b3b4b5b6/goc/tl/cfgt"
	"ibmp/wifimesh/meshifs"
)

type meshGroupType struct {
	ver string
	typ string
}

//MeshGroup is
type MeshGroup struct {
	typ     meshGroupType
	meshMap map[string]*mesh
}

var mg = &MeshGroup{meshMap: map[string]*mesh{}}
var log = logface.New(logface.TraceLevel)
var cfg = cfgt.New("conf.json")

//New is
func GetGroup() *MeshGroup {
	return mg
}

func (p *MeshGroup) Add(conn meshifs.MeshCon) {
	m := p.meshMap[conn.GetMeshID()]
	if m != nil {
		log.Trace("renew mesh[%s]", m.meshID)
		m.conn = conn
	} else {
		m = newMesh(conn)
		mg.meshMap[m.meshID] = m
		log.Debug("create new mesh[%s]", m.meshID)
	}
}

func (p *MeshGroup) Delete(meshID string) {
	m := p.meshMap[meshID]
	if m != nil {
		m.Delete(nil)
		log.Info("delete mesh[%s]", m.meshID)
	}
}

func (p *MeshGroup) GetMeshMap() map[string]*mesh {
	return p.meshMap
}

func (p *MeshGroup) GetMesh(meshID string) *mesh {
	return p.meshMap[meshID]
}

func (p *MeshGroup) Recv() *deviceRecv {
	return <-chMeshDataR
}

func (p *MeshGroup) Send(macs []string, data string) {
	device_send := &devicesSend{macs: macs, data: data}
	chMeshDataW <- device_send
}

func (p *MeshGroup) RecvDevsHB() []string {
	return <-chHeatBeat
}
