package server

import (
	"fmt"
	"goc/protocol/mom/mqtt"
	"github.com/lumosin/goc/tl/cfgt"
	"github.com/lumosin/goc/tl/errt"
	"github.com/lumosin/goc/tl/jsont"
	"math/rand"
	"time"
)

type mqttHandle struct {
	sendHandle   *mqtt.Handle
	recvHandle   *mqtt.Handle
	reportHanlde *mqtt.Handle
}

func NewMQTT() *mqttHandle {
	rand.Seed(time.Now().Unix())
	cfg := cfgtool.New("conf.json")

	sendJson, err := cfg.TakeJson("MQTTSendCfg")
	errtool.Errpanic(err)
	sendCfg := mqtt.Cfg{}
	sendCfg.ClientID = fmt.Sprintf("%+v", rand.Int())
	err = jsontool.Decode(sendJson, &sendCfg)
	errtool.Errpanic(err)

	recvJson, err := cfg.TakeJson("MQTTRecvCfg")
	errtool.Errpanic(err)
	recvCfg := mqtt.Cfg{}
	recvCfg.ClientID = fmt.Sprintf("%+v", rand.Int())
	err = jsontool.Decode(recvJson, &recvCfg)
	errtool.Errpanic(err)

	reprotJson, err := cfg.TakeJson("MQTTReportCfg")
	errtool.Errpanic(err)
	reportCfg := mqtt.Cfg{}
	reportCfg.ClientID = fmt.Sprintf("%+v", rand.Int())
	err = jsontool.Decode(reprotJson, &reportCfg)
	errtool.Errpanic(err)

	m := &mqttHandle{sendHandle: mqtt.New(sendCfg), recvHandle: mqtt.New(recvCfg), reportHanlde: mqtt.New(reportCfg)}

	log.Info("server mqtt init success")
	return m
}

func (p *mqttHandle) Connect() error {
	errtool.Errpanic(p.sendHandle.Connect())
	errtool.Errpanic(p.recvHandle.Connect())
	errtool.Errpanic(p.reportHanlde.Connect())
	return nil
}

func (p *mqttHandle) DisConnect() error {
	errtool.Errpanic(p.sendHandle.DisConnect())
	errtool.Errpanic(p.recvHandle.DisConnect())
	errtool.Errpanic(p.reportHanlde.DisConnect())
	return nil
}

func (p *mqttHandle) SendRequest(str string) error {
	log.Debug("mqtt send[%s]", str)
	return p.sendHandle.Send(str)
}

func (p *mqttHandle) WaitRequest() string {
	recv := p.recvHandle.Recv()
	log.Debug("mqtt recv[%s]", recv)
	return recv
}

func (p *mqttHandle) Report(str string) error {
	log.Debug("mqtt report[%s]", str)
	return p.reportHanlde.Send(str)
}
