package main

import (
	"goc/logface"
	"goc/toolcom/cfgtool"
	_ "goc/toolcom/debtool"
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
