package meshifs

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/lumosin/goc/tl/turnt"
	"github.com/lumosin/goc/tl/errt"
	"github.com/lumosin/goc/tl/jsont"
	"net"
	"time"
)

var tcpPortCfg struct{
	TcpRecvPort int
}

func init() {
	json, err := cfg.TakeJson("MeshIfs")
	errt.Errpanic(err)
	err = jsont.Decode(json, &tcpPortCfg)
	errt.Errpanic(err)
	go tcpRecv()
}

type Tcpifs struct {
	MeshID string
	Conn   *net.TCPConn
}

func (p *Tcpifs) Destroy() {
	err := p.Conn.Close()
	if err != nil {
		log.Error("tcp connect destroy fail[%+v]", err)
	}
}

func (p *Tcpifs) Send(data []byte) error {

	if p.Conn == nil {
		log.Warn("conn is nil")
		return errors.New("no connection")
	}

	err := p.Conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
	if err != nil {
		log.Warn("SetWriteDeadline fail, err[%s]", err)
		p.Conn.Close()
		return err
	}

	send_len, err := p.Conn.Write(data)
	log.Trace("tcp write data[%s]", string(data))
	if err != nil {
		log.Warn("tcp write fail, err[%s]", err)
		p.Conn.Close()
		return err
	}

	errt.Assert(send_len == len(data))
	time.Sleep(time.Millisecond * 100)
	return nil
}

func (p *Tcpifs) GetMeshID() string {
	return p.MeshID
}

func pakcetSplit(data []byte) {
	log.Debug("split data[%d] start", len(data))
	for {
		if len(data) == 0 {
			break
		}

		if len(data) < 2 {
			log.Error("recv data len[%d] error, must > 2", len(data))
			break
		}

		m := &RecvRaw{}
		data_len := binary.BigEndian.Uint16(data[:2])
		data = data[2:]

		if len(data) < int(data_len) {
			log.Error("recv data len[%d] error, suppose to be [%d]", len(data), data_len)
			break
		}

		m.MeshID = turnt.Mac2Str(data[:6])
		m.Typ = data[6]
		m.Data = data[7:data_len]
		log.Debug("recv data: mesh_id[%s] typ[%d] len[%d]", m.MeshID, m.Typ, len(m.Data))
		chRecvRaw <- m
		data = data[data_len:]
	}
}

func tcpHandle(conn *net.TCPConn) {
	defer func() {
		log.Debug("close socket[%s]", conn.RemoteAddr().String())
		conn.Close()
	}()
	log.Debug("tcp handle start[%s]", conn.RemoteAddr().String())
	for {
		data := make([]byte, MaxifsLen*100)
		recv_len, err := conn.Read(data)

		if err != nil {
			netErr := err.(net.Error)
			switch {
			case netErr.Timeout():
				continue
			default:
				log.Error("[%s]read fail [%s]", conn.RemoteAddr().String(), err.Error())
				return
			}
		}

		data = data[:recv_len]
		pakcetSplit(data)
	}

}

func tcpRecv() {

	localAddress, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", tcpPortCfg.TcpRecvPort))
	errt.Errpanic(err)
	tcpListener, err := net.ListenTCP("tcp", localAddress)
	errt.Errpanic(err)
	log.Info("tcp recv run in port[%d]", tcpPortCfg.TcpRecvPort)
	defer func() {
		tcpListener.Close()
	}()

	for {
		conn, err := tcpListener.AcceptTCP()
		errt.Errpanic(err)

		data := make([]byte, MaxifsLen)
		recv_len, err := conn.Read(data)
		if err != nil {
			log.Error("after connect, no meshid recv")
			conn.Close()
			continue
		}
		var m struct {
			MeshID string
		}
		jsont.Decode(string(data[:recv_len]), &m)

		mc := Tcpifs{}
		mc.MeshID = m.MeshID
		mc.Conn = conn
		chMeshCon <- &mc

		log.Debug("meshid[%s] ip[%s] connect", m.MeshID, conn.RemoteAddr().String())
		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(time.Second * 30)
		conn.SetNoDelay(true)
		go tcpHandle(conn)
	}

}
