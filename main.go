package main

import (
	"github.com/lumosin/goc/logface"
	"github.com/lumosin/goc/tl/cfgt"
	_ "github.com/lumosin/goc/tl/debt"
	_ "ibmp/bacnet"
	_ "ibmp/mcb"
	"ibmp/mcb/server"
	_ "ibmp/meshdebug"
	_ "ibmp/wifimesh"
	_ "net/http/pprof"
	"time"
)

var log = logface.New(logface.TraceLevel)
var cfg = cfgtool.New("conf.json")

func init() {

}

func main() {
	time.Sleep(time.Second * 5)
	server.New().SendRequest("{\"Typ\":\"askDevices\"}")

	select {} // block
}
