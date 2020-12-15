/*
 * @Author: your name
 * @Date: 2020-12-15 01:27:14
 * @LastEditTime: 2020-12-15 01:51:43
 * @LastEditors: your name
 * @Description: In User Settings Edit
 * @FilePath: /ibmp/main.go
 */
package main

import (
	"github.com/b1b2b3b4b5b6/goc/logface"
	"github.com/b1b2b3b4b5b6/goc/tl/cfgt"
	_ "github.com/b1b2b3b4b5b6/goc/tl/debt"
	_ "ibmp/bacnet"
	_ "ibmp/mcb"
	"ibmp/mcb/server"
	_ "ibmp/meshdebug"
	_ "ibmp/wifimesh"
	_ "net/http/pprof"
	"time"
)

var log = logface.New(logface.TraceLevel)
var cfg = cfgt.New("conf.json")

func init() {

}

func main() {
	time.Sleep(time.Second * 5)
	server.New().SendRequest("{\"Typ\":\"askDevices\"}")

	select {} // block
}
