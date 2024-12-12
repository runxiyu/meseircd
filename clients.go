package main

import (
	"net"
	"sync"

	"git.sr.ht/~runxiyu/meseircd/meselog"
)

type Client struct {
	conn   *net.Conn
	CID    uint64
	Nick   string
	Ident  string
	Gecos  string
	Host   string
	Caps   map[string]struct{}
	Extra  map[string]any
	Server *Server
	State  ClientState
}

func (client *Client) Send(msg SMsg) error {
	return client.SendRaw(msg.ClientSerialize())
}

func (client *Client) SendRaw(s string) error {
	if client.conn == nil {
		panic("not implemented")
	}
	meselog.Debug("send", "line", s, "client", client.CID)
	_, err := (*client.conn).Write([]byte(s))
	if err != nil {
		// TODO: Should shut down the netFd instead but the stdlib
		// doesn't expose a way to do this.
		(*client.conn).Close()
		return err
	}
	return nil
}

func (client *Client) ClientSource() string {
	// TODO: Edge cases where these aren't available
	return client.Nick + "!" + client.Ident + "@" + client.Host
}

func (client *Client) ServerSource() uint64 {
	return client.CID
}

func (client *Client) Teardown() {
	if client.conn != nil {
		(*client.conn).Close()
	}
	if !cidToClient.CompareAndDelete(client.CID, client) {
		meselog.Error("cid inconsistent", "cid", client.CID, "client", client)
	}
	if client.State >= ClientStateRegistered || client.Nick != "*" {
		if !nickToClient.CompareAndDelete(client.Nick, client) {
			meselog.Error("nick inconsistent", "nick", client.Nick, "client", client)
		}
	}
}

func NewLocalClient(conn *net.Conn) (*Client, error) {
	var cidPart uint32
	{
		cidPartCountLock.Lock()
		defer cidPartCountLock.Unlock()
		if cidPartCount == ^uint32(0) { // UINT32_MAX
			return nil, ErrFullClients
		}
		cidPartCount++
		cidPart = cidPartCount
	}
	client := &Client{
		conn:   conn,
		Server: self,
		State:  ClientStatePreRegistration,
		Nick:   "*",
		Caps:   make(map[string]struct{}),
		Extra:  make(map[string]any),
		CID:    uint64(self.SID)<<32 | uint64(cidPart),
	}
	return client, nil
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
	cidToClient      = sync.Map{}
	nickToClient     = sync.Map{}
	cidPartCount     uint32
	cidPartCountLock sync.Mutex
)
