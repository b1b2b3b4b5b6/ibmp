package meshdebug

import (
	"encoding/json"
	"fmt"
	"goc/logface"
	"goc/toolcom/cfgtool"
	"goc/toolcom/errtool"
	"goc/toolcom/jsontool"
	"ibmp/wifimesh"
	"net/http"
)

var log = logface.New(logface.DebugLevel)

type response struct {
	StatusCode int
	Message    string      `json:"Message,omitempty"`
	Result     interface{} `json:"Result,omitempty"`
}

type requestArg struct {
	Cmd     string
	MeshID  string
	Macs    []string
	CusData interface{}
}

func init() {
	http.HandleFunc("/debug/mesh", debugMesh)
	cfg := cfgtool.New("conf.json")
	port, err := cfg.TakeInt("MeshDebugPort")
	errtool.Errpanic(err)
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	log.Info("mesh debug work on port[%d]", port)
}

func debugMesh(w http.ResponseWriter, r *http.Request) {
	res := response{}
	defer func() {
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(jsontool.Encode(&res)))
	}()

	err := r.ParseForm()
	errtool.Errpanic(err)
	var request requestArg
	err = json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		res.StatusCode = -1
		res.Message = err.Error()
		return
	}

	meshGroup := wifimesh.GetGroup()
	mesh := meshGroup.GetMesh(request.MeshID)
	if mesh == nil {
		res.StatusCode = -1
		res.Message = "can not find mesh"
		return
	}

	switch request.Cmd {

	case "Ota":
		cus_data := request.CusData.(map[string]interface{})
		peroid_ms := cus_data["PeroidMs"].(float64)
		typ := cus_data["Typ"].(string)
		is_http := cus_data["IsHttp"].(float64)
		mesh.Ota(request.Macs, typ, int(peroid_ms), int(is_http))

	case "Destroy":
		mesh.Destroy()

	case "ReportGather":
		mesh.ReportGather()

	case "ReportNum":
		mesh.ReportNum()

	case "SendCustom":
		mesh.SendCustom(request.Macs, request.CusData)
	default:
		res.Message = "can not find Cmd:" + request.Cmd
		res.StatusCode = -1
		return
	}
	return
}
