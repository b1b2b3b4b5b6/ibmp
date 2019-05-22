package server

import (
	"fmt"
	"github.com/lumosin/goc/pt/mom/mqtt"
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
	cfg := cfgt.New("conf.json")

	sendJson, err := cfg.TakeJson("MQTTSendCfg")
	errt.Errpanic(err)
	sendCfg := mqtt.Cfg{}
	sendCfg.ClientID = fmt.Sprintf("%+v", rand.Int())
	err = jsont.Decode(sendJson, &sendCfg)
	errt.Errpanic(err)

	recvJson, err := cfg.TakeJson("MQTTRecvCfg")
	errt.Errpanic(err)
	recvCfg := mqtt.Cfg{}
	recvCfg.ClientID = fmt.Sprintf("%+v", rand.Int())
	err = jsont.Decode(recvJson, &recvCfg)
	errt.Errpanic(err)

	reprotJson, err := cfg.TakeJson("MQTTReportCfg")
	errt.Errpanic(err)
	reportCfg := mqtt.Cfg{}
	reportCfg.ClientID = fmt.Sprintf("%+v", rand.Int())
	err = jsont.Decode(reprotJson, &reportCfg)
	errt.Errpanic(err)

	m := &mqttHandle{sendHandle: mqtt.New(sendCfg), recvHandle: mqtt.New(recvCfg), reportHanlde: mqtt.New(reportCfg)}

	log.Info("server mqtt init success")
	return m
}

func (p *mqttHandle) Connect() error {
	errt.Errpanic(p.sendHandle.Connect())
	errt.Errpanic(p.recvHandle.Connect())
	errt.Errpanic(p.reportHanlde.Connect())
	return nil
}

func (p *mqttHandle) DisConnect() error {
	errt.Errpanic(p.sendHandle.DisConnect())
	errt.Errpanic(p.recvHandle.DisConnect())
	errt.Errpanic(p.reportHanlde.DisConnect())
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
