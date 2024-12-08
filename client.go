package main

import (
	"net"
)

type Client struct {
	conn net.Conn
	uid  [6]byte
}

func (client *Client) Send(msg SMsg) {
	client.SendRaw(msg.ClientSerialize())
}

// Send failures are not returned; broken connections detected and severed on
// the next receive.
func (client *Client) SendRaw(s string) {
	_, err := client.conn.Write([]byte(s))
	if err != nil {
		// TODO: Should shut down the netFd instead but the stdlib
		// doesn't expose a way to do this.
		client.conn.Close()
	}
}
