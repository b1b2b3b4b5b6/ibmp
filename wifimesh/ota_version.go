package wifimesh

import (
	"fmt"
	"goc/toolcom/debtool"
	"github.com/lumosin/goc/tl/errt"
	"goc/toolcom/iotool"
	"github.com/lumosin/goc/tl/jsont"
	"io"
	"net/http"
	"os"
)

const (
	LIG  = "LIG"
	SOD  = "SOD"
	TEMP = "TEMP"
)

var otaCfg struct {
	LocalIP   string
	LocalPort int
	LocalPath string
}

var typBinMap = map[string]string{}

func init() {
	json, err := cfg.TakeJson("OtaCfg")
	errtool.Errpanic(err)
	err = jsontool.Decode(json, &otaCfg)
	errtool.Errpanic(err)

	otaCfg.LocalPath = iotool.GetCurrentDirectory()
	go http.ListenAndServe(fmt.Sprintf(":%d", otaCfg.LocalPort), http.FileServer(http.Dir(otaCfg.LocalPath)))
	log.Info("ota version server running[%+v]", otaCfg)

	debtool.AddFunc("setBin", setBin)
}

func setBin(typ string, url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/%s", otaCfg.LocalPath, typ))
	if err != nil {
		return err
	}

	defer f.Close()
	io.Copy(f, res.Body)
	typBinMap[typ] = typ
	log.Debug("typ[%s] bin set url[%s]", typ, url)
	return nil
}

func getBinUrl(typ string) string {
	name, ok := typBinMap[typ]
	if !ok {
		return ""
	}

	url := fmt.Sprintf("http://%s:%d/%s", otaCfg.LocalIP, otaCfg.LocalPort, name)
	return url
}
