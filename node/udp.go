package node

import (
	"net"
	"github.com/msaldanha/realChain/ledger"
	log "github.com/sirupsen/logrus"
)

type UdpServer struct {
	ld *ledger.Ledger
	url string
}

func NewUdpServer(ld *ledger.Ledger, url string) (*UdpServer) {
	return &UdpServer{ld: ld, url: url}
}

func (n *UdpServer) Run() (error) {
	log.Info("Udp server starting")
	udpAddr, err := net.ResolveUDPAddr("udp", n.url)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

	log.Infof("Udp server listening at %s", n.url)
	for {
		n.handleClient(conn)
	}
}

func (n *UdpServer) handleClient(conn *net.UDPConn) {
	var buf [1024]byte
	size, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}
	err = n.ld.HandleTransactionBytes(buf[0:size])
	if err == nil {
		conn.WriteToUDP([]byte("OK"), addr)
	}
	conn.WriteToUDP([]byte(err.Error()), addr)
}