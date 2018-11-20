package client

import (
	"net"
	"github.com/msaldanha/realChain/network"
	"github.com/msaldanha/realChain/transaction"
	"bytes"
	"github.com/davecgh/go-xdr/xdr2"
	"github.com/kataras/iris/core/errors"
)

type Client struct {
	conn *net.UDPConn
}

func New(serverUrl string) (*Client, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", serverUrl)
	if err != nil {
		return nil, err
	}

	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		return nil, err
	}

	return &Client{conn:conn}, nil
}

func (rs *Client) HandleTransaction(tx *transaction.Transaction) error {
	err := rs.send("transaction.handle", tx.ToBytes())
	if err != nil {
		return err
	}
	data, err := rs.receive()
	if err != nil {
		return err
	}
	msg := network.NewMessageFromBytes(data)
	if msg != nil {
		return err
	}
	var errMsg string
	decoder := xdr.NewDecoder(bytes.NewReader(msg.Payload))
	decoder.Decode(&errMsg)
	if errMsg != "" {
		return errors.New(errMsg)
	}
	return nil
}

func (rs *Client) send(endPoint string, data []byte) error {
	msg := network.NewMessage()
	msg.EndPoint = endPoint
	msg.Payload = data
	_, err := rs.conn.Write(msg.ToBytes())
	return err
}

func (rs *Client) receive() ([]byte, error) {
	var buf [1024]byte
	size, _, err := rs.conn.ReadFrom(buf[0:])
	if err != nil {
		return nil, err
	}
	return buf[0:size], nil
}