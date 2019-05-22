package bacnet

import (
	"fmt"
	"github.com/lumosin/goc/logface"
	"github.com/lumosin/goc/tl/cfgt"
	"github.com/lumosin/goc/tl/errt"
	"github.com/lumosin/goc/tl/jsont"
	"ibmp/mcb/devser"

	"net/http"
)

var log = logface.New(logface.TraceLevel)

func init() {
	http.HandleFunc("/bacnet/get_info", get_handle)
	cfg := cfgt.New("conf.json")
	port, err := cfg.TakeInt("BacnetServerPort")
	errt.Errpanic(err)
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	log.Info("bacnet server on port[%d]", port)
}

func get_handle(w http.ResponseWriter, r *http.Request) {
	var response struct {
		ErrCode int
		Message string
		Data    interface{}
	}
	response.ErrCode = 0
	response.Message = "success"
	response.Data = devser.GetTypDev("UIT")
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(jsont.Encode(&response)))
	return
}
