package main

import (
	"crypto/rand"
	"log/slog"
	"math/big"
	"net"
	"sync"
)

type Client struct {
	conn   *net.Conn
	UID    string
	Nick   string
	Ident  string
	Gecos  string
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

func NewLocalClient(conn *net.Conn) (*Client, error) {
	client := &Client{
		conn:   conn,
		Server: self,
		State:  ClientStatePreRegistration,
		Nick:   "*",
	}
	for range 10 {
		uid_ := []byte(self.SID)
		for range 6 {
			randint, err := rand.Int(rand.Reader, big.NewInt(26))
			if err != nil {
				return nil, err
			}
			uid_ = append(uid_, byte(65+randint.Uint64()))
		}
		uid := string(uid_)
		_, exists := uidToClient.LoadOrStore(uid, client)
		if !exists {
			client.UID = uid
			return client, nil
		}
	}
	return nil, ErrUIDBusy
}

func (client *Client) checkRegistration() {
	if client.State != ClientStatePreRegistration {
		slog.Error("spurious call to checkRegistration", "client", client)
		return // TODO: Return an error?
	}
	if client.Nick != "*" && client.Ident != "" {
		client.Send(MakeMsg(self, RPL_WELCOME, client.Nick, "Welcome"))
	}
}

type ClientState uint8

const (
	ClientStateRemote ClientState = iota
	ClientStatePreRegistration
	ClientStateRegistered
)

var (
	uidToClient  = sync.Map{}
	nickToClient = sync.Map{}
)
