package main

import (
	"net"
	"log/slog"
	"sync"
)

type Client struct {
	conn   *net.Conn
	UID    string
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
	return client.UID
}

func (client *Client) Teardown() {
	if client.conn != nil {
		(*client.conn).Close()
	}
	if !uidToClient.CompareAndDelete(client.UID, client) {
		slog.Error("uid inconsistent", "uid", client.UID, "client", client)
	}
	if !nickToClient.CompareAndDelete(client.Nick, client) {
		slog.Error("nick inconsistent", "nick", client.Nick, "client", client)
	}
}

type ClientState uint8

const (
	ClientStateRemote ClientState = iota
	ClientStatePreRegistration
	ClientStateRegistered
)

var uidToClient = sync.Map{}
var nickToClient = sync.Map{}
