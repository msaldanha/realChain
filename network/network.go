package network

import (
	"bytes"
	"github.com/davecgh/go-xdr/xdr2"
	"net"
	"fmt"
)

type Network struct {
	conn     *net.UDPConn
	handlers map[string][]Handler
	peers    map[string]*Peer
	url      string
}

type Context struct {
	EndPoint string
	Peer     *Peer
	Data     []byte
}

type Handler func(*Context)

func NewNetwork(url string) (*Network, error) {
	destUdpAddr, err := net.ResolveUDPAddr("udp4", url)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp4", destUdpAddr)
	if err != nil {
		fmt.Println(err)
	}

	return &Network{conn: conn, handlers: make(map[string][]Handler), peers: make(map[string]*Peer), url:url}, nil
}

func (n *Network) InstallHandler(endPoint string, handlers ...Handler) error {
	if len(endPoint) == 0 {
		return nil
	}

	if len(handlers) == 0 {
		return nil
	}

	h := n.handlers[endPoint]
	if h == nil {
		h = handlers
	} else {
		h = append(h, handlers...)
	}
	n.handlers[endPoint] = h

	return nil
}

func (n *Network) UsePeers(urls ...string) {
	for _, url := range urls {
		peer, err := NewPeer(url, n.conn)
		if err == nil {
			n.peers[peer.String()] = peer
		}
	}
}

func (n *Network) handle(peer *Peer, data []byte) error {
	var msg message
	n.Log("Received msg from %s \n", peer)
	decoder := xdr.NewDecoder(bytes.NewReader(data))
	_, err := decoder.Decode(&msg)
	if err != nil {
		return nil
	}

	if msg.Magic != Magic {
		return nil
	}

	n.Log("Received msg from %s to endpoint %s\n", peer, msg.EndPoint)

	handlers := n.handlers[msg.EndPoint]
	if handlers == nil {
		n.Log("No handlers to handle endpoint %s\n", msg.EndPoint)
		return nil
	}

	if len(handlers) == 0 {
		return nil
	}

	ctx := &Context{EndPoint: msg.EndPoint, Peer: peer, Data: msg.Payload}
	for _, callback := range handlers {
		callback(ctx)
	}

	return nil
}

type handler func(peer *Peer, data []byte, err error) (error)

func (n *Network) Receive() ([]byte, *net.UDPAddr, error) {
	var buf [1024]byte
	size, addr, err := n.conn.ReadFromUDP(buf[0:])
	if err != nil {
		return nil, nil, err
	}
	return buf[0:size], addr, nil
}

func (n *Network) Run() {
	for {
		data, addr, err := n.Receive()
		if err != nil {
			return
		}
		key := addr.String()
		peer := n.peers[key]
		if peer != nil {
			n.handle(peer, data)
		}
	}
}

func (n *Network) Broadcast(endPoint string, data []byte) (error) {
	msg := NewMessage()
	msg.EndPoint = endPoint
	msg.Payload = data
	msgBytes := msg.ToBytes()
	for _, peer := range n.peers {
		peer.Send(msgBytes)
	}
	return nil
}

func (n *Network) Log(format string, values ...interface{}) {
	fmt.Printf(n.url + ": " + format, values...)
}
