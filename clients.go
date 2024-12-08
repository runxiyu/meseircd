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
	CID    string
	Nick   string
	Ident  string
	Gecos  string
	Host   string
	Caps   map[string]struct{}
	Extra  map[string]any
	Server Server
	State  ClientState
}

func (client *Client) Send(msg SMsg) error {
	return client.SendRaw(msg.ClientSerialize())
}

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
		return err
	}
	return nil
}

func (client Client) ClientSource() string {
	// TODO: Edge cases where these aren't available
	return client.Nick + "!" + client.Ident + "@" + client.Host
}

func (client Client) ServerSource() string {
	return client.CID
}

func (client *Client) Teardown() {
	if client.conn != nil {
		(*client.conn).Close()
	}
	if !cidToClient.CompareAndDelete(client.CID, client) {
		slog.Error("cid inconsistent", "cid", client.CID, "client", client)
	}
	if client.State >= ClientStateRegistered || client.Nick != "*" {
		if !nickToClient.CompareAndDelete(client.Nick, client) {
			slog.Error("nick inconsistent", "nick", client.Nick, "client", client)
		}
	}
}

func NewLocalClient(conn *net.Conn) (*Client, error) {
	client := &Client{
		conn:   conn,
		Server: self,
		State:  ClientStatePreRegistration,
		Nick:   "*",
		Caps:   make(map[string]struct{}),
		Extra:  make(map[string]any),
	}
	for range 10 {
		cid_ := []byte(self.SID)
		for range 6 {
			randint, err := rand.Int(rand.Reader, big.NewInt(26))
			if err != nil {
				return nil, err
			}
			cid_ = append(cid_, byte(65+randint.Uint64()))
		}
		cid := string(cid_)
		_, exists := cidToClient.LoadOrStore(cid, client)
		if !exists {
			client.CID = cid
			return client, nil
		}
	}
	return nil, ErrCIDBusy
}

func (client *Client) checkRegistration() error {
	switch client.State {
	case ClientStatePreRegistration:
	case ClientStateCapabilitiesFinished:
	default:
		return nil
	}
	if client.Nick == "*" || client.Ident == "" {
		return nil
	}
	client.State = ClientStateRegistered
	err := client.Send(MakeMsg(self, RPL_WELCOME, client.Nick, "Welcome to the rxIRC network, "+client.Nick))
	if err != nil {
		return err
	}
	err = client.Send(MakeMsg(self, RPL_YOURHOST, client.Nick, "Your host is "+self.Name+", running version "+VERSION))
	if err != nil {
		return err
	}
	err = client.Send(MakeMsg(self, RPL_CREATED, client.Nick, "This server was created 1970-01-01 00:00:00 UTC"))
	if err != nil {
		return err
	}
	err = client.Send(MakeMsg(self, RPL_MYINFO, client.Nick, self.Name, VERSION, "", "", ""))
	if err != nil {
		return err
	}
	err = client.Send(MakeMsg(self, RPL_ISUPPORT, "YAY=", "are supported by this server"))
	if err != nil {
		return err
	}
	return nil
}

type ClientState uint8

const (
	ClientStatePreRegistration ClientState = iota
	ClientStateCapabilities
	ClientStateCapabilitiesFinished
	ClientStateRegistered
	ClientStateRemote
)

var (
	cidToClient  = sync.Map{}
	nickToClient = sync.Map{}
)
