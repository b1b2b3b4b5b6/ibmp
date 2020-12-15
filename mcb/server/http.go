package server

import (
	"fmt"
	"github.com/b1b2b3b4b5b6/goc/tl/cfgt"
	"github.com/b1b2b3b4b5b6/goc/tl/errt"
	"github.com/b1b2b3b4b5b6/goc/tl/jsont"
	"math/rand"
	"net/http"
	"time"
)

type httpHandle struct {
	LocalPort int
	LocalPath string
}

var chHttpRecv = make(chan string, 1000)

func NewHttp() *httpHandle {
	rand.Seed(time.Now().Unix())
	cfg := cfgt.New("conf.json")

	sendJson, err := cfg.TakeJson("HttpCfg")
	errt.Errpanic(err)

	m := httpHandle{}
	jsont.Decode(sendJson, &m)
	log.Info("http server init success")
	return &m
}

func (p *httpHandle) Connect() error {
	http.HandleFunc(p.LocalPath, localHandle)
	go http.ListenAndServe(fmt.Sprintf(":%d", p.LocalPort), nil)
	log.Info("http server work on port[%d]", p.LocalPort)
	return nil
}

func (p *httpHandle) DisConnect() error {

	return nil
}

func (p *httpHandle) SendRequest(str string) error {
	log.Debug("http send[%s]", str)
	return nil
}

func (p *httpHandle) WaitRequest() string {
	recv := <-chHttpRecv
	log.Debug("http recv[%s]", recv)
	return recv
}

func (p *httpHandle) Report(str string) error {
	log.Debug("http report[%s]", str)
	return nil
}

func localHandle(w http.ResponseWriter, r *http.Request) {
	var res struct {
		StatusCode int
		Message    string      `json:"Message,omitempty"`
		Result     interface{} `json:"Result,omitempty"`
	}
	defer func() {
		w.Write([]byte(jsont.Encode(&res)))
	}()

	err := r.ParseForm()
	if err != nil {
		res.StatusCode = -1
		res.Message = err.Error()
		return
	}
	data := make([]byte, 1024*1024)
	rlen, err := r.Body.Read(data)
	data = data[:rlen]
	chHttpRecv <- string(data)
}
