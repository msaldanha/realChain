package network

import "net"

type Peer struct {
	conn *net.UDPConn
	addr *net.UDPAddr
	url  string
}

func NewPeer(url string, conn *net.UDPConn) (*Peer, error) {
	addr, err := net.ResolveUDPAddr("udp4", url)
	if err != nil {
		return nil, err
	}
	return &Peer{conn: conn, addr: addr, url: url}, nil
}

func (p *Peer) Send(endPoint string, data []byte) error {
	msg := NewMessage()
	msg.EndPoint = endPoint
	msg.Payload = data
	_, err := p.conn.WriteToUDP(msg.ToBytes(), p.addr)
	return err
}

func (p *Peer) String() string {
	return p.addr.String()
}
