package server

import (
	"fmt"
	"github.com/lumosin/goc/logface"
	"github.com/lumosin/goc/tl/cfgt"
	"github.com/lumosin/goc/tl/errt"
)

const (
	StatusCon    = 1
	StatusDisCon = 2
)

var log = logface.New(logface.DebugLevel)

type Server interface {
	Report(str string) error
	WaitRequest() string
	SendRequest(str string) error
	Connect() error
	DisConnect() error
}

var Ser Server

func init() {
	cfg := cfgt.New("conf.json")
	method, err := cfg.TakeString("ServerMethod")
	errt.Errpanic(err)

	switch method {
	case "MQTT":
		Ser = NewMQTT()
	case "HTTP":
		Ser = NewHttp()
	default:
		log.Panic(fmt.Sprintf("server meshthod[%s] not exist", method))
	}
	Ser.Connect()
	log.Info("server connect success")
}

func New() Server {
	return Ser
}
