package meshifs

import (
	"errors"
	"fmt"
	"goc/toolcom/converttool"
	"goc/toolcom/errtool"
	"goc/toolcom/jsontool"
	"net"
)

const (
	OK   = 0
	BUSY = 1
)

type Udpifs struct {
	MeshID string
	IP     string
}

func init() {
	go recvServer()
	go udpServer()
}

func (p *Udpifs) GetMeshID() string {
	return p.MeshID
}

var packetID = uint8(1)

func (p *Udpifs) Destroy() {

}

func getPacketID() uint8 {
	packetID++
	if packetID == 0 {
		packetID = 1
	}
	return packetID
}

func addShell(data []byte) []byte {
	m := make([]byte, 1)
	m[0] = getPacketID()
	m = append(m, data...)
	log.Trace("add packetID[%d]", m[0])
	return m
}

func (p *Udpifs) Send(data []byte) error {
	data = addShell(data)
	port, err := cfg.TakeInt("UdpSendPort")
	errtool.Errpanic(err)
	dstAddr := &net.UDPAddr{IP: net.ParseIP(p.IP), Port: port}
	conn, err := net.DialUDP("udp", nil, dstAddr)
	errtool.Errpanic(err)
	defer conn.Close()

	for n := 0; n < 3; n++ {
		if _, err = conn.Write(data); err != nil {
			log.Warn("udp send fail[%s]", p.IP)
			return errors.New(fmt.Sprintf("udp send fail[%s]", p.IP))
		}

		recv_data := make([]byte, MaxifsLen)

		len, _ := conn.Read(recv_data)
		if len == 1 {
			switch recv_data[0] {

			case OK:
				log.Trace("ucp write data[%s]", string(data))
				return nil

			case BUSY:
				return errors.New(fmt.Sprintf("ip[%s] is busy", p.IP))

			default:
				log.Panic("udp res undefine[%+v]", recv_data)
			}
		}
		log.Warn("ip[%s] ack fail[%+v]", p.IP, recv_data)
	}
	return errors.New(fmt.Sprintf("ip[%s] no ack", p.IP))
}

func recvServer() {
	recv_port, err := cfg.TakeInt("UdpRecvPort")
	errtool.Errpanic(err)
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", recv_port))
	errtool.Errpanic(err)

	conn, err := net.ListenUDP("udp", addr)
	errtool.Errpanic(err)
	log.Info("udp recv run in port[%d]", recv_port)

	for {
		data := make([]byte, MaxifsLen)
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		ip := remoteAddr.IP.String()
		data = data[:n]
		go func() {
			conn.WriteToUDP([]byte("yes, i received"), remoteAddr)
			m := &RecvRaw{}
			m.MeshID = converttool.Mac2Str(data[:6])
			m.Typ = data[6]
			m.Data = data[7:]
			log.Debug("recv data: mesh_id[%s] typ[%d] len[%d]", m.MeshID, m.Typ, len(m.Data))
			chRecvRaw <- m

			mc := Udpifs{}
			mc.IP = ip
			mc.MeshID = m.MeshID
			chMeshCon <- &mc
		}()
	}
}

func udpServer() {
	udp_port, err := cfg.TakeInt("UdpFindPort")
	errtool.Errpanic(err)
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", udp_port))
	errtool.Errpanic(err)

	conn, err := net.ListenUDP("udp", addr)
	errtool.Errpanic(err)

	data := make([]byte, MaxifsLen)
	for {
		len, remoteAddr, err := conn.ReadFromUDP(data)
		errtool.Errpanic(err)
		go func() {
			data = data[:len]
			conn.WriteToUDP([]byte("yes, i received"), remoteAddr)
			var m struct {
				MeshID string
			}
			jsontool.Decode(string(data), &m)

			mc := Udpifs{}
			mc.IP = remoteAddr.IP.String()
			mc.MeshID = m.MeshID
			chMeshCon <- &mc
			log.Debug("udp find meshid[%s] ip[%s]", mc.MeshID, mc.IP)
		}()

	}
}