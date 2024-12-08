package main

import (
	"net"
)

type Client struct {
	conn   net.Conn
	UID    [6]byte
	Nick   string
	Ident  string
	Host   string
	Server Server
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

func (client Client) ClientSource() string {
	// TODO: Edge cases where these aren't available
	return client.Nick + "!" + client.Ident + "@" + client.Host
}

func (client Client) ServerSource() string {
	return string(client.Server.SID[:]) + string(client.UID[:])
}
