package main

import (
	"net"
	"log/slog"
)

type Client struct {
	conn   *net.Conn
	UID    [6]byte
	Nick   string
	Ident  string
	Host   string
	Server Server
	State  ClientState
}

func (client *Client) Send(msg SMsg) error {
	return client.SendRaw(msg.ClientSerialize())
}

// Send failures are not returned; broken connections detected and severed on
// the next receive.
func (client *Client) SendRaw(s string) error {
	if client.conn == nil {
		panic("not implemented")
	}
	slog.Debug("send", "line", s, "conn", client.conn)
	_, err := (*client.conn).Write([]byte(s))
	if err != nil {
		// TODO: Should shut down the netFd instead but the stdlib
		// doesn't expose a way to do this.
		(*client.conn).Close()
	}
	return nil
}

func (client Client) ClientSource() string {
	// TODO: Edge cases where these aren't available
	return client.Nick + "!" + client.Ident + "@" + client.Host
}

func (client Client) ServerSource() string {
	return string(client.Server.SID[:]) + string(client.UID[:])
}

type ClientState uint8

const (
	ClientStateRemote ClientState = iota
	ClientStatePreRegistration
	ClientStateRegistered
)
