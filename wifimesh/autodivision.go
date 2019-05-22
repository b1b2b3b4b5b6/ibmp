package wifimesh

var chDivisionInfo = make(chan *divisionInfo, 1000)

type divisionInfo struct {
	ip  string
	num int
}

func init() {
	go divisionLoop()
}

func (p *divisionInfo) handle() {

}

func (p *divisionInfo) add() {
	log.Debug("new uninit mesh find ip[%s] num[%d]", p.ip, p.num)
	chDivisionInfo <- p
}

func divisionLoop() {
	for {
		m := <-chDivisionInfo
		log.Debug("division info recv: %+v", m)
	}
}
