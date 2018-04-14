package node

import (
	"net"
	"github.com/msaldanha/realChain/ledger"
	log "github.com/sirupsen/logrus"
)

type UdpServer struct {
	ld *ledger.Ledger
}

func NewUdpServer(ld *ledger.Ledger) (*UdpServer) {
	return &UdpServer{ld: ld}
}

func (n *UdpServer) Run() (error) {
	log.Info("Udp server starting")
	service := ":1200"
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

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